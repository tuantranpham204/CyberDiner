import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { authService } from '../../services/authService'
import { useAuth } from '../../store/authStore'

export default function SignInPage() {
  const [form, setForm]     = useState({ username: '', password: '' })
  const [error, setError]   = useState('')
  const [loading, setLoading] = useState(false)

  const { signIn } = useAuth()
  const navigate   = useNavigate()
  const location   = useLocation()
  const from       = location.state?.from?.pathname ?? '/'

  function handleChange(e) {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await authService.signIn(form)
      signIn(res.data.user, res.data)
      navigate(from, { replace: true })
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.wrapper}>
      <form style={styles.card} onSubmit={handleSubmit} noValidate>
        <h2 style={styles.title}>Sign in</h2>

        <div style={styles.field}>
          <label style={styles.label}>Username</label>
          <input
            name="username"
            value={form.username}
            onChange={handleChange}
            required
            autoComplete="username"
            style={styles.input}
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Password</label>
          <input
            name="password"
            type="password"
            value={form.password}
            onChange={handleChange}
            required
            autoComplete="current-password"
            style={styles.input}
          />
        </div>

        {error && <p style={styles.error}>{error}</p>}

        <button type="submit" disabled={loading} style={styles.btn}>
          {loading ? 'Signing in…' : 'Sign in'}
        </button>

        <p style={styles.footer}>
          Don&apos;t have an account?{' '}
          <Link to="/signup" style={styles.link}>Sign up</Link>
        </p>
      </form>
    </div>
  )
}

const styles = {
  wrapper: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: '#f5f5f5',
  },
  card: {
    background: '#fff',
    borderRadius: 12,
    padding: '2rem',
    width: '100%',
    maxWidth: 400,
    boxShadow: '0 4px 24px rgba(0,0,0,.08)',
    display: 'flex',
    flexDirection: 'column',
    gap: '0.75rem',
  },
  title:  { margin: '0 0 0.5rem', fontSize: '1.5rem', fontWeight: 700 },
  field:  { display: 'flex', flexDirection: 'column', gap: 4 },
  label:  { fontSize: '0.85rem', fontWeight: 500, color: '#444' },
  input: {
    padding: '0.5rem 0.75rem',
    border: '1px solid #ddd',
    borderRadius: 8,
    fontSize: '0.95rem',
    outline: 'none',
    width: '100%',
    boxSizing: 'border-box',
  },
  error:  { margin: 0, color: '#d32f2f', fontSize: '0.85rem' },
  btn: {
    marginTop: '0.5rem',
    padding: '0.65rem',
    background: '#1a1a2e',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    fontSize: '1rem',
    fontWeight: 600,
    cursor: 'pointer',
  },
  footer: { margin: 0, textAlign: 'center', fontSize: '0.875rem', color: '#666' },
  link:   { color: '#1a1a2e', fontWeight: 600 },
}
