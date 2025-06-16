import { ContactsFormInput } from '@/types';
import {
  AnyAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '@/utils/axios';

export interface Contact {
  id: number;
  name: string;
  email: string;
  phone: string;
}

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
        name: values.name,
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
  async (
    data: { id?: string; values: Pick<ContactsFormInput, 'email'> },
    { rejectWithValue }
  ) => {
    try {
      const newData = {
        email: data.values.email
      };
      const response = await axiosInstance.patch<Contact>(
        `/contacts/${data.id}`,
        newData
      );
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
  async (
    data: { id?: string; values: Pick<ContactsFormInput, 'phone'> },
    { rejectWithValue }
  ) => {
    try {
      const newData = {
        phone: data.values.phone
      };
      const response = await axiosInstance.patch<Contact>(
        `/contacts/${data.id}`,
        newData
      );
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to update phone:', err.message);
      return rejectWithValue('Failed to update phone');
    }
  }
);

export const editName = createAsyncThunk(
  'contacts/editName',
  async (
    data: { id?: string; values: Pick<ContactsFormInput, 'name'> },
    { rejectWithValue }
  ) => {
    try {
      const newData = {
        name: data.values.name
      };
      const response = await axiosInstance.patch<Contact>(
        `/contacts/${data.id}`,
        newData
      );
      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Failed to update name:', err.message);
      return rejectWithValue('Failed to update name');
    }
  }
);

const contactsSlice = createSlice({
  name: 'contacts',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      // Handle fetchContacts
      .addCase(fetchContacts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchContacts.fulfilled, (state, action) => {
        state.contacts = action.payload as Contact[];
        state.loading = false;
      })
      // Handle addContacts
      .addCase(addContacts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(addContacts.fulfilled, (state, action) => {
        // If contacts is empty, add the new contact, otherwise update the first one
        if (state.contacts.length === 0) {
          state.contacts = [action.payload];
        } else {
          state.contacts[0] = action.payload;
        }
        state.loading = false;
      })
      // Handle editEmail
      .addCase(editEmail.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(editEmail.fulfilled, (state, action) => {
        if (state.contacts.length > 0) {
          state.contacts[0].email = action.payload.email;
        }
        state.loading = false;
      })
      // Handle editPhone
      .addCase(editPhone.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(editPhone.fulfilled, (state, action) => {
        if (state.contacts.length > 0) {
          state.contacts[0].phone = action.payload.phone;
        }
        state.loading = false;
      })
      // Handle editName
      .addCase(editName.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(editName.fulfilled, (state, action) => {
        if (state.contacts.length > 0) {
          state.contacts[0].name = action.payload.name;
        }
        state.loading = false;
      })
      // Handle all rejected actions
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
