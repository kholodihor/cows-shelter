import { NewsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

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

export const fetchImages = createAsyncThunk('gallery/fetchImages', async (_, { rejectWithValue }) => {
  try {
    const response = await axiosInstance.get<Image[]>('/gallery');
    const data = response.data;
    return Array.isArray(data) ? data : [];
  } catch (error) {
    const err = error as AxiosError;
    console.error('Failed to fetch images:', err.message);
    return rejectWithValue('Failed to load images');
  }
});

export const fetchImageById = createAsyncThunk(
  'gallery/fetchImageById',
  async (id: string, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Image>(`/gallery/${id}`);
      return response.data;
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
      return response.data || { images: [], totalLength: 0 };
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

export const addNewImage = createAsyncThunk(
  'gallery/addNewImage',
  async (values: NewsFormInput, { rejectWithValue }) => {
    try {
      const file = values.image[0];
      const formData = new FormData();
      formData.append('image', file);
      const { data } = await axiosInstance.post<{ image_url: string }>('/upload-image', formData);
      const newImage = {
        image_url: data.image_url
      };
      await axiosInstance.post('/gallery', newImage);
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