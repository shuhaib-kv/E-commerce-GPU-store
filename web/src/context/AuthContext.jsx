import { createContext, useContext, useState } from 'react'

const AuthContext = createContext()

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [admin, setAdmin] = useState(null)

  const loginUser = (userData) => setUser(userData)
  const logoutUser = () => setUser(null)
  const loginAdmin = (adminData) => setAdmin(adminData)
  const logoutAdmin = () => setAdmin(null)

  return (
    <AuthContext.Provider value={{ user, admin, loginUser, logoutUser, loginAdmin, logoutAdmin }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
