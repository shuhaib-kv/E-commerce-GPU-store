import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { user, logoutUser } = useAuth()

  return (
    <nav className="bg-gray-900 text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <Link to="/" className="text-2xl font-bold text-green-400">
          GPU Store
        </Link>
        <div className="flex items-center gap-6">
          <Link to="/products" className="hover:text-green-400 transition">Products</Link>
          {user ? (
            <>
              <Link to="/cart" className="hover:text-green-400 transition">Cart</Link>
              <Link to="/orders" className="hover:text-green-400 transition">Orders</Link>
              <Link to="/wallet" className="hover:text-green-400 transition">Wallet</Link>
              <Link to="/address" className="hover:text-green-400 transition">Address</Link>
              <button
                onClick={logoutUser}
                className="bg-red-600 hover:bg-red-700 px-3 py-1 rounded text-sm transition"
              >
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="hover:text-green-400 transition">Login</Link>
              <Link to="/signup" className="bg-green-600 hover:bg-green-700 px-3 py-1 rounded transition">
                Sign Up
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  )
}
