import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Todo } from '../features/todos/todoSlice'

interface TodoFormProps {
  todo?: Todo
  onSubmit: (data: { title: string; description: string }) => void
  onCancel: () => void
  isLoading?: boolean
}

const TodoForm: React.FC<TodoFormProps> = ({ todo, onSubmit, onCancel, isLoading }) => {
  const { t } = useTranslation()
  const [formData, setFormData] = useState({
    title: todo?.title || '',
    description: todo?.description || '',
  })

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.title.trim()) return
    onSubmit(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="todo-form">
      <div className="form-group">
        <label htmlFor="title" className="form-label">
          {t('todos.todoTitle')}
        </label>
        <input
          type="text"
          id="title"
          name="title"
          value={formData.title}
          onChange={handleInputChange}
          className="form-control"
          required
          disabled={isLoading}
        />
      </div>
      
      <div className="form-group">
        <label htmlFor="description" className="form-label">
          {t('todos.todoDescription')}
        </label>
        <textarea
          id="description"
          name="description"
          value={formData.description}
          onChange={handleInputChange}
          className="form-control"
          rows={3}
          disabled={isLoading}
        />
      </div>
      
      <div className="form-actions">
        <button
          type="submit"
          className="btn btn-primary"
          disabled={isLoading || !formData.title.trim()}
        >
          {isLoading ? t('common.loading') : t('common.save')}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="btn btn-secondary"
          disabled={isLoading}
        >
          {t('common.cancel')}
        </button>
      </div>
    </form>
  )
}

export default TodoForm
