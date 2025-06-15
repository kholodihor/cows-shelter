import { PartnersFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

export type Partner = {
  id: string;
  name: string;
  link: string;
  logo: string;
};

type ResponseWithPagination = {
  partners: Partner[];
  totalLength: number;
};

type PartnersState = {
  partners: Partner[];
  loading: boolean;
  error: string | null;
  paginatedData: ResponseWithPagination;
};

const initialState: PartnersState = {
  partners: [],
  loading: false,
  error: null,
  paginatedData: {
    partners: [],
    totalLength: 0
  }
};

export const fetchPartners = createAsyncThunk(
  'partners/fetchPartners',
  async (_, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Partner[]>('/partners');
      const data = response.data;
      return Array.isArray(data) ? data : [];
    } catch (error) {
      const err = error as Error;
      console.error('Failed to fetch partners:', err.message);
      return rejectWithValue('Failed to load partners');
    }
  }
);

export const fetchPartnersWithPagination = createAsyncThunk(
  'partners/fetchPartnersWithPagination',
  async (query: { page: number; limit: number }, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<ResponseWithPagination>(
        `/partners/pagination?page=${query.page}&limit=${query.limit}`
      );
      const data = response.data;
      return data || { partners: [], totalLength: 0 };
    } catch (error) {
      const err = error as Error;
      console.error('Failed to fetch partners with pagination:', err.message);
      return rejectWithValue('Failed to load partners with pagination');
    }
  }
);

export const removePartner = createAsyncThunk(
  'partners/removePartner',
  async (id: string, { rejectWithValue }) => {
    try {
      await axiosInstance.delete(`/partners/${id}`);
      return id; // Return the deleted ID for state updates
    } catch (error) {
      const err = error as Error;
      console.error('Failed to delete partner:', err.message);
      return rejectWithValue('Failed to delete partner');
    }
  }
);

export const addNewPartner = createAsyncThunk(
  'partners/addNewPartner',
  async (values: PartnersFormInput, { rejectWithValue }) => {
    try {
      const file = values.logo[0];
      const formData = new FormData();
      formData.append('image', file);
      const { data: uploadData } = await axiosInstance.post<{ image_url: string }>('/upload-image', formData);
      
      const newPartner = {
        name: values.name,
        link: values.link,
        logo: uploadData.image_url
      };
      
      const response = await axiosInstance.post<Partner>('/partners', newPartner);
      return response.data;
    } catch (error) {
      const err = error as Error;
      console.error('Failed to add partner:', err.message);
      return rejectWithValue('Failed to add partner');
    }
  }
);

export const editPartner = createAsyncThunk(
  'partners/editPartner',
  async (partnersData: { id?: string; values: PartnersFormInput }, { rejectWithValue }) => {
    try {
      let logoUrl: string;
      
      if (partnersData.values.logo[0].size > 0) {
        const file = partnersData.values.logo[0];
        const formData = new FormData();
        formData.append('image', file);
        const { data } = await axiosInstance.post<{ image_url: string }>('/upload-image', formData);
        logoUrl = data.image_url;
      } else {
        logoUrl = partnersData.values.logo[0].name;
      }
      
      const updatedPartner = {
        name: partnersData.values.name,
        link: partnersData.values.link,
        logo: logoUrl
      };
      
      const response = await axiosInstance.patch<Partner>(`/partners/${partnersData.id}`, updatedPartner);
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to update partner:', err.message);
      return rejectWithValue('Failed to update partner');
    }
  }
);

const partnersSlice = createSlice({
  name: 'partners',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchPartners.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPartners.fulfilled, (state, action) => {
        state.partners = action.payload as Partner[];
        state.loading = false;
      })
      .addCase(fetchPartnersWithPagination.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPartnersWithPagination.fulfilled, (state, action) => {
        state.paginatedData = action.payload as ResponseWithPagination;
        state.loading = false;
      })
      .addCase(removePartner.fulfilled, (state, action) => {
        state.partners = state.partners.filter(
          (item) => item.id !== (action.meta.arg as string)
        );
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.error = action.payload;
        state.loading = false;
      });
  }
});

export default partnersSlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}