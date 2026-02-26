import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import api from '../../api/axios'

export default function Coupons() {
  const [coupons, setCoupons] = useState([])
  const [form, setForm] = useState({ couponname: '', couponpercentage: '', expiresat: '' })
  const [showForm, setShowForm] = useState(false)

  const fetchCoupons = () => {
    api.get('/admin/list/coupons')
      .then((res) => setCoupons(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load coupons'))
  }

  useEffect(() => { fetchCoupons() }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await api.post('/admin/add/coupon', {
        couponname: form.couponname,
        couponpercentage: parseInt(form.couponpercentage),
        expiresat: parseInt(form.expiresat),
      })
      toast.success('Coupon added')
      setForm({ couponname: '', couponpercentage: '', expiresat: '' })
      setShowForm(false)
      fetchCoupons()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to add coupon')
    }
  }

  const deleteCoupon = async (couponName) => {
    if (!confirm('Delete this coupon?')) return
    try {
      await api.delete('/admin/delete/coupon', { data: { couponname: couponName } })
      toast.success('Coupon deleted')
      fetchCoupons()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to delete coupon')
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Coupons</h1>
        <button onClick={() => setShowForm(!showForm)} className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded transition">
          {showForm ? 'Cancel' : 'Add Coupon'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow mb-8 grid grid-cols-1 md:grid-cols-2 gap-4">
          <input placeholder="Coupon Name" value={form.couponname} onChange={(e) => setForm({ ...form, couponname: e.target.value })} className="border rounded px-3 py-2" required />
          <input type="number" placeholder="Discount %" value={form.couponpercentage} onChange={(e) => setForm({ ...form, couponpercentage: e.target.value })} className="border rounded px-3 py-2" required />
          <input type="number" placeholder="Expires in (days)" value={form.expiresat} onChange={(e) => setForm({ ...form, expiresat: e.target.value })} className="border rounded px-3 py-2" required />
          <button className="bg-green-600 hover:bg-green-700 text-white py-2 rounded transition font-semibold md:col-span-2">Add Coupon</button>
        </form>
      )}

      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full text-left">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Code</th>
              <th className="px-4 py-3">Discount</th>
              <th className="px-4 py-3">Expiry</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {coupons.map((c) => (
              <tr key={c.id} className="border-b hover:bg-gray-50">
                <td className="px-4 py-3">{c.coupon_name}</td>
                <td className="px-4 py-3 font-mono">{c.coupon_code}</td>
                <td className="px-4 py-3">{c.coupon_percentage}%</td>
                <td className="px-4 py-3">{new Date(c.expiry_date).toLocaleDateString()}</td>
                <td className="px-4 py-3">
                  <button onClick={() => deleteCoupon(c.coupon_name)} className="text-red-600 hover:underline text-sm">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {coupons.length === 0 && <p className="text-center text-gray-500 py-8">No coupons yet.</p>}
      </div>
    </div>
  )
}
