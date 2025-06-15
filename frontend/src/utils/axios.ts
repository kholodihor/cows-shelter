import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { store } from '@/store';
import { openAlert } from '@/store/slices/responseAlertSlice';

// Use the Vite proxy in development
const instance = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  withCredentials: true,
  timeout: 10000,
});

// Request interceptor
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Don't set Content-Type for FormData, let the browser set it with the correct boundary
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type'];
    } else if (!config.headers['Content-Type']) {
      config.headers['Content-Type'] = 'application/json';
    }

    // Log request details
    console.log(`Request: ${config.method?.toUpperCase()} ${config.baseURL}${config.url}`, {
      headers: config.headers,
      data: config.data
    });

    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Response interceptor
instance.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    if (error.response) {
      const { status, data } = error.response;
      let message = 'An error occurred';

      if (typeof data === 'object' && data !== null && 'message' in data) {
        message = (data as any).message;
      }

      // Show error alert
      store.dispatch(openAlert(message || 'An error occurred while processing your request.'));

      // Handle specific status codes
      if (status === 401) {
        // Handle unauthorized (e.g., redirect to login)
        localStorage.removeItem('token');
        window.location.href = '/login';
      }
    } else if (error.request) {
      // The request was made but no response was received
      store.dispatch(openAlert('No response from server. Please check your connection.'));
    } else {
      // Something happened in setting up the request
      store.dispatch(openAlert('Error setting up the request. Please try again.'));
    }
    return Promise.reject(error);
  }
);

export default instance;
