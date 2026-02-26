import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { Users, Package, ShoppingCart, TrendingUp, DollarSign, AlertCircle } from 'lucide-react'
import api from '../../api/axios'

export default function Dashboard() {
  const [stats, setStats] = useState({ users: 0, products: 0, orders: 0, revenue: 0, pendingOrders: 0, lowStock: 0 })
  const [recentOrders, setRecentOrders] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get('/admin/users').catch(() => ({ data: { data: [] } })),
      api.get('/admin/product/view').catch(() => ({ data: { data: [] } })),
      api.get('/admin/order/view').catch(() => ({ data: { data: [] } })),
    ]).then(([usersRes, productsRes, ordersRes]) => {
      const users = usersRes.data.data || []
      const products = productsRes.data.data || []
      const orders = ordersRes.data.data || []

      const revenue = orders.reduce((sum, o) => sum + (o.payment_status ? (o.total_amount || 0) : 0), 0)
      const pendingOrders = orders.filter(o => !o.status).length
      const lowStock = products.filter(p => (p.stock || 0) < 10).length

      setStats({
        users: users.length,
        products: products.length,
        orders: orders.length,
        revenue,
        pendingOrders,
        lowStock,
      })
      setRecentOrders(orders.slice(0, 5))
    }).catch(() => toast.error('Failed to load dashboard')).finally(() => setLoading(false))
  }, [])

  const cards = [
    { label: 'Total Users', value: stats.users, icon: Users, color: 'from-blue-500 to-blue-600', iconBg: 'bg-blue-400/20' },
    { label: 'Products', value: stats.products, icon: Package, color: 'from-emerald-500 to-emerald-600', iconBg: 'bg-emerald-400/20' },
    { label: 'Total Orders', value: stats.orders, icon: ShoppingCart, color: 'from-purple-500 to-purple-600', iconBg: 'bg-purple-400/20' },
    { label: 'Revenue', value: `₹${stats.revenue.toLocaleString()}`, icon: DollarSign, color: 'from-amber-500 to-amber-600', iconBg: 'bg-amber-400/20' },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-green-600" />
      </div>
    )
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500 text-sm mt-1">Overview of your GPU store</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-5 mb-8">
        {cards.map((card) => {
          const Icon = card.icon
          return (
            <div key={card.label} className={`bg-gradient-to-br ${card.color} rounded-xl p-5 text-white shadow-lg shadow-black/5`}>
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-white/80">{card.label}</p>
                  <p className="text-3xl font-bold mt-1">{card.value}</p>
                </div>
                <div className={`w-12 h-12 ${card.iconBg} rounded-xl flex items-center justify-center`}>
                  <Icon size={24} className="text-white" />
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Alert Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 mb-8">
        <div className="bg-white rounded-xl border border-gray-200 p-5 flex items-center gap-4">
          <div className="w-11 h-11 bg-orange-100 rounded-xl flex items-center justify-center shrink-0">
            <TrendingUp size={22} className="text-orange-600" />
          </div>
          <div>
            <p className="text-sm text-gray-500">Pending Deliveries</p>
            <p className="text-2xl font-bold text-gray-900">{stats.pendingOrders}</p>
          </div>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5 flex items-center gap-4">
          <div className="w-11 h-11 bg-red-100 rounded-xl flex items-center justify-center shrink-0">
            <AlertCircle size={22} className="text-red-600" />
          </div>
          <div>
            <p className="text-sm text-gray-500">Low Stock Items</p>
            <p className="text-2xl font-bold text-gray-900">{stats.lowStock}</p>
          </div>
        </div>
      </div>

      {/* Recent Orders */}
      <div className="bg-white rounded-xl border border-gray-200 shadow-sm">
        <div className="px-5 py-4 border-b border-gray-100">
          <h2 className="font-semibold text-gray-900">Recent Orders</h2>
        </div>
        {recentOrders.length === 0 ? (
          <p className="text-center text-gray-400 py-10">No orders yet</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-gray-500 text-xs uppercase tracking-wider">
                  <th className="px-5 py-3 font-medium">Order ID</th>
                  <th className="px-5 py-3 font-medium">Amount</th>
                  <th className="px-5 py-3 font-medium">Payment</th>
                  <th className="px-5 py-3 font-medium">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {recentOrders.map((o) => (
                  <tr key={o.id} className="hover:bg-gray-50/50">
                    <td className="px-5 py-3 font-mono text-xs text-gray-600">{o.order_id}</td>
                    <td className="px-5 py-3 font-semibold text-gray-900">₹{(o.total_amount || 0).toLocaleString()}</td>
                    <td className="px-5 py-3">
                      <span className="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-gray-100 text-gray-700 capitalize">
                        {o.payment_method}
                      </span>
                    </td>
                    <td className="px-5 py-3">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium ${
                        o.status ? 'bg-green-50 text-green-700' : 'bg-amber-50 text-amber-700'
                      }`}>
                        {o.status ? 'Delivered' : 'Pending'}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
