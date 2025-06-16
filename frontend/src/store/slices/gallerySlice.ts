import { NewsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';
import { transformMinioUrlsInData } from '@/utils/minioUrlHelper';

export type Image = {
  id: string;
  image_url: string;
};

type ResponseWithPagination = {
  images: Image[];
  totalLength: number;
};

type ImageState = {
  images: Image[];
  paginatedData: ResponseWithPagination;
  loading: boolean;
  error: string | null;
};

const initialState: ImageState = {
  images: [],
  loading: false,
  error: null,
  paginatedData: {
    images: [],
    totalLength: 0
  }
};

export const fetchImages = createAsyncThunk(
  'gallery/fetchImages',
  async (_, { rejectWithValue }) => {
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

export const fetchImageById = createAsyncThunk(
  'gallery/fetchImageById',
  async (id: string, { rejectWithValue }) => {
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

export const fetchImagesWithPagination = createAsyncThunk(
  'gallery/fetchImagesWithPagination',
  async (query: { page: number; limit: number }, { rejectWithValue }) => {
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

export const removeImage = createAsyncThunk(
  'gallery/removeImage',
  async (id: string, { rejectWithValue }) => {
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

// Note: These helper functions are now used in the thunks
// They're kept here for reference and potential future use
// Helper function to convert file to base64
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });
};

export const addNewImage = createAsyncThunk(
  'gallery/addNewImage',
  async (values: NewsFormInput, { rejectWithValue }) => {
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

      // Convert file to base64
      const imageData = await fileToBase64(file);

      // Send the post with base64-encoded image data
      const response = await axiosInstance.post<Image>('/gallery', {
        image_data: imageData
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
  extraReducers: (builder) => {
    builder
      .addCase(fetchImages.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchImages.fulfilled, (state, action) => {
        state.images = action.payload as Image[];
        state.loading = false;
      })
      .addCase(fetchImagesWithPagination.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchImagesWithPagination.fulfilled, (state, action) => {
        state.paginatedData = action.payload as ResponseWithPagination;
        state.loading = false;
      })
      .addCase(fetchImageById.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.images = [];
      })
      .addCase(fetchImageById.fulfilled, (state, action) => {
        state.images.push(action.payload as Image);
        state.loading = false;
      })
      .addCase(removeImage.fulfilled, (state, action) => {
        state.images = state.images.filter(
          (item) => item.id !== (action.meta.arg as string)
        );
      })
      .addCase(addNewImage.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(addNewImage.fulfilled, (state, action) => {
        if (action.payload) {
          state.images.push(action.payload as Image);
        }
        state.loading = false;
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.error = action.payload;
        state.loading = false;
      });
  }
});

export default gallerySlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}
