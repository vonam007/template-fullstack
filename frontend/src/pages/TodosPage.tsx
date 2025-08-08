import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useTranslation } from 'react-i18next'
import {
  fetchTodos,
  createTodo,
  updateTodo,
  deleteTodo,
  selectTodos,
  selectTodosLoading,
  selectTodosError,
  selectTodosPagination,
  clearError,
  Todo
} from '../features/todos/todoSlice'
import { AppDispatch } from '../app/store'
import { TodoForm, TodoList } from '../components'

const TodosPage: React.FC = () => {
  const { t } = useTranslation()
  const dispatch = useDispatch<AppDispatch>()
  const todos = useSelector(selectTodos)
  const isLoading = useSelector(selectTodosLoading)
  const error = useSelector(selectTodosError)
  const pagination = useSelector(selectTodosPagination)
  
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null)

  useEffect(() => {
    dispatch(fetchTodos({ page: 1, page_size: 10 }))
  }, [dispatch])

  const handleCreateTodo = async (data: { title: string; description: string }) => {
    await dispatch(createTodo(data))
    setShowCreateForm(false)
  }

  const handleUpdateTodo = async (id: string, data: Partial<Todo>) => {
    await dispatch(updateTodo({ id, data }))
    setEditingTodo(null)
  }

  const handleDeleteTodo = async (id: string) => {
    if (window.confirm(t('todos.deleteConfirm'))) {
      await dispatch(deleteTodo(id))
    }
  }

  const handlePageChange = (page: number) => {
    dispatch(fetchTodos({ page, page_size: pagination.page_size }))
  }

  return (
    <div className="container">
      <div className="todos-header">
        <h1>{t('todos.title')}</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="btn btn-primary"
        >
          {t('todos.createTodo')}
        </button>
      </div>

      {error && (
        <div className="alert alert-error">
          {error}
          <button onClick={() => dispatch(clearError())}>×</button>
        </div>
      )}

      {showCreateForm && (
        <div className="card">
          <h2>{t('todos.createTodo')}</h2>
          <TodoForm
            onSubmit={handleCreateTodo}
            onCancel={() => setShowCreateForm(false)}
            isLoading={isLoading}
          />
        </div>
      )}

      {editingTodo && (
        <div className="card">
          <h2>{t('todos.editTodo')}</h2>
          <TodoForm
            todo={editingTodo}
            onSubmit={(data: { title: string; description: string }) => handleUpdateTodo(editingTodo.id, data)}
            onCancel={() => setEditingTodo(null)}
            isLoading={isLoading}
          />
        </div>
      )}

      <TodoList
        todos={todos}
        isLoading={isLoading}
        onEdit={setEditingTodo}
        onDelete={handleDeleteTodo}
        onToggleComplete={(id: string, completed: boolean) => 
          handleUpdateTodo(id, { completed: !completed })
        }
      />

      {pagination.total_pages > 1 && (
        <div className="pagination">
          {Array.from({ length: pagination.total_pages }, (_, i) => i + 1).map(page => (
            <button
              key={page}
              onClick={() => handlePageChange(page)}
              className={`btn ${page === pagination.page ? 'btn-primary' : 'btn-secondary'}`}
            >
              {page}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}

export default TodosPage
