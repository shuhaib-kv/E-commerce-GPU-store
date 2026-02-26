import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ShoppingCart, CreditCard, MapPin, Tag } from 'lucide-react'
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
      .catch(() => {})
    api.get('/user/address')
      .then((res) => {
        const addr = res.data.data
        setAddresses(addr ? [addr] : [])
      })
      .catch(() => {})
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
      <div className="flex items-center gap-3 mb-8">
        <div className="bg-green-50 text-green-600 p-2 rounded-lg">
          <ShoppingCart size={22} />
        </div>
        <h1 className="text-3xl font-bold text-gray-900">Shopping Cart</h1>
        {items.length > 0 && (
          <span className="bg-green-100 text-green-700 text-sm font-semibold px-2.5 py-0.5 rounded-full">
            {items.length} item{items.length > 1 ? 's' : ''}
          </span>
        )}
      </div>

      {items.length === 0 ? (
        <div className="text-center py-20">
          <ShoppingCart size={64} className="text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 text-lg">Your cart is empty.</p>
        </div>
      ) : (
        <>
          <div className="flex flex-col gap-3 mb-8">
            {items.map((item) => (
              <CartItem key={item.id} item={item} />
            ))}
          </div>

          <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
            <div className="flex justify-between items-center text-xl font-bold mb-6 pb-4 border-b border-gray-100">
              <span className="text-gray-900">Total</span>
              <span className="text-green-600">₹{total.toLocaleString()}</span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                  <MapPin size={14} /> Delivery Address
                </label>
                <select
                  value={addressId}
                  onChange={(e) => setAddressId(e.target.value)}
                  className="border border-gray-200 rounded-lg px-3 py-2 w-full text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
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
                <label className="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                  <CreditCard size={14} /> Payment Method
                </label>
                <select
                  value={paymentMethod}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  className="border border-gray-200 rounded-lg px-3 py-2 w-full text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                >
                  <option value="razorpay">Razorpay</option>
                  <option value="wallet">Wallet</option>
                </select>
              </div>
            </div>

            <div className="mb-6">
              <label className="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-1.5">
                <Tag size={14} /> Coupon Code
              </label>
              <input
                value={coupon}
                onChange={(e) => setCoupon(e.target.value)}
                placeholder="Enter coupon code (optional)"
                className="border border-gray-200 rounded-lg px-3 py-2 w-full text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>

            <button
              onClick={handleOrder}
              className="w-full bg-green-600 hover:bg-green-700 text-white py-3 rounded-xl text-lg font-semibold transition shadow-lg shadow-green-600/25"
            >
              Place Order
            </button>
          </div>
        </>
      )}
    </div>
  )
}
