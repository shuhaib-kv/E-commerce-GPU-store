import { useEffect, useState } from 'react'
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
      <h1 className="text-3xl font-bold mb-6">My Orders</h1>

      {orders.length === 0 ? (
        <p className="text-gray-500 text-center py-12">No orders yet.</p>
      ) : (
        <div className="flex flex-col gap-4">
          {orders.map((order) => (
            <div key={order.id} className="bg-white p-6 rounded-lg shadow">
              <div className="flex justify-between items-start mb-3">
                <div>
                  <p className="font-semibold">Order #{order.order_id?.slice(0, 8)}</p>
                  <p className="text-sm text-gray-500">
                    {new Date(order.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="text-right">
                  <span className="text-green-600 font-bold text-lg">₹{order.total_amount}</span>
                  <p className="text-sm">
                    <span className={`px-2 py-0.5 rounded text-xs ${order.status ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
                      {order.status ? 'Delivered' : 'Pending'}
                    </span>
                    {' '}
                    <span className={`px-2 py-0.5 rounded text-xs ${order.payment_status ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                      {order.payment_status ? 'Paid' : 'Unpaid'}
                    </span>
                  </p>
                </div>
              </div>
              <p className="text-sm text-gray-600">Payment: {order.payment_method}</p>
              {!order.status && (
                <button
                  onClick={() => cancelOrder(order.order_id)}
                  className="mt-3 bg-red-500 hover:bg-red-600 text-white px-4 py-1 rounded text-sm transition"
                >
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
