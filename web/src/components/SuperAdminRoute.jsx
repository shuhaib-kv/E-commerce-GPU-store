import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function SuperAdminRoute({ children }) {
  const { admin } = useAuth()
  if (!admin) return <Navigate to="/admin/login" />
  if (admin.role !== 'superadmin') return <Navigate to="/admin/dashboard" />
  return children
}
