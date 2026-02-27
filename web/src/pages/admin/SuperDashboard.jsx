import { useEffect, useState } from 'react'
import { Users, Package, ShoppingCart, TrendingUp, ShieldCheck } from 'lucide-react'
import api from '../../api/axios'

export default function SuperDashboard() {
  const [admins, setAdmins] = useState([])
  const [stats, setStats] = useState({ total_products: 0, total_orders: 0, total_revenue: 0, paid_orders: 0, pending_orders: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get('/admin/manage/admins').catch(() => ({ data: { data: [] } })),
      api.get('/admin/manage/stats').catch(() => ({ data: { data: {} } })),
    ]).then(([adminsRes, statsRes]) => {
      setAdmins(adminsRes.data.data || [])
      setStats(statsRes.data.data || {})
      setLoading(false)
    })
  }, [])

  if (loading) {
    return <div className="flex items-center justify-center h-64"><div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" /></div>
  }

  const cards = [
    { label: 'Total Admins', value: admins.length, icon: Users, color: 'purple' },
    { label: 'Total Products', value: stats.total_products, icon: Package, color: 'blue' },
    { label: 'Total Orders', value: stats.total_orders, icon: ShoppingCart, color: 'green' },
    { label: 'Total Revenue', value: `₹${(stats.total_revenue || 0).toLocaleString()}`, icon: TrendingUp, color: 'amber' },
  ]

  const colorMap = {
    purple: 'bg-purple-500/10 text-purple-400',
    blue: 'bg-blue-500/10 text-blue-400',
    green: 'bg-green-500/10 text-green-400',
    amber: 'bg-amber-500/10 text-amber-400',
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-8">
        <ShieldCheck size={28} className="text-purple-400" />
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Super Admin Dashboard</h1>
          <p className="text-sm text-gray-500">Platform overview and admin management</p>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 mb-8">
        {cards.map((card) => {
          const Icon = card.icon
          return (
            <div key={card.label} className="bg-white rounded-xl border border-gray-200 p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm text-gray-500">{card.label}</span>
                <div className={`w-9 h-9 rounded-lg flex items-center justify-center ${colorMap[card.color]}`}>
                  <Icon size={18} />
                </div>
              </div>
              <p className="text-2xl font-bold text-gray-900">{card.value}</p>
            </div>
          )
        })}
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Admins Overview</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500 border-b border-gray-100">
                <th className="pb-3 font-medium">Name</th>
                <th className="pb-3 font-medium">Email</th>
                <th className="pb-3 font-medium">Store Name</th>
                <th className="pb-3 font-medium">Status</th>
                <th className="pb-3 font-medium">Payment Gateway</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {admins.map((admin) => (
                <tr key={admin.id} className="text-gray-700">
                  <td className="py-3 font-medium">{admin.name}</td>
                  <td className="py-3">{admin.email}</td>
                  <td className="py-3">{admin.store_name || '-'}</td>
                  <td className="py-3">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${admin.block_status ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                      {admin.block_status ? 'Blocked' : 'Active'}
                    </span>
                  </td>
                  <td className="py-3">{admin.payment_gateway?.provider || 'Not configured'}</td>
                </tr>
              ))}
              {admins.length === 0 && (
                <tr><td colSpan={5} className="py-8 text-center text-gray-400">No admins found</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
