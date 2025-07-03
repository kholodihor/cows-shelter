import { PartnersFormInput } from '@/types';
import { 
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction,
  ActionReducerMapBuilder
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

export interface PartnersState {
  partners: Partner[];
  paginatedData: ResponseWithPagination;
  loading: boolean;
  error: string | null;
}

// Helper type for error handling in async thunks
type RejectWithValue = (value: string) => any;

const initialState: PartnersState = {
  partners: [],
  loading: false,
  error: null,
  paginatedData: {
    partners: [],
    totalLength: 0
  }
};

export const fetchPartners = createAsyncThunk<Partner[], void, { rejectValue: string }>(
  'partners/fetchPartners',
  async (_, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
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

export const fetchPartnersWithPagination = createAsyncThunk<ResponseWithPagination, { page: number; limit: number }, { rejectValue: string }>(
  'partners/fetchPartnersWithPagination',
  async (query: { page: number; limit: number }, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
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

export const removePartner = createAsyncThunk<string, string, { rejectValue: string }>(
  'partners/removePartner',
  async (id: string, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
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

export const addNewPartner = createAsyncThunk<Partner, PartnersFormInput, { rejectValue: string }>(
  'partners/addNewPartner',
  async (values: PartnersFormInput, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
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
export const editPartner = createAsyncThunk<Partner, { id: string; values: PartnersFormInput }, { rejectValue: string }>(
  'partners/updatePartner',
  async ({ id, values }: { id: string; values: PartnersFormInput }, { rejectWithValue }: { rejectWithValue: RejectWithValue }) => {
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
  extraReducers: (builder: ActionReducerMapBuilder<PartnersState>) => {
    builder
      .addCase(fetchPartners.pending, (state: PartnersState) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPartners.fulfilled, (state: PartnersState, action: PayloadAction<Partner[]>) => {
        state.partners = action.payload;
        state.loading = false;
      })
      .addCase(fetchPartnersWithPagination.pending, (state: PartnersState) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchPartnersWithPagination.fulfilled, (state: PartnersState, action: PayloadAction<ResponseWithPagination>) => {
        state.paginatedData = action.payload;
        state.loading = false;
      })
      .addCase(removePartner.fulfilled, (state: PartnersState, action: PayloadAction<string, string, { arg: string }>) => {
        state.partners = state.partners.filter(
          (item: Partner) => item.ID !== action.meta.arg
        );
      })
      .addCase(editPartner.pending, (state: PartnersState) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(editPartner.fulfilled, (state: PartnersState, action: PayloadAction<Partner>) => {
        state.loading = false;
        const index = state.partners.findIndex(
          (partner: Partner) => partner.ID === action.payload.ID
        );
        if (index !== -1) {
          state.partners[index] = action.payload;
        }
        // Also update in paginated data if it exists
        const paginatedIndex = state.paginatedData.partners.findIndex(
          (partner: Partner) => partner.ID === action.payload.ID
        );
        if (paginatedIndex !== -1) {
          state.paginatedData.partners[paginatedIndex] = action.payload;
        }
      })
      .addCase(editPartner.rejected, (state: PartnersState, action: PayloadAction<string | undefined>) => {
        state.loading = false;
        state.error = action.payload || 'Failed to update partner';
        // Don't modify the partners array on error, just update the loading and error states
      })
      .addMatcher(
        (action: AnyAction) => action.type.endsWith('rejected'),
        (state: PartnersState, action: PayloadAction<string | undefined>) => {
          state.error = action.payload || 'An error occurred';
          state.loading = false;
        }
      );
  }
});

export default partnersSlice.reducer;
