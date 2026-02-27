import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import {
  LayoutDashboard, Users, ShieldCheck, BarChart3,
  LogOut, Cpu
} from 'lucide-react'

const links = [
  { to: '/admin/super/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/admin/super/admins', label: 'Manage Admins', icon: Users },
  { to: '/admin/super/stats', label: 'Platform Stats', icon: BarChart3 },
]

export default function SuperAdminSidebar() {
  const location = useLocation()
  const { logoutAdmin } = useAuth()

  return (
    <aside className="w-64 bg-gray-950 text-gray-300 min-h-screen flex flex-col border-r border-gray-800">
      <div className="px-5 py-6 border-b border-gray-800">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 bg-purple-600 rounded-lg flex items-center justify-center">
            <ShieldCheck size={20} className="text-white" />
          </div>
          <div>
            <h2 className="text-white font-bold text-base leading-tight">GPU Store</h2>
            <p className="text-xs text-purple-400">Super Admin</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 px-3 py-4 space-y-1">
        <p className="px-3 text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-2">Menu</p>
        {links.map((link) => {
          const Icon = link.icon
          const active = location.pathname === link.to
          return (
            <Link
              key={link.to}
              to={link.to}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all ${
                active
                  ? 'bg-purple-600/15 text-purple-400'
                  : 'text-gray-400 hover:bg-gray-800/60 hover:text-gray-200'
              }`}
            >
              <Icon size={18} strokeWidth={active ? 2.2 : 1.8} />
              {link.label}
            </Link>
          )
        })}
      </nav>

      <div className="px-3 pb-4 border-t border-gray-800 pt-4">
        <button
          onClick={logoutAdmin}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-red-400 hover:bg-red-500/10 transition"
        >
          <LogOut size={18} strokeWidth={1.8} />
          Logout
        </button>
      </div>
    </aside>
  )
}
