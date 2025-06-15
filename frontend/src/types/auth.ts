export interface LoginCredentials {
  email: string;
  password: string;
}

export interface User {
  email: string;
  // Add other user properties as needed
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

export interface AuthResponse {
  token: string;
  user: User;
}
