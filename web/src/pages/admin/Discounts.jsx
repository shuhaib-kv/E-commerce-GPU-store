import { useEffect, useState } from 'react'
import api from '../../api/axios'

export default function Discounts() {
  const [discounts, setDiscounts] = useState([])
  const [form, setForm] = useState({ discount_name: '', discount_percentage: '' })
  const [showForm, setShowForm] = useState(false)

  const fetchDiscounts = () => {
    api.get('/admin/discount')
      .then((res) => setDiscounts(res.data.discounts || []))
      .catch(() => {})
  }

  useEffect(() => { fetchDiscounts() }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await api.post('/admin/add/discount', {
        ...form,
        discount_percentage: parseInt(form.discount_percentage),
      })
      setForm({ discount_name: '', discount_percentage: '' })
      setShowForm(false)
      fetchDiscounts()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to add discount')
    }
  }

  const deleteDiscount = async (id) => {
    if (!confirm('Delete this discount?')) return
    await api.delete('/admin/delete/discount', { data: { discount_id: id } })
    fetchDiscounts()
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
          <input placeholder="Discount Name" value={form.discount_name} onChange={(e) => setForm({ ...form, discount_name: e.target.value })} className="border rounded px-3 py-2 flex-1" required />
          <input type="number" placeholder="Percentage" value={form.discount_percentage} onChange={(e) => setForm({ ...form, discount_percentage: e.target.value })} className="border rounded px-3 py-2 w-32" required />
          <button className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded transition">Add</button>
        </form>
      )}

      <div className="bg-white rounded-lg shadow">
        {discounts.map((d) => (
          <div key={d._id} className="flex justify-between items-center px-4 py-3 border-b hover:bg-gray-50">
            <div>
              <span className="font-medium">{d.discount_name}</span>
              <span className="ml-3 text-green-600 font-bold">{d.discount_percentage}% off</span>
            </div>
            <button onClick={() => deleteDiscount(d._id)} className="text-red-600 hover:underline text-sm">Delete</button>
          </div>
        ))}
        {discounts.length === 0 && <p className="text-center text-gray-500 py-8">No discounts yet.</p>}
      </div>
    </div>
  )
}
