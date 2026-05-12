import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authService } from '../../services/authService'

const GENDERS = ['Male', 'Female', 'Other']

export default function SignUpPage() {
  const [form, setForm] = useState({
    name: '', surname: '', username: '', email: '',
    password: '', confirmPassword: '', gender: '', dob: '',
  })
  const [error, setError]     = useState('')
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()

  function handleChange(e) {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')

    if (form.password !== form.confirmPassword) {
      setError('Passwords do not match')
      return
    }

    setLoading(true)
    try {
      const { name, surname, username, email, password, gender, dob } = form
      await authService.signUp({ name, surname, username, email, password, gender, dob })
      navigate('/signin')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.wrapper}>
      <form style={styles.card} onSubmit={handleSubmit} noValidate>
        <h2 style={styles.title}>Create account</h2>

        <div style={styles.row}>
          <Field label="First name" name="name"    value={form.name}    onChange={handleChange} required />
          <Field label="Last name"  name="surname" value={form.surname} onChange={handleChange} required />
        </div>

        <Field label="Username" name="username" value={form.username} onChange={handleChange} required />
        <Field label="Email"    name="email"    type="email" value={form.email} onChange={handleChange} required />
        <Field label="Password" name="password" type="password" value={form.password} onChange={handleChange} required />
        <Field label="Confirm password" name="confirmPassword" type="password"
          value={form.confirmPassword} onChange={handleChange} required />

        <div style={styles.field}>
          <label style={styles.label}>Gender</label>
          <select name="gender" value={form.gender} onChange={handleChange} style={styles.input}>
            <option value="">Select gender</option>
            {GENDERS.map((g) => (
              <option key={g} value={g.toLowerCase()}>{g}</option>
            ))}
          </select>
        </div>

        <Field label="Date of birth" name="dob" type="date" value={form.dob} onChange={handleChange} />

        {error && <p style={styles.error}>{error}</p>}

        <button type="submit" disabled={loading} style={styles.btn}>
          {loading ? 'Creating account…' : 'Sign up'}
        </button>

        <p style={styles.footer}>
          Already have an account?{' '}
          <Link to="/signin" style={styles.link}>Sign in</Link>
        </p>
      </form>
    </div>
  )
}

function Field({ label, name, type = 'text', value, onChange, required }) {
  return (
    <div style={styles.field}>
      <label style={styles.label}>{label}</label>
      <input name={name} type={type} value={value} onChange={onChange}
        required={required} style={styles.input} />
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
    maxWidth: 460,
    boxShadow: '0 4px 24px rgba(0,0,0,.08)',
    display: 'flex',
    flexDirection: 'column',
    gap: '0.75rem',
  },
  title:  { margin: '0 0 0.5rem', fontSize: '1.5rem', fontWeight: 700 },
  row:    { display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem' },
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
