import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import api from '../../api/axios'
import CartItem from '../../components/CartItem'

export default function Cart() {
  const [items, setItems] = useState([])
  const [total, setTotal] = useState(0)
  const [paymentMethod, setPaymentMethod] = useState('razorpay')
  const [addressId, setAddressId] = useState('')
  const [addresses, setAddresses] = useState([])
  const [coupon, setCoupon] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    api.get('/user/cart/view')
      .then((res) => {
        setItems(res.data.data || [])
        setTotal(res.data.total_amount || 0)
      })
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load cart'))
    api.get('/user/address')
      .then((res) => {
        const addr = res.data.data
        setAddresses(addr ? [addr] : [])
      })
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load addresses'))
  }, [])

  const handleOrder = async () => {
    if (!addressId) return toast.error('Please select an address')
    try {
      const body = { address_id: addressId, payment_method: paymentMethod }
      if (coupon) body.coupon_code = coupon
      const res = await api.post('/cart/order', body)
      if (paymentMethod === 'razorpay') {
        window.location.href = `/razorpay`
      } else {
        toast.success('Order placed successfully!')
        navigate('/orders')
      }
    } catch (err) {
      toast.error(err.response?.data?.message || 'Order failed')
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Shopping Cart</h1>

      {items.length === 0 ? (
        <p className="text-gray-500 text-center py-12">Your cart is empty.</p>
      ) : (
        <>
          <div className="flex flex-col gap-4 mb-8">
            {items.map((item) => (
              <CartItem key={item.id} item={item} />
            ))}
          </div>

          <div className="bg-white p-6 rounded-lg shadow">
            <div className="flex justify-between text-xl font-bold mb-6">
              <span>Total</span>
              <span className="text-green-600">₹{total}</span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium mb-1">Delivery Address</label>
                <select
                  value={addressId}
                  onChange={(e) => setAddressId(e.target.value)}
                  className="border rounded px-3 py-2 w-full"
                >
                  <option value="">Select address</option>
                  {addresses.map((a) => (
                    <option key={a.id} value={a.id}>
                      {a.name} - {a.house}, {a.city} ({a.pincode})
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Payment Method</label>
                <select
                  value={paymentMethod}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  className="border rounded px-3 py-2 w-full"
                >
                  <option value="razorpay">Razorpay</option>
                  <option value="wallet">Wallet</option>
                </select>
              </div>
            </div>

            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Coupon Code (optional)</label>
              <input
                value={coupon}
                onChange={(e) => setCoupon(e.target.value)}
                placeholder="Enter coupon code"
                className="border rounded px-3 py-2 w-full"
              />
            </div>

            <button
              onClick={handleOrder}
              className="w-full bg-green-600 hover:bg-green-700 text-white py-3 rounded-lg text-lg font-semibold transition"
            >
              Place Order
            </button>
          </div>
        </>
      )}
    </div>
  )
}
