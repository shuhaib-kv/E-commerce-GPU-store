import { useEffect, useState } from 'react'
import { BarChart3, Package, ShoppingCart, TrendingUp, Clock, CheckCircle } from 'lucide-react'
import api from '../../api/axios'

export default function PlatformStats() {
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/admin/manage/stats')
      .then(res => setStats(res.data.data))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return <div className="flex items-center justify-center h-64"><div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" /></div>
  }

  if (!stats) return <p className="text-gray-500 text-center py-12">Failed to load stats</p>

  const cards = [
    { label: 'Total Products', value: stats.total_products, icon: Package, color: 'bg-blue-500/10 text-blue-500' },
    { label: 'Total Orders', value: stats.total_orders, icon: ShoppingCart, color: 'bg-green-500/10 text-green-500' },
    { label: 'Total Revenue', value: `₹${(stats.total_revenue || 0).toLocaleString()}`, icon: TrendingUp, color: 'bg-amber-500/10 text-amber-500' },
    { label: 'Paid Orders', value: stats.paid_orders, icon: CheckCircle, color: 'bg-emerald-500/10 text-emerald-500' },
    { label: 'Pending Orders', value: stats.pending_orders, icon: Clock, color: 'bg-orange-500/10 text-orange-500' },
  ]

  return (
    <div>
      <div className="flex items-center gap-3 mb-8">
        <BarChart3 size={28} className="text-purple-400" />
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Platform Statistics</h1>
          <p className="text-sm text-gray-500">Overall platform performance metrics</p>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        {cards.map(card => {
          const Icon = card.icon
          return (
            <div key={card.label} className="bg-white rounded-xl border border-gray-200 p-6">
              <div className="flex items-center justify-between mb-4">
                <span className="text-sm text-gray-500">{card.label}</span>
                <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${card.color}`}>
                  <Icon size={20} />
                </div>
              </div>
              <p className="text-3xl font-bold text-gray-900">{card.value}</p>
            </div>
          )
        })}
      </div>
    </div>
  )
}
