import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useTranslation } from 'react-i18next'
import { loginUser, clearError, selectAuthLoading, selectAuthError } from '../features/auth/authSlice'
import { AppDispatch } from '../app/store'

const LoginPage: React.FC = () => {
  const { t } = useTranslation()
  const dispatch = useDispatch<AppDispatch>()
  const isLoading = useSelector(selectAuthLoading)
  const error = useSelector(selectAuthError)
  
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  })

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }))
    
    // Clear error when user starts typing
    if (error) {
      dispatch(clearError())
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.email || !formData.password) {
      return
    }

    await dispatch(loginUser(formData))
  }

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-card card">
          <h1>{t('auth.login')}</h1>
          
          {error && (
            <div className="alert alert-error">
              {error}
            </div>
          )}
          
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="email" className="form-label">
                {t('auth.email')}
              </label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleInputChange}
                className="form-control"
                required
                disabled={isLoading}
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="password" className="form-label">
                {t('auth.password')}
              </label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                className="form-control"
                required
                disabled={isLoading}
              />
            </div>
            
            <button
              type="submit"
              className="btn btn-primary"
              disabled={isLoading || !formData.email || !formData.password}
            >
              {isLoading ? (
                <span className="loading">
                  <span className="spinner"></span>
                  {t('common.loading')}
                </span>
              ) : (
                t('auth.loginButton')
              )}
            </button>
          </form>
          
          <div className="demo-credentials">
            <h3>Demo Credentials:</h3>
            <p>Email: admin@example.com</p>
            <p>Password: admin123</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage
