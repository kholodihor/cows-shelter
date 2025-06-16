import { ReviewsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

export type Review = {
  id: string;
  name_ua: string;
  name_en: string;
  review_ua: string;
  review_en: string;
};

type ResponseWithPagination = {
  reviews: Review[];
  totalLength: number;
};

type ReviewsState = {
  reviews: Review[];
  loading: boolean;
  error: string | null;
  paginatedData: ResponseWithPagination;
};

const initialState: ReviewsState = {
  reviews: [],
  loading: false,
  error: null,
  paginatedData: {
    reviews: [],
    totalLength: 0
  }
};

export const fetchReviews = createAsyncThunk(
  'reviews/fetchReviews',
  async (_, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Review[]>('/reviews');
      const data = response.data;
      return Array.isArray(data) ? data : [];
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch reviews:', err.message);
      return rejectWithValue('Failed to load reviews');
    }
  }
);

export const fetchReviewById = createAsyncThunk(
  'reviews/fetchReviewById',
  async (id: string, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Review>(`/reviews/${id}`);
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error(`Failed to fetch review with id ${id}:`, err.message);
      return rejectWithValue('Failed to load review');
    }
  }
);

export const fetchReviewsWithPagination = createAsyncThunk(
  'reviews/fetchReviewsWithPagination',
  async (query: { page: number; limit: number }, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<ResponseWithPagination>(
        `/reviews/pagination?page=${query.page}&limit=${query.limit}`
      );
      return response.data || { reviews: [], totalLength: 0 };
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch paginated reviews:', err.message);
      return rejectWithValue('Failed to load paginated reviews');
    }
  }
);

export const removeReview = createAsyncThunk(
  'reviews/removeReview',
  async (id: string, { rejectWithValue }) => {
    try {
      await axiosInstance.delete(`/reviews/${id}`);
      return id; // Return the deleted ID for state updates
    } catch (error) {
      const err = error as AxiosError;
      console.error(`Failed to delete review with id ${id}:`, err.message);
      return rejectWithValue('Failed to delete review');
    }
  }
);

export const addNewReview = createAsyncThunk(
  'reviews/addNewReview',
  async (values: ReviewsFormInput, { rejectWithValue }) => {
    try {
      const newReview = {
        name_ua: values.nameUa,
        name_en: values.nameEn,
        review_ua: values.reviewUa,
        review_en: values.reviewEn
      };
      const response = await axiosInstance.post<Review>('/reviews', newReview);
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to add new review:', err.message);
      return rejectWithValue('Failed to add new review');
    }
  }
);

export const editReview = createAsyncThunk(
  'reviews/editReview',
  async (
    reviewData: { id?: string; values: ReviewsFormInput },
    { rejectWithValue }
  ) => {
    try {
      const updatedReview = {
        name_ua: reviewData.values.nameUa,
        name_en: reviewData.values.nameEn,
        review_ua: reviewData.values.reviewUa,
        review_en: reviewData.values.reviewEn
      };
      const response = await axiosInstance.patch<Review>(
        `/reviews/${reviewData.id}`,
        updatedReview
      );
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error(
        `Failed to update review with id ${reviewData.id}:`,
        err.message
      );
      return rejectWithValue('Failed to update review');
    }
  }
);

const reviewsSlice = createSlice({
  name: 'reviews',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchReviews.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchReviews.fulfilled, (state, action) => {
        state.reviews = action.payload as Review[];
        state.loading = false;
      })
      .addCase(fetchReviewById.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.reviews = [];
      })
      .addCase(fetchReviewById.fulfilled, (state, action) => {
        state.reviews.push(action.payload as Review);
        state.loading = false;
      })
      .addCase(fetchReviewsWithPagination.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchReviewsWithPagination.fulfilled, (state, action) => {
        state.loading = false;
        state.paginatedData = action.payload as ResponseWithPagination;
      })
      .addCase(removeReview.fulfilled, (state, action) => {
        state.reviews = state.reviews.filter(
          (item) => item.id !== (action.meta.arg as string)
        );
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.error = action.payload;
        state.loading = false;
      });
  }
});

export default reviewsSlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}
