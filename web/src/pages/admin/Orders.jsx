import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
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
      <h1 className="text-3xl font-bold mb-6">Order Management</h1>
      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full text-left">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-4 py-3">Order ID</th>
              <th className="px-4 py-3">Amount</th>
              <th className="px-4 py-3">Payment</th>
              <th className="px-4 py-3">Delivery</th>
              <th className="px-4 py-3">Payment Status</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((o) => (
              <tr key={o.id} className="border-b hover:bg-gray-50">
                <td className="px-4 py-3 font-mono text-sm">{o.order_id?.slice(0, 8)}</td>
                <td className="px-4 py-3">₹{o.total_amount}</td>
                <td className="px-4 py-3">{o.payment_method}</td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs ${o.status ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
                    {o.status ? 'Delivered' : 'Pending'}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs ${o.payment_status ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                    {o.payment_status ? 'Paid' : 'Unpaid'}
                  </span>
                </td>
                <td className="px-4 py-3 flex gap-2">
                  {!o.status && (
                    <button onClick={() => updateOrder(o.order_id, 'status', 'true')} className="text-green-600 hover:underline text-sm">
                      Mark Delivered
                    </button>
                  )}
                  {!o.payment_status && (
                    <button onClick={() => updateOrder(o.order_id, 'paymentstatus', 'true')} className="text-blue-600 hover:underline text-sm">
                      Mark Paid
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {orders.length === 0 && <p className="text-center text-gray-500 py-8">No orders found.</p>}
      </div>
    </div>
  )
}
