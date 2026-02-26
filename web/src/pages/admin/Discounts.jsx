import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { Percent, Plus, X, Trash2 } from 'lucide-react'
import api from '../../api/axios'

export default function Discounts() {
  const [discounts, setDiscounts] = useState([])
  const [form, setForm] = useState({ discountname: '', discountpercentage: '' })
  const [showForm, setShowForm] = useState(false)

  const fetchDiscounts = () => {
    api.get('/admin/discount')
      .then((res) => setDiscounts(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load discounts'))
  }

  useEffect(() => { fetchDiscounts() }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await api.post('/admin/add/discount', {
        discountname: form.discountname,
        discountpercentage: parseInt(form.discountpercentage),
      })
      toast.success('Discount added')
      setForm({ discountname: '', discountpercentage: '' })
      setShowForm(false)
      fetchDiscounts()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to add discount')
    }
  }

  const deleteDiscount = async (id) => {
    if (!confirm('Delete this discount?')) return
    try {
      await api.delete('/admin/delete/discount', { data: { id } })
      toast.success('Discount deleted')
      fetchDiscounts()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to delete discount')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <Percent size={24} className="text-teal-600" />
            Discounts
          </h1>
          <p className="text-gray-500 text-sm mt-1">{discounts.length} discounts</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition ${
            showForm ? 'bg-gray-100 text-gray-700 hover:bg-gray-200' : 'bg-green-600 text-white hover:bg-green-700 shadow-sm'
          }`}
        >
          {showForm ? <><X size={16} /> Cancel</> : <><Plus size={16} /> Add Discount</>}
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-6 mb-6">
          <h3 className="font-semibold text-gray-900 mb-4">New Discount</h3>
          <form onSubmit={handleSubmit} className="flex flex-wrap gap-4 items-end">
            <div className="flex-1 min-w-[200px]">
              <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider mb-1.5">Discount Name</label>
              <input value={form.discountname} onChange={(e) => setForm({ ...form, discountname: e.target.value })} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent" required />
            </div>
            <div className="w-36">
              <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider mb-1.5">Percentage</label>
              <input type="number" value={form.discountpercentage} onChange={(e) => setForm({ ...form, discountpercentage: e.target.value })} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent" required />
            </div>
            <button className="flex items-center gap-2 bg-green-600 hover:bg-green-700 text-white px-5 py-2 rounded-lg transition text-sm font-medium shadow-sm">
              <Plus size={16} /> Add
            </button>
          </form>
        </div>
      )}

      <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        {discounts.length === 0 ? (
          <p className="text-center text-gray-400 py-10">No discounts yet.</p>
        ) : (
          <div className="divide-y divide-gray-100">
            {discounts.map((d) => (
              <div key={d.id} className="flex justify-between items-center px-5 py-4 hover:bg-gray-50/50 transition">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-teal-100 text-teal-700 rounded-xl flex items-center justify-center">
                    <Percent size={18} />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{d.discount_name}</p>
                    <p className="text-sm text-teal-600 font-semibold">{d.discount_percentage}% off</p>
                  </div>
                </div>
                <button onClick={() => deleteDiscount(d.id)} className="p-1.5 rounded-lg hover:bg-red-50 text-red-500 transition" title="Delete">
                  <Trash2 size={16} />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
