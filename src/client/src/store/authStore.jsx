import { createContext, useCallback, useContext, useEffect, useReducer } from 'react'
import { tokenStorage } from '../api/apiClient'
import { authService } from '../services/authService'

// --- State ---

const initialState = { user: null, isAuthenticated: false, isLoading: true }

function reducer(state, action) {
  switch (action.type) {
    case 'RESTORE':
      return { ...state, user: action.user, isAuthenticated: !!action.user, isLoading: false }
    case 'SIGN_IN':
      return { user: action.user, isAuthenticated: true, isLoading: false }
    case 'SIGN_OUT':
      return { user: null, isAuthenticated: false, isLoading: false }
    default:
      return state
  }
}

// --- Context ---

const AuthContext = createContext(null)

// --- Provider ---

export function AuthProvider({ children }) {
  const [state, dispatch] = useReducer(reducer, initialState)

  // UC-02: Restore session from localStorage on mount
  useEffect(() => {
    const token = tokenStorage.getAccess()
    const stored = localStorage.getItem('cd_user')
    if (token && stored) {
      try {
        dispatch({ type: 'RESTORE', user: JSON.parse(stored) })
      } catch {
        tokenStorage.clear()
        dispatch({ type: 'RESTORE', user: null })
      }
    } else {
      dispatch({ type: 'RESTORE', user: null })
    }
  }, [])

  const signIn = useCallback((user, tokens) => {
    tokenStorage.set(tokens.access_token, tokens.refresh_token)
    localStorage.setItem('cd_user', JSON.stringify(user))
    dispatch({ type: 'SIGN_IN', user })
  }, [])

  // UC-04: Always clears local state even if the server call fails (expired session)
  const signOut = useCallback(async () => {
    await authService.signOut()
    tokenStorage.clear()
    localStorage.removeItem('cd_user')
    dispatch({ type: 'SIGN_OUT' })
  }, [])

  return (
    <AuthContext.Provider value={{ ...state, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  )
}

// --- Hook ---

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
