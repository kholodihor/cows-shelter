import { ContactsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

export type Contact = {
  id: string;
  email: string;
  phone: string;
};

type ContactState = {
  contacts: Contact[];
  loading: boolean;
  error: string | null;
};

const initialState: ContactState = {
  contacts: [],
  loading: false,
  error: null
};

export const fetchContacts = createAsyncThunk(
  'contacts/fetchContacts',
  async (_, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.get<Contact[]>('/contacts');
      const data = response.data;
      return Array.isArray(data) ? data : [];
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to fetch contacts:', err.message);
      return rejectWithValue('Failed to load contacts');
    }
  }
);

export const addContacts = createAsyncThunk(
  'contacts/addContacts',
  async (values: ContactsFormInput, { rejectWithValue }) => {
    try {
      const newData = {
        email: values.email,
        phone: values.phone
      };
      const response = await axiosInstance.post<Contact>('/contacts', newData);
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to add contact:', err.message);
      return rejectWithValue('Failed to add contact');
    }
  }
);

export const editEmail = createAsyncThunk(
  'contacts/editEmail',
  async (data: { id?: string; values: ContactsFormInput }, { rejectWithValue }) => {
    try {
      const newData = {
        email: data.values.email
      };
      const response = await axiosInstance.patch<Contact>(`/contacts/${data.id}`, newData);
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to update email:', err.message);
      return rejectWithValue('Failed to update email');
    }
  }
);

export const editPhone = createAsyncThunk(
  'contacts/editPhone',
  async (data: { id?: string; values: ContactsFormInput }, { rejectWithValue }) => {
    try {
      const newData = {
        phone: data.values.phone
      };
      const response = await axiosInstance.patch<Contact>(`/contacts/${data.id}`, newData);
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to update phone:', err.message);
      return rejectWithValue('Failed to update phone');
    }
  }
);

const contactsSlice = createSlice({
  name: 'contacts',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchContacts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchContacts.fulfilled, (state, action) => {
        state.contacts = action.payload as Contact[];
        state.loading = false;
      })
      .addMatcher(isError, (state, action: PayloadAction<string>) => {
        state.error = action.payload;
        state.loading = false;
      });
  }
});

export default contactsSlice.reducer;

function isError(action: AnyAction) {
  return action.type.endsWith('rejected');
}