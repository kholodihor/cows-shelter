import { PartnersFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import axiosInstance from '@/utils/axios';
import { transformMinioUrlsInData } from '@/utils/minioUrlHelper';

export type Partner = {
  ID: string; // Note: Uppercase ID to match backend response
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
      let data = response.data;
      data = Array.isArray(data) ? data : [];
      // Transform MinIO URLs in the response data
      return transformMinioUrlsInData(data);
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
      const data = response.data || { partners: [], totalLength: 0 };
      // Transform MinIO URLs in the response data
      if (data.partners && Array.isArray(data.partners)) {
        data.partners = transformMinioUrlsInData(data.partners);
      }
      return data;
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

// Helper function to extract base64 data from a data URL
const extractBase64Data = (dataUrl: string): string => {
  // Handle data URL format: data:image/png;base64,iVBORw0KGgo...
  const parts = dataUrl.split(',');
  if (parts.length !== 2 || !parts[0].includes('base64')) {
    throw new Error('Invalid data URL format');
  }
  // Return just the base64-encoded part
  return parts[1];
};

// Helper function to convert file to base64 data URL
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });
};

export const addNewPartner = createAsyncThunk(
  'partners/addNewPartner',
  async (values: PartnersFormInput, { rejectWithValue }) => {
    try {
      const file = values.logo[0];
      const dataUrl = await fileToBase64(file);
      // Extract just the base64 data part
      const logoData = extractBase64Data(dataUrl);

      const newPartner = {
        name: values.name,
        link: values.link,
        logo_data: logoData
      };

      const response = await axiosInstance.post<Partner>(
        '/partners',
        newPartner
      );
      return response.data;
    } catch (error) {
      const err = error as Error;
      console.error('Failed to add new partner:', err.message);
      return rejectWithValue('Failed to add new partner');
    }
  }
);

// Alias updatePartner as editPartner for backward compatibility
export const editPartner = createAsyncThunk(
  'partners/updatePartner',
  async (
    { id, values }: { id: string; values: PartnersFormInput },
    { rejectWithValue }
  ) => {
    try {
      // Convert logo to base64 if a new logo is provided
      let logoData = '';

      if (values.logo?.[0] && values.logo[0] instanceof File) {
        // Check if this is a real file with actual content (not a dummy file)
        // Dummy files created in the edit form have size 0
        if (values.logo[0].size > 0) {
          try {
            const file = values.logo[0];
            // Check if file size is reasonable (e.g., less than 5MB)
            const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
            if (file.size > MAX_FILE_SIZE) {
              return rejectWithValue('Logo size should be less than 5MB');
            }

            const dataUrl = await fileToBase64(file);
            // Extract just the base64 data part
            logoData = extractBase64Data(dataUrl);
          } catch (uploadError) {
            console.error('Logo processing failed:', uploadError);
            return rejectWithValue('Failed to process logo');
          }
        }
      }

      // Prepare the update data
      const updateData: Record<string, any> = {
        name: values.name || '',
        link: values.link || ''
      };

      // Only include logo_data if a new logo was provided
      // This ensures we keep the existing logo URL if no new logo is uploaded
      if (logoData) {
        updateData.logo_data = logoData;
      }

      console.log('Updating partner with data:', updateData);

      // Update the partner with JSON data
      const response = await axiosInstance.patch<Partner>(
        `/partners/${id}`,
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
      const err = error as Error;
      console.error('Failed to update partner:', error);
      return rejectWithValue(err.message || 'Failed to update partner');
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
          (item) => item.ID !== (action.meta.arg as string)
        );
      })
      .addCase(editPartner.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        editPartner.fulfilled,
        (state, action: PayloadAction<Partner>) => {
          state.loading = false;
          const index = state.partners.findIndex(
            (partner) => partner.ID === action.payload.ID
          );
          if (index !== -1) {
            state.partners[index] = action.payload;
          }
          // Also update in paginated data if it exists
          const paginatedIndex = state.paginatedData.partners.findIndex(
            (partner) => partner.ID === action.payload.ID
          );
          if (paginatedIndex !== -1) {
            state.paginatedData.partners[paginatedIndex] = action.payload;
          }
        }
      )
      .addCase(editPartner.rejected, (state, action) => {
        state.loading = false;
        state.error = (action.payload as string) || 'Failed to update partner';
        // Don't modify the partners array on error, just update the loading and error states
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
