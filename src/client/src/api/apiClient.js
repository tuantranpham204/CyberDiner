const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'

async function request(path, options = {}) {
  const token = tokenStorage.getAccess()
  const authHeader = token ? { Authorization: `Bearer ${token}` } : {}

  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...authHeader, ...options.headers },
    ...options,
  })

  const body = await res.json()

  if (!res.ok) {
    throw new Error(body.message ?? 'Request failed')
  }

  return body
}

export const apiClient = {
  get:  (path, headers = {}) =>
    request(path, { method: 'GET', headers }),
  post: (path, data, headers = {}) =>
    request(path, { method: 'POST', body: JSON.stringify(data), headers }),
  put:  (path, data, headers = {}) =>
    request(path, { method: 'PUT', body: JSON.stringify(data), headers }),
  del:  (path, headers = {}) =>
    request(path, { method: 'DELETE', headers }),
}

export const tokenStorage = {
  set: (access, refresh) => {
    localStorage.setItem('cd_access_token', access)
    localStorage.setItem('cd_refresh_token', refresh)
  },
  getAccess:  () => localStorage.getItem('cd_access_token'),
  getRefresh: () => localStorage.getItem('cd_refresh_token'),
  clear: () => {
    localStorage.removeItem('cd_access_token')
    localStorage.removeItem('cd_refresh_token')
  },
}
