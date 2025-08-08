import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import api, { ApiResponse, PaginatedResponse } from '../../app/axios'

export interface Todo {
  id: string
  title: string
  description: string
  completed: boolean
  user_id: string
  created_at: string
  updated_at: string
}

export interface TodoState {
  todos: Todo[]
  currentTodo: Todo | null
  isLoading: boolean
  error: string | null
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

interface CreateTodoRequest {
  title: string
  description: string
}

interface UpdateTodoRequest {
  title?: string
  description?: string
  completed?: boolean
}

interface FetchTodosParams {
  page?: number
  page_size?: number
}

const initialState: TodoState = {
  todos: [],
  currentTodo: null,
  isLoading: false,
  error: null,
  pagination: {
    page: 1,
    page_size: 10,
    total: 0,
    total_pages: 0,
  },
}

// Async thunks
export const fetchTodos = createAsyncThunk<
  PaginatedResponse<Todo>,
  FetchTodosParams,
  { rejectValue: string }
>('todos/fetchTodos', async (params, { rejectWithValue }) => {
  try {
    const response = await api.get<ApiResponse<PaginatedResponse<Todo>>>('/todos', {
      params,
    })
    
    if (response.data.success && response.data.data) {
      return response.data.data
    } else {
      return rejectWithValue(response.data.error?.message || 'Failed to fetch todos')
    }
  } catch (error: any) {
    return rejectWithValue(
      error.response?.data?.error?.message || 'An error occurred while fetching todos'
    )
  }
})

export const createTodo = createAsyncThunk<
  Todo,
  CreateTodoRequest,
  { rejectValue: string }
>('todos/createTodo', async (todoData, { rejectWithValue }) => {
  try {
    const response = await api.post<ApiResponse<Todo>>('/todos', todoData)
    
    if (response.data.success && response.data.data) {
      return response.data.data
    } else {
      return rejectWithValue(response.data.error?.message || 'Failed to create todo')
    }
  } catch (error: any) {
    return rejectWithValue(
      error.response?.data?.error?.message || 'An error occurred while creating todo'
    )
  }
})

export const updateTodo = createAsyncThunk<
  Todo,
  { id: string; data: UpdateTodoRequest },
  { rejectValue: string }
>('todos/updateTodo', async ({ id, data }, { rejectWithValue }) => {
  try {
    const response = await api.put<ApiResponse<Todo>>(`/todos/${id}`, data)
    
    if (response.data.success && response.data.data) {
      return response.data.data
    } else {
      return rejectWithValue(response.data.error?.message || 'Failed to update todo')
    }
  } catch (error: any) {
    return rejectWithValue(
      error.response?.data?.error?.message || 'An error occurred while updating todo'
    )
  }
})

export const deleteTodo = createAsyncThunk<
  string,
  string,
  { rejectValue: string }
>('todos/deleteTodo', async (id, { rejectWithValue }) => {
  try {
    await api.delete(`/todos/${id}`)
    return id
  } catch (error: any) {
    return rejectWithValue(
      error.response?.data?.error?.message || 'An error occurred while deleting todo'
    )
  }
})

const todoSlice = createSlice({
  name: 'todos',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null
    },
    setCurrentTodo: (state, action: PayloadAction<Todo | null>) => {
      state.currentTodo = action.payload
    },
    toggleTodo: (state, action: PayloadAction<string>) => {
      const todo = state.todos.find(t => t.id === action.payload)
      if (todo) {
        todo.completed = !todo.completed
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch todos
      .addCase(fetchTodos.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(fetchTodos.fulfilled, (state, action) => {
        state.isLoading = false
        state.todos = action.payload.data
        state.pagination = action.payload.pagination
        state.error = null
      })
      .addCase(fetchTodos.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload || 'Failed to fetch todos'
      })
      // Create todo
      .addCase(createTodo.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(createTodo.fulfilled, (state, action) => {
        state.isLoading = false
        state.todos.unshift(action.payload)
        state.error = null
      })
      .addCase(createTodo.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload || 'Failed to create todo'
      })
      // Update todo
      .addCase(updateTodo.pending, (state) => {
        state.error = null
      })
      .addCase(updateTodo.fulfilled, (state, action) => {
        const index = state.todos.findIndex(todo => todo.id === action.payload.id)
        if (index !== -1) {
          state.todos[index] = action.payload
        }
        if (state.currentTodo?.id === action.payload.id) {
          state.currentTodo = action.payload
        }
        state.error = null
      })
      .addCase(updateTodo.rejected, (state, action) => {
        state.error = action.payload || 'Failed to update todo'
      })
      // Delete todo
      .addCase(deleteTodo.pending, (state) => {
        state.error = null
      })
      .addCase(deleteTodo.fulfilled, (state, action) => {
        state.todos = state.todos.filter(todo => todo.id !== action.payload)
        if (state.currentTodo?.id === action.payload) {
          state.currentTodo = null
        }
        state.error = null
      })
      .addCase(deleteTodo.rejected, (state, action) => {
        state.error = action.payload || 'Failed to delete todo'
      })
  },
})

export const { clearError, setCurrentTodo, toggleTodo } = todoSlice.actions

// Selectors
export const selectTodos = (state: { todos: TodoState }) => state.todos.todos
export const selectCurrentTodo = (state: { todos: TodoState }) => state.todos.currentTodo
export const selectTodosLoading = (state: { todos: TodoState }) => state.todos.isLoading
export const selectTodosError = (state: { todos: TodoState }) => state.todos.error
export const selectTodosPagination = (state: { todos: TodoState }) => state.todos.pagination

export default todoSlice.reducer
