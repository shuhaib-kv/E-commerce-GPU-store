import { useEffect, useState } from 'react'
import api from '../../api/axios'

export default function Coupons() {
  const [coupons, setCoupons] = useState([])
  const [form, setForm] = useState({ coupon_name: '', coupon_code: '', coupon_percentage: '', expiry_date: '' })
  const [showForm, setShowForm] = useState(false)

  const fetchCoupons = () => {
    api.get('/admin/list/coupons')
      .then((res) => setCoupons(res.data.coupons || []))
      .catch(() => {})
  }

  useEffect(() => { fetchCoupons() }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await api.post('/admin/add/coupon', {
        ...form,
        coupon_percentage: parseInt(form.coupon_percentage),
      })
      setForm({ coupon_name: '', coupon_code: '', coupon_percentage: '', expiry_date: '' })
      setShowForm(false)
      fetchCoupons()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to add coupon')
    }
  }

  const deleteCoupon = async (id) => {
    if (!confirm('Delete this coupon?')) return
    await api.delete('/admin/delete/coupon', { data: { coupon_id: id } })
    fetchCoupons()
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
          <input placeholder="Coupon Name" value={form.coupon_name} onChange={(e) => setForm({ ...form, coupon_name: e.target.value })} className="border rounded px-3 py-2" required />
          <input placeholder="Coupon Code" value={form.coupon_code} onChange={(e) => setForm({ ...form, coupon_code: e.target.value })} className="border rounded px-3 py-2" required />
          <input type="number" placeholder="Discount %" value={form.coupon_percentage} onChange={(e) => setForm({ ...form, coupon_percentage: e.target.value })} className="border rounded px-3 py-2" required />
          <input type="date" value={form.expiry_date} onChange={(e) => setForm({ ...form, expiry_date: e.target.value })} className="border rounded px-3 py-2" required />
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
              <tr key={c._id} className="border-b hover:bg-gray-50">
                <td className="px-4 py-3">{c.coupon_name}</td>
                <td className="px-4 py-3 font-mono">{c.coupon_code}</td>
                <td className="px-4 py-3">{c.coupon_percentage}%</td>
                <td className="px-4 py-3">{new Date(c.expiry_date).toLocaleDateString()}</td>
                <td className="px-4 py-3">
                  <button onClick={() => deleteCoupon(c._id)} className="text-red-600 hover:underline text-sm">Delete</button>
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
