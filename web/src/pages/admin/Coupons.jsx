import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { Ticket, Plus, X, Trash2, Clock } from 'lucide-react'
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

  const isExpired = (date) => new Date(date) < new Date()

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <Ticket size={24} className="text-amber-600" />
            Coupons
          </h1>
          <p className="text-gray-500 text-sm mt-1">{coupons.length} coupons</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition ${
            showForm ? 'bg-gray-100 text-gray-700 hover:bg-gray-200' : 'bg-green-600 text-white hover:bg-green-700 shadow-sm'
          }`}
        >
          {showForm ? <><X size={16} /> Cancel</> : <><Plus size={16} /> Add Coupon</>}
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-6 mb-6">
          <h3 className="font-semibold text-gray-900 mb-4">New Coupon</h3>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider mb-1.5">Coupon Name</label>
              <input value={form.couponname} onChange={(e) => setForm({ ...form, couponname: e.target.value })} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent" required />
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider mb-1.5">Discount %</label>
              <input type="number" value={form.couponpercentage} onChange={(e) => setForm({ ...form, couponpercentage: e.target.value })} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent" required />
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider mb-1.5">Expires in (days)</label>
              <input type="number" value={form.expiresat} onChange={(e) => setForm({ ...form, expiresat: e.target.value })} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent" required />
            </div>
            <div className="md:col-span-3">
              <button className="w-full bg-green-600 hover:bg-green-700 text-white py-2.5 rounded-lg transition font-semibold text-sm shadow-sm">Add Coupon</button>
            </div>
          </form>
        </div>
      )}

      <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-gray-500 text-xs uppercase tracking-wider bg-gray-50/80">
              <th className="px-5 py-3.5 font-medium">Name</th>
              <th className="px-5 py-3.5 font-medium">Code</th>
              <th className="px-5 py-3.5 font-medium">Discount</th>
              <th className="px-5 py-3.5 font-medium">Expiry</th>
              <th className="px-5 py-3.5 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {coupons.map((c) => {
              const expired = isExpired(c.expiry_date)
              return (
                <tr key={c.id} className="hover:bg-gray-50/50 transition">
                  <td className="px-5 py-3.5 font-medium text-gray-900">{c.coupon_name}</td>
                  <td className="px-5 py-3.5">
                    <code className="bg-gray-100 text-gray-700 px-2 py-0.5 rounded text-xs font-mono">{c.coupon_code}</code>
                  </td>
                  <td className="px-5 py-3.5">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-green-50 text-green-700">{c.coupon_percentage}% off</span>
                  </td>
                  <td className="px-5 py-3.5">
                    <span className={`inline-flex items-center gap-1 text-xs ${expired ? 'text-red-500' : 'text-gray-600'}`}>
                      <Clock size={12} />
                      {new Date(c.expiry_date).toLocaleDateString()}
                      {expired && <span className="ml-1 font-medium">(Expired)</span>}
                    </span>
                  </td>
                  <td className="px-5 py-3.5 text-right">
                    <button onClick={() => deleteCoupon(c.coupon_name)} className="p-1.5 rounded-lg hover:bg-red-50 text-red-500 transition" title="Delete">
                      <Trash2 size={16} />
                    </button>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {coupons.length === 0 && <p className="text-center text-gray-400 py-10">No coupons yet.</p>}
      </div>
    </div>
  )
}
