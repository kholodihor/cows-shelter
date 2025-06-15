import { PdfFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

export type Pdf = {
  id: string;
  title: string;
  document_url: string;
  document_id: string;
};

type PdfState = {
  documents: Pdf[];
  loading: boolean;
  error: string | null;
};

const initialState: PdfState = {
  documents: [],
  loading: false,
  error: null
};

export const fetchPdfs = createAsyncThunk('pdf/fetchPdfs', async (_, { rejectWithValue }) => {
  try {
    const response = await axiosInstance.get<Pdf[]>('/pdf');
    const data = response.data;
    // Return empty array if no documents are found
    return Array.isArray(data) ? data : [];
  } catch (error) {
    const err = error as AxiosError;
    console.error('Failed to fetch PDFs:', err.message);
    return rejectWithValue('Failed to load PDF documents');
  }
});

export const fetchPdfById = createAsyncThunk(
  'pdf/fetchPdfById',
  async (id: string) => {
    try {
      const response = await axiosInstance.get<Pdf>(`/pdf/${id}`);
      const data = response.data;
      return data;
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const removePdf = createAsyncThunk(
  'pdf/removePdf',
  async (id: string) => {
    try {
      await axiosInstance.delete(`/pdf/${id}`);
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

export const addNewPdf = createAsyncThunk(
  'pdf/addNewPdf',
  async (values: PdfFormInput) => {
    try {
      const file = values.document[0];
      const formData = new FormData();
      formData.append('file', file);
      const { data } = await axiosInstance.post('/upload-pdf', formData);
      const newPdf = {
        title: values.title,
        document_url: data.document_url,
        document_id: data.document_id
      };
      await axiosInstance.post('/pdf', newPdf);
    } catch (error) {
      const err = error as AxiosError;
      return err.message;
    }
  }
);

const pdfSlice = createSlice({
  name: 'pdf',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchPdfs.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPdfs.fulfilled, (state, action: PayloadAction<Pdf[] | string>) => {
        state.loading = false;
        if (Array.isArray(action.payload)) {
          state.documents = action.payload;
        } else {
          state.error = 'Invalid response format from server';
          state.documents = [];
        }
      })
      .addCase(fetchPdfs.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string || action.error.message || 'Failed to fetch PDFs';
        state.documents = [];
      })
      .addCase(fetchPdfById.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.documents = [];
      })
      .addCase(fetchPdfById.fulfilled, (state, action) => {
        state.documents.push(action.payload as Pdf);
        state.loading = false;
      })
      .addCase(removePdf.fulfilled, (state, action) => {
        state.documents = state.documents.filter(
          (item) => item.id !== (action.meta.arg as string)
        );
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.error = action.payload;
        state.loading = false;
      });
  }
});

export default pdfSlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}