import { createContext, useContext, useState } from 'react'

const AuthContext = createContext()

function loadFromStorage(key) {
  try {
    const val = localStorage.getItem(key)
    return val ? JSON.parse(val) : null
  } catch {
    return null
  }
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => loadFromStorage('user'))
  const [admin, setAdmin] = useState(() => loadFromStorage('admin'))

  const loginUser = (userData) => {
    setUser(userData)
    localStorage.setItem('user', JSON.stringify(userData))
  }
  const logoutUser = () => {
    setUser(null)
    localStorage.removeItem('user')
  }
  const loginAdmin = (adminData) => {
    setAdmin(adminData)
    localStorage.setItem('admin', JSON.stringify(adminData))
  }
  const logoutAdmin = () => {
    setAdmin(null)
    localStorage.removeItem('admin')
  }

  return (
    <AuthContext.Provider value={{ user, admin, loginUser, logoutUser, loginAdmin, logoutAdmin }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
