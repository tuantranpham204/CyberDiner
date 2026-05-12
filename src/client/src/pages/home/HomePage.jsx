import { useAuth } from '../../store/authStore'
import { useNavigate } from 'react-router-dom'

export default function HomePage() {
  const { user, signOut } = useAuth()
  const navigate = useNavigate()

  async function handleSignOut() {
    await signOut()
    navigate('/signin', { replace: true })
  }

  return (
    <div style={styles.wrapper}>
      <div style={styles.card}>
        <h1 style={styles.title}>Welcome to CyberDiner{user ? `, ${user.name}` : ''}!</h1>
        <p style={styles.sub}>You are signed in.</p>
        <button style={styles.btn} onClick={handleSignOut}>Sign out</button>
      </div>
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
    padding: '2.5rem',
    textAlign: 'center',
    boxShadow: '0 4px 24px rgba(0,0,0,.08)',
  },
  title: { margin: '0 0 0.5rem', fontSize: '1.75rem', fontWeight: 700 },
  sub:   { margin: '0 0 1.5rem', color: '#666' },
  btn: {
    padding: '0.6rem 1.75rem',
    background: '#1a1a2e',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    fontSize: '1rem',
    fontWeight: 600,
    cursor: 'pointer',
  },
}
