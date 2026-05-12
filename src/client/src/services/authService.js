import { apiClient } from '../api/apiClient'

export const authService = {
  signUp: (data) => apiClient.post('/auth/signup', data),
  signIn: (data) => apiClient.post('/auth/signin', data),
  // UC-04: fire-and-forget — caller clears local state regardless of outcome
  signOut: async () => {
    try {
      await apiClient.post('/auth/signout')
    } catch {
      // Expired or already-invalid token; server-side revocation skipped.
    }
  },
}
