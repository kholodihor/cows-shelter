import { PdfFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';
import { transformMinioUrlsInData } from '@/utils/minioUrlHelper';

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

export const fetchPdfs = createAsyncThunk(
  'pdf/fetchPdfs',
  async (_, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Pdf[]>('/pdf');
      let data = response.data;
      // Return empty array if no documents are found
      data = Array.isArray(data) ? data : [];
      // Transform MinIO URLs in the response data
      return transformMinioUrlsInData(data);
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch PDFs:', err.message);
      return rejectWithValue('Failed to load PDF documents');
    }
  }
);

export const fetchPdfById = createAsyncThunk(
  'pdf/fetchPdfById',
  async (id: string) => {
    try {
      const response = await axiosInstance.get<Pdf>(`/pdf/${id}`);
      // Transform MinIO URL in the response data
      return transformMinioUrlsInData(response.data);
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

// Helper function to convert file to base64
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });
};

export const addNewPdf = createAsyncThunk(
  'pdf/addNewPdf',
  async (values: PdfFormInput, { rejectWithValue }) => {
    try {
      const file = values.document[0];
      const documentData = await fileToBase64(file);

      const newPdf = {
        title: values.title,
        document_data: documentData
      };

      const response = await axiosInstance.post<Pdf>('/pdf', newPdf);
      // Transform MinIO URL in the response data
      return transformMinioUrlsInData(response.data);
    } catch (error) {
      const err = error as Error;
      console.error('Failed to add new PDF:', err.message);
      return rejectWithValue('Failed to add new PDF');
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
      .addCase(
        fetchPdfs.fulfilled,
        (state, action: PayloadAction<Pdf[] | string>) => {
          state.loading = false;
          if (Array.isArray(action.payload)) {
            state.documents = action.payload;
          } else {
            state.error = 'Invalid response format from server';
            state.documents = [];
          }
        }
      )
      .addCase(fetchPdfs.rejected, (state, action) => {
        state.loading = false;
        state.error =
          (action.payload as string) ||
          action.error.message ||
          'Failed to fetch PDFs';
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
