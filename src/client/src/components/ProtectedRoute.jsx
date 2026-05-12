import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../store/authStore'

export default function ProtectedRoute({ children }) {
  const { isAuthenticated, isLoading } = useAuth()
  const location = useLocation()

  if (isLoading) return null

  return isAuthenticated
    ? children
    : <Navigate to="/signin" state={{ from: location }} replace />
}
