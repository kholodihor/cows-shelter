import { store } from '../store';
import { login as loginAction, logout as logoutAction, clearError } from '../store/slices/authSlice';
import { LoginCredentials } from '../types/auth';

/**
 * Authentication service that works with Redux
 */

/**
 * Login user and update Redux state
 */
export const login = async (credentials: LoginCredentials) => {
  try {
    // Clear any previous errors
    store.dispatch(clearError());
    
    // Dispatch login action which will handle the API call
    const result = await store.dispatch(loginAction(credentials));
    
    // If login was rejected, throw an error
    if (loginAction.rejected.match(result)) {
      throw new Error(result.payload as string || 'Login failed');
    }
    
    return result.payload;
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
};

/**
 * Logout user and clear Redux state
 */
export const logout = async () => {
  try {
    await store.dispatch(logoutAction());
  } catch (error) {
    console.error('Logout error:', error);
  }
};

/**
 * Get current authentication state from Redux
 */
export const getAuthState = () => {
  return store.getState().auth;
};

/**
 * Check if user is authenticated
 */
export const isAuthenticated = (): boolean => {
  const { isAuthenticated } = store.getState().auth;
  return isAuthenticated;
};

/**
 * Get auth headers for API requests
 */
export const getAuthHeader = () => {
  const { token } = store.getState().auth;
  return token ? { Authorization: `Bearer ${token}` } : {};
};
