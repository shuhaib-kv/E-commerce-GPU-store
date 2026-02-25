import { useEffect, useState } from 'react'
import api from '../../api/axios'

export default function Address() {
  const [addresses, setAddresses] = useState([])
  const [form, setForm] = useState({ name: '', phone_number: '', pincode: '', house: '', area: '', landmark: '', city: '' })
  const [editing, setEditing] = useState(null)
  const [showForm, setShowForm] = useState(false)

  const fetchAddresses = () => {
    api.get('/user/address')
      .then((res) => setAddresses(res.data.address || []))
      .catch(() => {})
  }

  useEffect(() => { fetchAddresses() }, [])

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      if (editing) {
        await api.patch('/user/edit/address', { ...form, address_id: editing })
      } else {
        await api.post('/user/add/address', form)
      }
      setForm({ name: '', phone_number: '', pincode: '', house: '', area: '', landmark: '', city: '' })
      setEditing(null)
      setShowForm(false)
      fetchAddresses()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to save address')
    }
  }

  const startEdit = (addr) => {
    setForm({ name: addr.name, phone_number: addr.phone_number, pincode: addr.pincode, house: addr.house, area: addr.area, landmark: addr.landmark, city: addr.city })
    setEditing(addr._id)
    setShowForm(true)
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">My Addresses</h1>
        <button onClick={() => { setShowForm(!showForm); setEditing(null) }} className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded transition">
          {showForm ? 'Cancel' : 'Add Address'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow mb-8 grid grid-cols-1 md:grid-cols-2 gap-4">
          <input name="name" placeholder="Full Name" value={form.name} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="phone_number" placeholder="Phone" value={form.phone_number} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="house" placeholder="House/Flat" value={form.house} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="area" placeholder="Area/Street" value={form.area} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="landmark" placeholder="Landmark" value={form.landmark} onChange={handleChange} className="border rounded px-3 py-2" />
          <input name="city" placeholder="City" value={form.city} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="pincode" placeholder="Pincode" value={form.pincode} onChange={handleChange} className="border rounded px-3 py-2" required />
          <button className="bg-green-600 hover:bg-green-700 text-white py-2 rounded transition font-semibold md:col-span-2">
            {editing ? 'Update Address' : 'Save Address'}
          </button>
        </form>
      )}

      <div className="flex flex-col gap-4">
        {addresses.map((a) => (
          <div key={a._id} className="bg-white p-4 rounded-lg shadow flex justify-between items-center">
            <div>
              <p className="font-semibold">{a.name}</p>
              <p className="text-sm text-gray-600">{a.house}, {a.area}, {a.city} - {a.pincode}</p>
              <p className="text-sm text-gray-500">{a.phone_number}</p>
            </div>
            <button onClick={() => startEdit(a)} className="text-green-600 hover:underline text-sm">Edit</button>
          </div>
        ))}
        {addresses.length === 0 && <p className="text-gray-500 text-center py-8">No addresses saved.</p>}
      </div>
    </div>
  )
}
