import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { UserPlus, Shield, Ban, CheckCircle, Eye, X } from 'lucide-react'
import api from '../../api/axios'

export default function ManageAdmins() {
  const [admins, setAdmins] = useState([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)
  const [showDetail, setShowDetail] = useState(null)
  const [form, setForm] = useState({
    name: '', email: '', password: '', store_name: '',
    payment_gateway: { provider: 'razorpay', key: '', secret: '' }
  })

  const fetchAdmins = async () => {
    try {
      const res = await api.get('/admin/manage/admins')
      setAdmins(res.data.data || [])
    } catch { toast.error('Failed to fetch admins') }
    setLoading(false)
  }

  useEffect(() => { fetchAdmins() }, [])

  const handleCreate = async (e) => {
    e.preventDefault()
    try {
      await api.post('/admin/signup', form)
      toast.success('Admin created successfully')
      setShowCreate(false)
      setForm({ name: '', email: '', password: '', store_name: '', payment_gateway: { provider: 'razorpay', key: '', secret: '' } })
      fetchAdmins()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to create admin')
    }
  }

  const handleBlock = async (id) => {
    try {
      await api.patch(`/admin/manage/admins/${id}/block`)
      toast.success('Admin blocked')
      fetchAdmins()
    } catch (err) { toast.error(err.response?.data?.message || 'Failed to block') }
  }

  const handleUnblock = async (id) => {
    try {
      await api.patch(`/admin/manage/admins/${id}/unblock`)
      toast.success('Admin unblocked')
      fetchAdmins()
    } catch (err) { toast.error(err.response?.data?.message || 'Failed to unblock') }
  }

  const handleUpdate = async (id, data) => {
    try {
      await api.patch(`/admin/manage/admins/${id}`, data)
      toast.success('Admin updated')
      fetchAdmins()
      setShowDetail(null)
    } catch (err) { toast.error(err.response?.data?.message || 'Failed to update') }
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64"><div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" /></div>
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Manage Admins</h1>
          <p className="text-sm text-gray-500">Create, edit and manage store admins</p>
        </div>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="flex items-center gap-2 bg-purple-600 hover:bg-purple-700 text-white px-4 py-2.5 rounded-xl text-sm font-medium transition"
        >
          <UserPlus size={16} />
          Create Admin
        </button>
      </div>

      {/* Create Admin Form */}
      {showCreate && (
        <div className="bg-white rounded-xl border border-gray-200 p-6 mb-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Create New Admin</h3>
          <form onSubmit={handleCreate} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Name</label>
              <input type="text" required value={form.name} onChange={e => setForm({ ...form, name: e.target.value })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Email</label>
              <input type="email" required value={form.email} onChange={e => setForm({ ...form, email: e.target.value })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Password</label>
              <input type="password" required value={form.password} onChange={e => setForm({ ...form, password: e.target.value })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Store Name</label>
              <input type="text" required value={form.store_name} onChange={e => setForm({ ...form, store_name: e.target.value })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>

            <div className="md:col-span-2">
              <p className="text-sm font-semibold text-gray-700 mb-2">Payment Gateway</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">Provider</label>
              <select value={form.payment_gateway.provider}
                onChange={e => setForm({ ...form, payment_gateway: { ...form.payment_gateway, provider: e.target.value } })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500">
                <option value="razorpay">Razorpay</option>
                <option value="stripe">Stripe</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">API Key</label>
              <input type="text" value={form.payment_gateway.key}
                onChange={e => setForm({ ...form, payment_gateway: { ...form.payment_gateway, key: e.target.value } })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">API Secret</label>
              <input type="password" value={form.payment_gateway.secret}
                onChange={e => setForm({ ...form, payment_gateway: { ...form.payment_gateway, secret: e.target.value } })}
                className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>

            <div className="md:col-span-2 flex gap-3">
              <button type="submit" className="bg-purple-600 hover:bg-purple-700 text-white px-6 py-2.5 rounded-lg text-sm font-medium transition">
                Create Admin
              </button>
              <button type="button" onClick={() => setShowCreate(false)} className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-6 py-2.5 rounded-lg text-sm font-medium transition">
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Edit Admin Modal */}
      {showDetail && (
        <EditAdminModal admin={showDetail} onClose={() => setShowDetail(null)} onSave={handleUpdate} />
      )}

      {/* Admins Table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-gray-50 text-left text-gray-500 border-b border-gray-200">
              <th className="px-6 py-3 font-medium">Admin</th>
              <th className="px-6 py-3 font-medium">Store Name</th>
              <th className="px-6 py-3 font-medium">Payment</th>
              <th className="px-6 py-3 font-medium">Status</th>
              <th className="px-6 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {admins.map((admin) => (
              <tr key={admin.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-9 h-9 bg-purple-100 rounded-full flex items-center justify-center text-purple-600 font-semibold text-sm">
                      {admin.name?.charAt(0)?.toUpperCase()}
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">{admin.name}</p>
                      <p className="text-xs text-gray-500">{admin.email}</p>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 text-gray-700">{admin.store_name || '-'}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded-md text-xs font-medium ${admin.payment_gateway?.key ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'}`}>
                    {admin.payment_gateway?.provider || 'None'}
                    {admin.payment_gateway?.key ? ' (configured)' : ''}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <span className={`px-2.5 py-1 rounded-full text-xs font-medium ${admin.block_status ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                    {admin.block_status ? 'Blocked' : 'Active'}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2">
                    <button onClick={() => setShowDetail(admin)} title="Edit"
                      className="p-1.5 rounded-lg hover:bg-gray-100 text-gray-500 hover:text-gray-700 transition">
                      <Eye size={16} />
                    </button>
                    {admin.block_status ? (
                      <button onClick={() => handleUnblock(admin.id)} title="Unblock"
                        className="p-1.5 rounded-lg hover:bg-green-50 text-green-600 hover:text-green-700 transition">
                        <CheckCircle size={16} />
                      </button>
                    ) : (
                      <button onClick={() => handleBlock(admin.id)} title="Block"
                        className="p-1.5 rounded-lg hover:bg-red-50 text-red-500 hover:text-red-700 transition">
                        <Ban size={16} />
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
            {admins.length === 0 && (
              <tr><td colSpan={5} className="px-6 py-12 text-center text-gray-400">No admins found. Create one to get started.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function EditAdminModal({ admin, onClose, onSave }) {
  const [storeName, setStoreName] = useState(admin.store_name || '')
  const [gateway, setGateway] = useState(admin.payment_gateway || { provider: 'razorpay', key: '', secret: '' })

  const handleSubmit = (e) => {
    e.preventDefault()
    onSave(admin.id, { store_name: storeName, payment_gateway: gateway })
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-2xl w-full max-w-lg p-6 shadow-xl">
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-semibold text-gray-900">Edit Admin: {admin.name}</h3>
          <button onClick={onClose} className="p-1 rounded-lg hover:bg-gray-100"><X size={18} /></button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Store Name</label>
            <input type="text" value={storeName} onChange={e => setStoreName(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Payment Provider</label>
            <select value={gateway.provider} onChange={e => setGateway({ ...gateway, provider: e.target.value })}
              className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500">
              <option value="razorpay">Razorpay</option>
              <option value="stripe">Stripe</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">API Key</label>
            <input type="text" value={gateway.key} onChange={e => setGateway({ ...gateway, key: e.target.value })}
              className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">API Secret</label>
            <input type="password" value={gateway.secret} onChange={e => setGateway({ ...gateway, secret: e.target.value })}
              className="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
          </div>
          <div className="flex gap-3 pt-2">
            <button type="submit" className="bg-purple-600 hover:bg-purple-700 text-white px-6 py-2.5 rounded-lg text-sm font-medium transition">Save Changes</button>
            <button type="button" onClick={onClose} className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-6 py-2.5 rounded-lg text-sm font-medium transition">Cancel</button>
          </div>
        </form>
      </div>
    </div>
  )
}
