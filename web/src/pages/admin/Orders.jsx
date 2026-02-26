import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { ShoppingCart, CheckCircle2, CreditCard } from 'lucide-react'
import api from '../../api/axios'

export default function AdminOrders() {
  const [orders, setOrders] = useState([])

  const fetchOrders = () => {
    api.get('/admin/order/view')
      .then((res) => setOrders(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load orders'))
  }

  useEffect(() => { fetchOrders() }, [])

  const updateOrder = async (orderId, field, value) => {
    try {
      await api.patch('/admin/order/view', { orderid: orderId, [field]: value })
      toast.success('Order updated')
      fetchOrders()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Update failed')
    }
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <ShoppingCart size={24} className="text-purple-600" />
          Order Management
        </h1>
        <p className="text-gray-500 text-sm mt-1">{orders.length} total orders</p>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-gray-500 text-xs uppercase tracking-wider bg-gray-50/80">
              <th className="px-5 py-3.5 font-medium">Order ID</th>
              <th className="px-5 py-3.5 font-medium">Amount</th>
              <th className="px-5 py-3.5 font-medium">Payment</th>
              <th className="px-5 py-3.5 font-medium">Delivery</th>
              <th className="px-5 py-3.5 font-medium">Payment Status</th>
              <th className="px-5 py-3.5 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {orders.map((o) => (
              <tr key={o.id} className="hover:bg-gray-50/50 transition">
                <td className="px-5 py-3.5 font-mono text-xs text-gray-600">{o.order_id}</td>
                <td className="px-5 py-3.5 font-semibold text-gray-900">₹{(o.total_amount || 0).toLocaleString()}</td>
                <td className="px-5 py-3.5">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-medium bg-gray-100 text-gray-700 capitalize">
                    {o.payment_method}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    o.status ? 'bg-green-50 text-green-700' : 'bg-amber-50 text-amber-700'
                  }`}>
                    {o.status ? 'Delivered' : 'Pending'}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    o.payment_status ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
                  }`}>
                    {o.payment_status ? 'Paid' : 'Unpaid'}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <div className="flex items-center justify-end gap-1">
                    {!o.status && (
                      <button
                        onClick={() => updateOrder(o.order_id, 'status', 'true')}
                        className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-green-50 text-green-700 hover:bg-green-100 transition"
                      >
                        <CheckCircle2 size={14} /> Deliver
                      </button>
                    )}
                    {!o.payment_status && (
                      <button
                        onClick={() => updateOrder(o.order_id, 'paymentstatus', 'true')}
                        className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-blue-50 text-blue-700 hover:bg-blue-100 transition"
                      >
                        <CreditCard size={14} /> Paid
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {orders.length === 0 && <p className="text-center text-gray-400 py-10">No orders found.</p>}
      </div>
    </div>
  )
}
