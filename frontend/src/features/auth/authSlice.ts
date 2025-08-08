import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import api, { ApiResponse } from '../../app/axios'

export interface User {
  id: string
  email: string
  name: string
  created_at: string
  updated_at: string
}

export interface AuthState {
  user: User | null
  token: string | null
  isLoading: boolean
  error: string | null
  isAuthenticated: boolean
}

interface LoginRequest {
  email: string
  password: string
}

interface LoginResponse {
  token: string
  user: User
}

const initialState: AuthState = {
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  token: localStorage.getItem('authToken'),
  isLoading: false,
  error: null,
  isAuthenticated: !!localStorage.getItem('authToken'),
}

// Async thunks
export const loginUser = createAsyncThunk<
  LoginResponse,
  LoginRequest,
  { rejectValue: string }
>('auth/login', async (credentials, { rejectWithValue }) => {
  try {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', credentials)
    
    if (response.data.success && response.data.data) {
      const { token, user } = response.data.data
      localStorage.setItem('authToken', token)
      localStorage.setItem('user', JSON.stringify(user))
      return response.data.data
    } else {
      return rejectWithValue(response.data.error?.message || 'Login failed')
    }
  } catch (error: any) {
    return rejectWithValue(
      error.response?.data?.error?.message || 'An error occurred during login'
    )
  }
})

export const logoutUser = createAsyncThunk('auth/logout', async () => {
  localStorage.removeItem('authToken')
  localStorage.removeItem('user')
})

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null
    },
    setCredentials: (state, action: PayloadAction<{ user: User; token: string }>) => {
      state.user = action.payload.user
      state.token = action.payload.token
      state.isAuthenticated = true
      localStorage.setItem('authToken', action.payload.token)
      localStorage.setItem('user', JSON.stringify(action.payload.user))
    },
  },
  extraReducers: (builder) => {
    builder
      // Login
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.isLoading = false
        state.user = action.payload.user
        state.token = action.payload.token
        state.isAuthenticated = true
        state.error = null
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload || 'Login failed'
        state.isAuthenticated = false
      })
      // Logout
      .addCase(logoutUser.fulfilled, (state) => {
        state.user = null
        state.token = null
        state.isAuthenticated = false
        state.error = null
      })
  },
})

export const { clearError, setCredentials } = authSlice.actions

// Selectors
export const selectCurrentUser = (state: { auth: AuthState }) => state.auth.user
export const selectIsAuthenticated = (state: { auth: AuthState }) => state.auth.isAuthenticated
export const selectAuthLoading = (state: { auth: AuthState }) => state.auth.isLoading
export const selectAuthError = (state: { auth: AuthState }) => state.auth.error

export default authSlice.reducer
