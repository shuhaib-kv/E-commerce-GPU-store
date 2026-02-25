import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const links = [
  { to: '/admin/dashboard', label: 'Dashboard' },
  { to: '/admin/users', label: 'Users' },
  { to: '/admin/products', label: 'Products' },
  { to: '/admin/categories', label: 'Categories' },
  { to: '/admin/orders', label: 'Orders' },
  { to: '/admin/coupons', label: 'Coupons' },
  { to: '/admin/discounts', label: 'Discounts' },
]

export default function AdminSidebar() {
  const location = useLocation()
  const { logoutAdmin } = useAuth()

  return (
    <aside className="w-64 bg-gray-900 text-white min-h-screen p-4 flex flex-col">
      <h2 className="text-xl font-bold text-green-400 mb-8">Admin Panel</h2>
      <nav className="flex flex-col gap-1 flex-1">
        {links.map((link) => (
          <Link
            key={link.to}
            to={link.to}
            className={`px-4 py-2 rounded transition ${
              location.pathname === link.to
                ? 'bg-green-600 text-white'
                : 'hover:bg-gray-800'
            }`}
          >
            {link.label}
          </Link>
        ))}
      </nav>
      <button
        onClick={logoutAdmin}
        className="mt-4 bg-red-600 hover:bg-red-700 px-4 py-2 rounded transition"
      >
        Logout
      </button>
    </aside>
  )
}
