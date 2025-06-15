import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { AxiosError } from 'axios';
import axiosInstance from '../../utils/axios';
import { AuthState, LoginCredentials, AuthResponse, User } from '../../types/auth';

// Helper to get initial state from localStorage
const getInitialState = (): AuthState => {
  const token = localStorage.getItem('token');
  const userString = localStorage.getItem('user');
  let user: User | null = null;

  if (userString) {
    try {
      user = JSON.parse(userString);
    } catch (_e) {
      // Invalid JSON in localStorage, clear it
      localStorage.removeItem('user');
      // No need to handle the error further as we're just clearing invalid data
    }
  }

  return {
    user,
    token,
    isAuthenticated: !!token,
    loading: false,
    error: null
  };
};

// Async thunks
export const login = createAsyncThunk(
  'auth/login',
  async (credentials: LoginCredentials, { rejectWithValue }) => {
    try {
      const response = await axiosInstance.post<AuthResponse>('/login', credentials, {
        withCredentials: true
      });

      if (!response.data.token) {
        return rejectWithValue('No token received');
      }

      return response.data;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Login failed:', err.message);
      // Handle the error response data safely
      const errorMessage = err.response?.data && typeof err.response.data === 'object'
        ? (err.response.data as any).message || 'Login failed'
        : 'Login failed. Please check your credentials.';
      return rejectWithValue(errorMessage);
    }
  }
);

export const logout = createAsyncThunk(
  'auth/logout',
  async (_, { rejectWithValue }) => {
    try {
      // If you have a logout endpoint, call it here
      // await axiosInstance.post('/api/logout');

      return true;
    } catch (error) {
      const err = error as AxiosError;
      console.error('Logout failed:', err.message);
      return rejectWithValue('Logout failed');
    }
  }
);

// Create the slice
const authSlice = createSlice({
  name: 'auth',
  initialState: getInitialState(),
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;
      localStorage.setItem('user', JSON.stringify(action.payload));
    },
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload;
      state.isAuthenticated = true;
      localStorage.setItem('token', action.payload);
    }
  },
  extraReducers: (builder) => {
    // Login
    builder.addCase(login.pending, (state) => {
      state.loading = true;
      state.error = null;
    });
    builder.addCase(login.fulfilled, (state, action) => {
      state.loading = false;
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.isAuthenticated = true;

      // Save to localStorage
      localStorage.setItem('token', action.payload.token);
      localStorage.setItem('user', JSON.stringify(action.payload.user));
    });
    builder.addCase(login.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload as string;
      state.isAuthenticated = false;
    });

    // Logout
    builder.addCase(logout.pending, (state) => {
      state.loading = true;
    });
    builder.addCase(logout.fulfilled, (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.loading = false;

      // Clear localStorage
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    });
    builder.addCase(logout.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload as string;
      // Even if API logout fails, we clear the local state
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    });
  }
});

export const { clearError, setUser, setToken } = authSlice.actions;
export default authSlice.reducer;
