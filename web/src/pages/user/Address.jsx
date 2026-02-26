import { useEffect, useState } from 'react'
import { MapPin, Plus, Pencil, X } from 'lucide-react'
import toast from 'react-hot-toast'
import api from '../../api/axios'

export default function Address() {
  const [addresses, setAddresses] = useState([])
  const [form, setForm] = useState({ name: '', phone_number: '', pincode: '', house: '', area: '', landmark: '', city: '' })
  const [editing, setEditing] = useState(null)
  const [showForm, setShowForm] = useState(false)

  const fetchAddresses = () => {
    api.get('/user/address')
      .then((res) => {
        const addr = res.data.data
        setAddresses(addr ? [addr] : [])
      })
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
      toast.success(editing ? 'Address updated' : 'Address added')
      setForm({ name: '', phone_number: '', pincode: '', house: '', area: '', landmark: '', city: '' })
      setEditing(null)
      setShowForm(false)
      fetchAddresses()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to save address')
    }
  }

  const startEdit = (addr) => {
    setForm({ name: addr.name, phone_number: addr.phone_number, pincode: addr.pincode, house: addr.house, area: addr.area, landmark: addr.landmark, city: addr.city })
    setEditing(addr.id)
    setShowForm(true)
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <div className="flex items-center gap-3">
          <div className="bg-green-50 text-green-600 p-2 rounded-lg">
            <MapPin size={22} />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">My Addresses</h1>
        </div>
        <button
          onClick={() => { setShowForm(!showForm); setEditing(null) }}
          className={`flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition ${
            showForm
              ? 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              : 'bg-green-600 text-white hover:bg-green-700'
          }`}
        >
          {showForm ? <><X size={16} /> Cancel</> : <><Plus size={16} /> Add Address</>}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 mb-8 grid grid-cols-1 md:grid-cols-2 gap-4">
          <input name="name" placeholder="Full Name" value={form.name} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="phone_number" placeholder="Phone" value={form.phone_number} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="house" placeholder="House/Flat" value={form.house} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="area" placeholder="Area/Street" value={form.area} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="landmark" placeholder="Landmark" value={form.landmark} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          <input name="city" placeholder="City" value={form.city} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="pincode" placeholder="Pincode" value={form.pincode} onChange={handleChange} className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <button className="bg-green-600 hover:bg-green-700 text-white py-2.5 rounded-lg transition font-medium text-sm md:col-span-2">
            {editing ? 'Update Address' : 'Save Address'}
          </button>
        </form>
      )}

      <div className="flex flex-col gap-4">
        {addresses.map((a) => (
          <div key={a.id} className="bg-white p-5 rounded-xl shadow-sm border border-gray-100 flex justify-between items-start">
            <div className="flex gap-3">
              <div className="bg-gray-50 text-gray-400 p-2 rounded-lg mt-0.5">
                <MapPin size={18} />
              </div>
              <div>
                <p className="font-semibold text-gray-900">{a.name}</p>
                <p className="text-sm text-gray-600 mt-0.5">{a.house}, {a.area}, {a.city} - {a.pincode}</p>
                <p className="text-sm text-gray-500 mt-0.5">{a.phone_number}</p>
              </div>
            </div>
            <button onClick={() => startEdit(a)} className="flex items-center gap-1 text-green-600 hover:text-green-700 text-sm font-medium">
              <Pencil size={14} /> Edit
            </button>
          </div>
        ))}
        {addresses.length === 0 && !showForm && (
          <div className="text-center py-16">
            <MapPin size={64} className="text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500 text-lg">No addresses saved.</p>
          </div>
        )}
      </div>
    </div>
  )
}
