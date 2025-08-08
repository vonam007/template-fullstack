import React from 'react'
import { useTranslation } from 'react-i18next'
import { useSelector } from 'react-redux'
import { Link } from 'react-router-dom'
import { selectCurrentUser } from '../features/auth/authSlice'

const HomePage: React.FC = () => {
  const { t } = useTranslation()
  const user = useSelector(selectCurrentUser)

  return (
    <div className="container">
      <div className="card">
        <h1>{t('navigation.home')}</h1>
        <p>Welcome back, {user?.name}!</p>
        
        <div className="quick-actions">
          <h2>Quick Actions</h2>
          <div className="action-buttons">
            <Link to="/todos" className="btn btn-primary">
              {t('navigation.todos')}
            </Link>
          </div>
        </div>

        <div className="stats">
          <h2>Overview</h2>
          <p>This is your dashboard. Navigate to different sections using the menu above.</p>
        </div>
      </div>
    </div>
  )
}

export default HomePage
