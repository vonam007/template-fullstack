import React from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { logoutUser, selectCurrentUser } from '../features/auth/authSlice'
import LanguageSelector from './LanguageSelector'
import { AppDispatch } from '../app/store'

const Navbar: React.FC = () => {
  const { t } = useTranslation()
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const user = useSelector(selectCurrentUser)

  const handleLogout = async () => {
    await dispatch(logoutUser())
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-brand">
          Fullstack Template
        </Link>
        
        <div className="navbar-menu">
          <Link to="/" className="navbar-item">
            {t('navigation.home')}
          </Link>
          <Link to="/todos" className="navbar-item">
            {t('navigation.todos')}
          </Link>
        </div>

        <div className="navbar-actions">
          <LanguageSelector />
          
          <div className="user-menu">
            <span className="user-name">
              {t('common.hello')}, {user?.name}
            </span>
            <button onClick={handleLogout} className="btn btn-secondary">
              {t('auth.logout')}
            </button>
          </div>
        </div>
      </div>
    </nav>
  )
}

export default Navbar
