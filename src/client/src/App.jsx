import { Navigate, Route, Routes } from 'react-router-dom'
import ProtectedRoute from './components/ProtectedRoute'
import SignInPage from './pages/auth/SignInPage'
import SignUpPage from './pages/auth/SignUpPage'
import HomePage from './pages/home/HomePage'

export default function App() {
  return (
    <Routes>
      <Route path="/signin"  element={<SignInPage />} />
      <Route path="/signup"  element={<SignUpPage />} />
      <Route path="/"        element={<ProtectedRoute><HomePage /></ProtectedRoute>} />
      <Route path="*"        element={<Navigate to="/" replace />} />
    </Routes>
  )
}
