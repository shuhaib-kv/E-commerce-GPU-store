import { useEffect, useState } from 'react'
import { Package, CheckCircle, Clock, CreditCard, XCircle, Truck } from 'lucide-react'
import toast from 'react-hot-toast'
import api from '../../api/axios'

export default function Orders() {
  const [orders, setOrders] = useState([])

  useEffect(() => {
    api.get('/user/orderview')
      .then((res) => setOrders(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load orders'))
  }, [])

  const cancelOrder = async (orderId) => {
    if (!confirm('Cancel this order?')) return
    try {
      await api.post('/user/cancel/order', { orderid: orderId })
      setOrders((prev) => prev.filter((o) => o.order_id !== orderId))
      toast.success('Order cancelled')
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to cancel')
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <div className="bg-green-50 text-green-600 p-2 rounded-lg">
          <Package size={22} />
        </div>
        <h1 className="text-3xl font-bold text-gray-900">My Orders</h1>
      </div>

      {orders.length === 0 ? (
        <div className="text-center py-20">
          <Package size={64} className="text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 text-lg">No orders yet.</p>
        </div>
      ) : (
        <div className="flex flex-col gap-4">
          {orders.map((order) => (
            <div key={order.id} className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <p className="font-semibold text-gray-900 flex items-center gap-2">
                    <Package size={16} className="text-gray-500" />
                    Order #{order.order_id?.slice(0, 8)}
                  </p>
                  <p className="text-sm text-gray-500 mt-1">
                    {new Date(order.created_at).toLocaleDateString()}
                  </p>
                </div>
                <span className="text-green-600 font-bold text-xl">₹{order.total_amount?.toLocaleString()}</span>
              </div>

              <div className="flex flex-wrap gap-2 mb-3">
                <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${
                  order.status ? 'bg-green-50 text-green-700' : 'bg-amber-50 text-amber-700'
                }`}>
                  {order.status ? <CheckCircle size={12} /> : <Clock size={12} />}
                  {order.status ? 'Delivered' : 'Pending'}
                </span>
                <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${
                  order.payment_status ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
                }`}>
                  <CreditCard size={12} />
                  {order.payment_status ? 'Paid' : 'Unpaid'}
                </span>
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-gray-50 text-gray-600">
                  <Truck size={12} />
                  {order.payment_method}
                </span>
              </div>

              {!order.status && (
                <button
                  onClick={() => cancelOrder(order.order_id)}
                  className="flex items-center gap-1.5 bg-red-50 hover:bg-red-100 text-red-600 px-4 py-1.5 rounded-lg text-sm font-medium transition"
                >
                  <XCircle size={14} />
                  Cancel Order
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
