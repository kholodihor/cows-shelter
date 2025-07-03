import { NewsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction,
  ActionReducerMapBuilder
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';
import { transformMinioUrlsInData } from '@/utils/minioUrlHelper';

export type Image = {
  ID: string; // Note: Uppercase ID to match backend response
  image_url: string;
};

type ResponseWithPagination = {
  images: Image[];
  totalLength: number;
};

export interface ImageState {
  images: Image[];
  paginatedData: ResponseWithPagination;
  loading: boolean;
  error: string | null;
}

// Helper type for error handling in async thunks
type RejectWithValue = (value: string) => any;

const initialState: ImageState = {
  images: [],
  loading: false,
  error: null,
  paginatedData: {
    images: [],
    totalLength: 0
  }
};

export const fetchImages = createAsyncThunk<Image[], void, { rejectValue: string }>(
  'gallery/fetchImages',
  async (_, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Image[]>('/gallery');
      const data = response.data;
      console.log('Original gallery data:', data);
      // Transform MinIO URLs in the response data
      const transformedData = transformMinioUrlsInData(
        Array.isArray(data) ? data : []
      );
      console.log('Transformed gallery data:', transformedData);
      return transformedData;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch images:', err.message);
      return rejectWithValue('Failed to load images');
    }
  }
);

export const fetchImageById = createAsyncThunk<Image, string, { rejectValue: string }>(
  'gallery/fetchImageById',
  async (id: string, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Image>(`/gallery/${id}`);
      // Transform MinIO URLs in the response data
      const transformedData = transformMinioUrlsInData(response.data);
      return transformedData;
    } catch (error) {
      const err = error as AxiosError;
      console.error(`Failed to fetch image with id ${id}:`, err.message);
      return rejectWithValue('Failed to load image');
    }
  }
);

export const fetchImagesWithPagination = createAsyncThunk<ResponseWithPagination, { page: number; limit: number }, { rejectValue: string }>(
  'gallery/fetchImagesWithPagination',
  async (query: { page: number; limit: number }, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
    try {
      const response = await axiosInstance.get<ResponseWithPagination>(
        `/gallery/pagination?page=${query.page}&limit=${query.limit}`
      );
      // Transform MinIO URLs in the response data
      const data = response.data || { images: [], totalLength: 0 };
      if (data.images) {
        data.images = transformMinioUrlsInData(data.images);
      }
      return data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch paginated images:', err.message);
      return rejectWithValue('Failed to load paginated images');
    }
  }
);

export const removeImage = createAsyncThunk<string, string, { rejectValue: string }>(
  'gallery/removeImage',
  async (id: string, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
    try {
      await axiosInstance.delete(`/gallery/${id}`);
      return id; // Return the deleted ID for state updates
    } catch (error) {
      const err = error as AxiosError;
      console.error(`Failed to delete image with id ${id}:`, err.message);
      return rejectWithValue('Failed to delete image');
    }
  }
);

// Helper function to ensure base64 data is in the correct format for the backend
// Backend expects: data:[content-type];base64,[base64-data]
const extractBase64Data = (data: string): string => {
  // If it's already a properly formatted data URL, return it as is
  if (data.startsWith('data:') && data.includes(';base64,')) {
    const parts = data.split(',');
    if (parts.length !== 2) {
      throw new Error('Invalid data URL format: missing comma separator');
    }
    if (!parts[0].includes('base64')) {
      throw new Error('Invalid data URL format: not base64 encoded');
    }
    return data; // Return the full data URL
  }
  
  // If it's raw base64 data, convert it to a proper data URL
  const base64Regex = /^[A-Za-z0-9+/=]+$/;
  if (!base64Regex.test(data)) {
    console.warn('Warning: Input might not be valid base64');
  }
  
  // Convert to proper data URL format with image/jpeg as default content type
  return `data:image/jpeg;base64,${data}`;
};

// Helper function to convert file to base64 data URL
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      const result = reader.result;
      if (typeof result === 'string') {
        resolve(result);
      } else {
        reject(new Error('FileReader did not return a string'));
      }
    };
    reader.onerror = error => {
      console.error('Error reading file:', error);
      reject(error);
    };
  });
};

export const addNewImage = createAsyncThunk<Image, NewsFormInput, { rejectValue: string }>(
  'gallery/addNewImage',
  async (values: NewsFormInput, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
    try {
      if (
        !values.image ||
        !values.image[0] ||
        !(values.image[0] instanceof File)
      ) {
        return rejectWithValue('Image file is required');
      }

      const file = values.image[0];
      const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

      if (file.size > MAX_FILE_SIZE) {
        return rejectWithValue('Image size should be less than 5MB');
      }

      // Convert file to base64 data URL
      const dataUrl = await fileToBase64(file);
      // Extract just the base64 data part
      const base64Data = extractBase64Data(dataUrl);

      // Send the post with base64-encoded image data
      const response = await axiosInstance.post<Image>('/gallery', {
        image_data: base64Data
      });

      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to add new image:', err.message);
      return rejectWithValue('Failed to add new image');
    }
  }
);

const gallerySlice = createSlice({
  name: 'gallery',
  initialState,
  reducers: {},
  extraReducers: (builder: ActionReducerMapBuilder<ImageState>) => {
    builder
      .addCase(fetchImages.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchImages.fulfilled, (state, action: PayloadAction<Image[]>) => {
        state.images = action.payload;
        state.loading = false;
      })
      .addCase(fetchImagesWithPagination.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchImagesWithPagination.fulfilled, (state, action: PayloadAction<ResponseWithPagination>) => {
        state.paginatedData = action.payload;
        state.loading = false;
      })
      .addCase(fetchImageById.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.images = [];
      })
      .addCase(fetchImageById.fulfilled, (state, action: PayloadAction<Image>) => {
        state.images = [action.payload];
        state.loading = false;
      })
      .addCase(fetchImageById.rejected, (state, action: PayloadAction<string | undefined>) => {
        state.loading = false;
        state.error = action.payload || 'Failed to fetch image';
      })
      .addCase(removeImage.fulfilled, (state, action) => {
        state.images = state.images.filter(
          (item) => item.ID !== action.meta.arg
        );
      })
      .addCase(addNewImage.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(addNewImage.fulfilled, (state, action: PayloadAction<Image>) => {
        state.images.unshift(action.payload);
        state.loading = false;
      })
      .addCase(addNewImage.rejected, (state, action: PayloadAction<string | undefined>) => {
        state.loading = false;
        state.error = action.payload || 'Failed to add new image';
      })
      .addMatcher(
        (action: AnyAction) => action.type.endsWith('rejected'),
        (state, action: PayloadAction<string | undefined>) => {
          state.error = action.payload || 'An error occurred';
          state.loading = false;
        }
      );
  }
});

export default gallerySlice.reducer;
