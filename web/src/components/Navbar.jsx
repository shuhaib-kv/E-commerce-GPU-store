import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { Home, ShoppingCart, Package, Wallet, MapPin, LogIn, LogOut, UserPlus, Cpu, Menu, X, ShoppingBag } from 'lucide-react'

export default function Navbar() {
  const { user, logoutUser } = useAuth()
  const location = useLocation()
  const [open, setOpen] = useState(false)

  const isActive = (path) => location.pathname === path

  const navLink = (to, label, Icon) => (
    <Link
      key={to}
      to={to}
      onClick={() => setOpen(false)}
      className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition ${
        isActive(to)
          ? 'bg-green-600 text-white'
          : 'text-gray-300 hover:bg-gray-800 hover:text-green-400'
      }`}
    >
      <Icon size={16} />
      {label}
    </Link>
  )

  return (
    <nav className="bg-gray-900 text-white shadow-lg sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 text-xl font-bold text-green-400 hover:text-green-300 transition">
          <Cpu size={26} />
          GPU Store
        </Link>

        {/* Desktop */}
        <div className="hidden md:flex items-center gap-2">
          {navLink('/', 'Home', Home)}
          {navLink('/products', 'Products', ShoppingBag)}
          {user ? (
            <>
              {navLink('/cart', 'Cart', ShoppingCart)}
              {navLink('/orders', 'Orders', Package)}
              {navLink('/wallet', 'Wallet', Wallet)}
              {navLink('/address', 'Address', MapPin)}
              <button
                onClick={logoutUser}
                className="flex items-center gap-1.5 bg-red-600 hover:bg-red-700 px-3 py-1.5 rounded-lg text-sm font-medium transition ml-2"
              >
                <LogOut size={16} />
                Logout
              </button>
            </>
          ) : (
            <>
              {navLink('/login', 'Login', LogIn)}
              <Link
                to="/signup"
                className="flex items-center gap-1.5 bg-green-600 hover:bg-green-700 px-3 py-1.5 rounded-lg text-sm font-medium transition"
              >
                <UserPlus size={16} />
                Sign Up
              </Link>
            </>
          )}
        </div>

        {/* Mobile toggle */}
        <button onClick={() => setOpen(!open)} className="md:hidden p-1 text-gray-300 hover:text-white">
          {open ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>

      {/* Mobile menu */}
      {open && (
        <div className="md:hidden bg-gray-900 border-t border-gray-800 px-4 pb-4 flex flex-col gap-1">
          {navLink('/', 'Home', Home)}
          {navLink('/products', 'Products', ShoppingBag)}
          {user ? (
            <>
              {navLink('/cart', 'Cart', ShoppingCart)}
              {navLink('/orders', 'Orders', Package)}
              {navLink('/wallet', 'Wallet', Wallet)}
              {navLink('/address', 'Address', MapPin)}
              <button
                onClick={() => { logoutUser(); setOpen(false) }}
                className="flex items-center gap-1.5 bg-red-600 hover:bg-red-700 px-3 py-1.5 rounded-lg text-sm font-medium transition mt-2"
              >
                <LogOut size={16} />
                Logout
              </button>
            </>
          ) : (
            <>
              {navLink('/login', 'Login', LogIn)}
              <Link
                to="/signup"
                onClick={() => setOpen(false)}
                className="flex items-center gap-1.5 bg-green-600 hover:bg-green-700 px-3 py-1.5 rounded-lg text-sm font-medium transition"
              >
                <UserPlus size={16} />
                Sign Up
              </Link>
            </>
          )}
        </div>
      )}
    </nav>
  )
}
