import React from 'react'
import { useTranslation } from 'react-i18next'
import { Todo } from '../features/todos/todoSlice'

interface TodoListProps {
  todos: Todo[]
  isLoading: boolean
  onEdit: (todo: Todo) => void
  onDelete: (id: string) => void
  onToggleComplete: (id: string, completed: boolean) => void
}

const TodoList: React.FC<TodoListProps> = ({ 
  todos, 
  isLoading, 
  onEdit, 
  onDelete, 
  onToggleComplete 
}) => {
  const { t } = useTranslation()

  if (isLoading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
        {t('common.loading')}
      </div>
    )
  }

  if (todos.length === 0) {
    return (
      <div className="empty-state card">
        <h3>{t('todos.noTodos')}</h3>
        <p>Create your first todo to get started!</p>
      </div>
    )
  }

  return (
    <div className="todo-list">
      {todos.map(todo => (
        <div key={todo.id} className={`todo-item card ${todo.completed ? 'completed' : ''}`}>
          <div className="todo-content">
            <div className="todo-header">
              <h3 className="todo-title">{todo.title}</h3>
              <span className={`todo-status ${todo.completed ? 'completed' : 'pending'}`}>
                {todo.completed ? t('todos.completed') : t('todos.pending')}
              </span>
            </div>
            {todo.description && (
              <p className="todo-description">{todo.description}</p>
            )}
            <div className="todo-meta">
              <small>Created: {new Date(todo.created_at).toLocaleDateString()}</small>
              {todo.updated_at !== todo.created_at && (
                <small>Updated: {new Date(todo.updated_at).toLocaleDateString()}</small>
              )}
            </div>
          </div>
          
          <div className="todo-actions">
            <button
              onClick={() => onToggleComplete(todo.id, todo.completed)}
              className={`btn ${todo.completed ? 'btn-secondary' : 'btn-primary'}`}
              title={todo.completed ? 'Mark as pending' : 'Mark as complete'}
            >
              {todo.completed ? '↩️' : '✅'}
            </button>
            
            <button
              onClick={() => onEdit(todo)}
              className="btn btn-secondary"
              title={t('common.edit')}
            >
              ✏️
            </button>
            
            <button
              onClick={() => onDelete(todo.id)}
              className="btn btn-danger"
              title={t('common.delete')}
            >
              🗑️
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}

export default TodoList
