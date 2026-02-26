import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import api from '../../api/axios'

export default function Dashboard() {
  const [stats, setStats] = useState({ users: 0, products: 0, orders: 0 })

  useEffect(() => {
    Promise.all([
      api.get('/admin/users').catch(() => ({ data: { total: 0 } })),
      api.get('/admin/product/view').catch(() => ({ data: { totalItems: 0 } })),
      api.get('/admin/order/view').catch(() => ({ data: { totalItems: 0 } })),
    ]).then(([users, products, orders]) => {
      setStats({
        users: users.data.total || users.data.data?.length || 0,
        products: products.data.totalItems || products.data.data?.length || 0,
        orders: orders.data.totalItems || orders.data.data?.length || 0,
      })
    }).catch(() => toast.error('Failed to load dashboard stats'))
  }, [])

  const cards = [
    { label: 'Total Users', value: stats.users, color: 'bg-blue-600' },
    { label: 'Total Products', value: stats.products, color: 'bg-green-600' },
    { label: 'Total Orders', value: stats.orders, color: 'bg-purple-600' },
  ]

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {cards.map((card) => (
          <div key={card.label} className={`${card.color} text-white p-6 rounded-lg shadow-lg`}>
            <p className="text-sm opacity-80">{card.label}</p>
            <p className="text-4xl font-bold mt-2">{card.value}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
