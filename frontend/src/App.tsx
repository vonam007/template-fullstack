import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { selectIsAuthenticated } from './features/auth/authSlice'
import { Navbar } from './components'
import { LoginPage, TodosPage, HomePage } from './pages'
import './App.css'

const App: React.FC = () => {
  const isAuthenticated = useSelector(selectIsAuthenticated)

  return (
    <div className="App">
      {isAuthenticated && <Navbar />}
      <main className="main-content">
        <Routes>
          <Route 
            path="/login" 
            element={!isAuthenticated ? <LoginPage /> : <Navigate to="/" />} 
          />
          <Route 
            path="/" 
            element={isAuthenticated ? <HomePage /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/todos" 
            element={isAuthenticated ? <TodosPage /> : <Navigate to="/login" />} 
          />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </main>
    </div>
  )
}

export default App
