import { ExcursionsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';
import { transformMinioUrlsInData } from '@/utils/minioUrlHelper';

export type Excursion = {
  id: string;
  title_ua: string;
  title_en: string;
  description_ua: string;
  description_en: string;
  amount_of_persons: string;
  time_from: string;
  time_to: string;
  image_url: string;
};

type ResponseWithPagination = {
  excursions: Excursion[];
  totalLength: number;
};

type ExcursionsState = {
  excursions: Excursion[];
  loading: boolean;
  error: string | null;
  paginatedData: ResponseWithPagination;
};

const initialState: ExcursionsState = {
  excursions: [],
  loading: false,
  error: null,
  paginatedData: {
    excursions: [],
    totalLength: 0
  }
};

// Helper function to convert file to base64
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });
};

export const fetchExcursion = createAsyncThunk(
  'excursions/fetchExcursion',
  async () => {
    try {
      const response = await axiosInstance.get<Excursion[]>('/excursions');
      const data = response.data;
      // Transform MinIO URLs in the response data
      return transformMinioUrlsInData(Array.isArray(data) ? data : []);
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const fetchExcursionById = createAsyncThunk(
  'excursions/fetchExcursionById',
  async (id: string) => {
    try {
      const response = await axiosInstance.get<Excursion>(`/excursions/${id}`);
      const data = response.data;
      // Transform MinIO URLs in the response data
      return transformMinioUrlsInData(data);
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const fetchExcursionsWithPagination = createAsyncThunk(
  'excursions/fetchExcursionsWithPagination',
  async (query: { page: number; limit: number }) => {
    try {
      const response = await axiosInstance.get<ResponseWithPagination>(
        `/excursions/pagination?page=${query.page}&limit=${query.limit}`
      );
      const data = response.data;
      // Transform MinIO URLs in the response data
      if (data.excursions) {
        data.excursions = transformMinioUrlsInData(data.excursions);
      }
      return data;
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const removeExcursion = createAsyncThunk(
  'excursions/removeExcursion',
  async (id: string) => {
    try {
      await axiosInstance.delete(`/excursions/${id}`);
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const addNewExcursion = createAsyncThunk(
  'excursions/addNewExcursion',
  async (values: ExcursionsFormInput, { rejectWithValue }) => {
    try {
      // Convert image to base64 if it exists
      let imageData = '';

      if (values.image && values.image[0] && values.image[0] instanceof File) {
        try {
          const file = values.image[0];
          // Check if file size is reasonable (e.g., less than 5MB)
          const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
          if (file.size > MAX_FILE_SIZE) {
            return rejectWithValue('Image size should be less than 5MB');
          }

          imageData = await fileToBase64(file);
        } catch (uploadError) {
          console.error('Image processing failed:', uploadError);
          return rejectWithValue('Failed to process image');
        }
      } else {
        return rejectWithValue('Image is required');
      }

      // Prepare the excursion data with base64 image
      const excursionData = {
        title_ua: values.titleUa || '',
        title_en: values.titleEn || values.titleUa || '',
        description_ua: values.descriptionUa || '',
        description_en: values.descriptionEn || values.descriptionUa || '',
        amount_of_persons: values.visitorsNumber,
        time_from: values.timeFrom,
        time_to: values.timeTill,
        image_data: imageData
      };

      // Create the excursion with JSON data
      const response = await axiosInstance.post<Excursion>('/excursions', excursionData, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.data) {
        throw new Error('No data received from server');
      }

      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Error adding excursion:', err.response?.data || err.message);
      return rejectWithValue(err.response?.data || 'Failed to add excursion');
    }
  }
);

export const editExcursion = createAsyncThunk(
  'excursions/editExcursion',
  async (excursionsData: { id?: string; values: ExcursionsFormInput }, { rejectWithValue }) => {
    try {
      if (!excursionsData.id) {
        throw new Error('Excursion ID is required');
      }

      // Convert image to base64 if a new image is provided
      let imageData = '';

      if (excursionsData.values.image?.[0] && excursionsData.values.image[0] instanceof File) {
        // Check if this is a real file with actual content (not a dummy file)
        // Dummy files created in the edit form have size 0
        if (excursionsData.values.image[0].size > 0) {
          try {
            const file = excursionsData.values.image[0];
            // Check if file size is reasonable (e.g., less than 5MB)
            const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
            if (file.size > MAX_FILE_SIZE) {
              return rejectWithValue('Image size should be less than 5MB');
            }

            imageData = await fileToBase64(file);
          } catch (uploadError) {
            console.error('Image processing failed:', uploadError);
            return rejectWithValue('Failed to process image');
          }
        }
      }

      // Prepare the update data with the same structure as addNewExcursion
      const updateData: Record<string, any> = {
        title_ua: excursionsData.values.titleUa || '',
        title_en: excursionsData.values.titleEn || excursionsData.values.titleUa || '',
        description_ua: excursionsData.values.descriptionUa || '',
        description_en: excursionsData.values.descriptionEn || excursionsData.values.descriptionUa || '',
        amount_of_persons: excursionsData.values.visitorsNumber,
        time_from: excursionsData.values.timeFrom,
        time_to: excursionsData.values.timeTill
      };

      // Only include image_data if a new image was provided
      // This ensures we keep the existing image URL if no new image is uploaded
      if (imageData) {
        updateData.image_data = imageData;
      }

      console.log('Updating excursion with data:', updateData);

      // Update the excursion with JSON data
      const response = await axiosInstance.patch<Excursion>(
        `/excursions/${excursionsData.id}`,
        updateData,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      if (!response.data) {
        throw new Error('No data received from server');
      }

      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Error updating excursion:', err.response?.data || err.message);
      return rejectWithValue(err.response?.data || 'Failed to update excursion');
    }
  }
);

const excursionsSlice = createSlice({
  name: 'excursions',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchExcursion.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchExcursion.fulfilled, (state, action) => {
        state.excursions = action.payload as Excursion[];
        state.loading = false;
      })
      .addCase(fetchExcursionById.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.excursions = [];
      })
      .addCase(fetchExcursionById.fulfilled, (state, action) => {
        state.excursions.push(action.payload as Excursion);
        state.loading = false;
      })
      .addCase(fetchExcursionsWithPagination.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchExcursionsWithPagination.fulfilled, (state, action) => {
        state.paginatedData = action.payload as ResponseWithPagination;
        state.loading = false;
      })
      .addCase(removeExcursion.fulfilled, (state, action) => {
        state.excursions = state.excursions.filter(
          (item) => item.id !== (action.meta.arg as string)
        );
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.error = action.payload;
        state.loading = false;
      });
  }
});

export default excursionsSlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}