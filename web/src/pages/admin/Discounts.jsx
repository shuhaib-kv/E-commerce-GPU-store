import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Discounts</h1>
        <button onClick={() => setShowForm(!showForm)} className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded transition">
          {showForm ? 'Cancel' : 'Add Discount'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow mb-8 flex flex-wrap gap-4">
          <input placeholder="Discount Name" value={form.discountname} onChange={(e) => setForm({ ...form, discountname: e.target.value })} className="border rounded px-3 py-2 flex-1" required />
          <input type="number" placeholder="Percentage" value={form.discountpercentage} onChange={(e) => setForm({ ...form, discountpercentage: e.target.value })} className="border rounded px-3 py-2 w-32" required />
          <button className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded transition">Add</button>
        </form>
      )}

      <div className="bg-white rounded-lg shadow">
        {discounts.map((d) => (
          <div key={d.id} className="flex justify-between items-center px-4 py-3 border-b hover:bg-gray-50">
            <div>
              <span className="font-medium">{d.discount_name}</span>
              <span className="ml-3 text-green-600 font-bold">{d.discount_percentage}% off</span>
            </div>
            <button onClick={() => deleteDiscount(d.id)} className="text-red-600 hover:underline text-sm">Delete</button>
          </div>
        ))}
        {discounts.length === 0 && <p className="text-center text-gray-500 py-8">No discounts yet.</p>}
      </div>
    </div>
  )
}
