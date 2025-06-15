import { ExcursionsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

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

export const fetchExcursion = createAsyncThunk(
  'excursions/fetchExcursion',
  async () => {
    try {
      const response = await axiosInstance.get<Excursion[]>('/excursions');
      const data = response.data;
      return data;
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
      return data;
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
      console.log(data);
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
  async (values: ExcursionsFormInput) => {
    try {
      const file = values.image[0];
      const formData = new FormData();
      formData.append('image', file);
      const { data } = await axiosInstance.post('/upload-image', formData);
      const newPost = {
        title_ua: values.titleUa,
        title_en: values.titleEn,
        description_ua: values.descriptionUa,
        description_en: values.descriptionEn,
        amount_of_persons: values.visitorsNumber,
        time_from: values.timeFrom,
        time_to: values.timeTill,
        image_url: data.image_url
      };
      await axiosInstance.post('/excursions', newPost);
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const editExcursion = createAsyncThunk(
  'excursions/editExcursion',
  async (excursionsData: { id?: string; values: ExcursionsFormInput }) => {
    try {
      if (excursionsData.values.image[0].size > 0) {
        const file = excursionsData.values.image[0];
        const formData = new FormData();
        formData.append('image', file);
        const { data } = await axiosInstance.post('/upload-image', formData);
        const newExcursion = {
          title_ua: excursionsData.values.titleUa,
          title_en: excursionsData.values.titleEn,
          description_ua: excursionsData.values.descriptionUa,
          description_en: excursionsData.values.descriptionEn,
          amount_of_persons: excursionsData.values.visitorsNumber,
          time_from: excursionsData.values.timeFrom,
          time_to: excursionsData.values.timeTill,
          image_url: data.image_url
        };
        await axiosInstance.patch(`/excursions/${excursionsData.id}`, newExcursion);
      } else {
        const newExcursion = {
          title_ua: excursionsData.values.titleUa,
          title_en: excursionsData.values.titleEn,
          description_ua: excursionsData.values.descriptionUa,
          description_en: excursionsData.values.descriptionEn,
          amount_of_persons: excursionsData.values.visitorsNumber,
          time_from: excursionsData.values.timeFrom,
          time_to: excursionsData.values.timeTill,
          image_url: excursionsData.values.image[0].name
        };
        await axiosInstance.patch(`/excursions/${excursionsData.id}`, newExcursion);
      }
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
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