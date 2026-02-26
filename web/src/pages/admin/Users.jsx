import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { Users as UsersIcon, Search, ShieldOff, ShieldCheck, Trash2 } from 'lucide-react'
import api from '../../api/axios'

export default function Users() {
  const [users, setUsers] = useState([])
  const [search, setSearch] = useState('')

  const fetchUsers = () => {
    api.get('/admin/users')
      .then((res) => setUsers(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load users'))
  }

  useEffect(() => { fetchUsers() }, [])

  const blockUser = async (id) => {
    try {
      await api.patch('/admin/users/block', { id })
      toast.success('User blocked')
      fetchUsers()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to block user')
    }
  }

  const unblockUser = async (id) => {
    try {
      await api.patch('/admin/users/unblock', { id })
      toast.success('User unblocked')
      fetchUsers()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to unblock user')
    }
  }

  const deleteUser = async (id) => {
    if (!confirm('Delete this user?')) return
    try {
      await api.delete('/admin/users/delete', { data: { id } })
      toast.success('User deleted')
      fetchUsers()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to delete user')
    }
  }

  const filtered = users.filter(u =>
    `${u.first_name} ${u.last_name} ${u.email}`.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <UsersIcon size={24} className="text-blue-600" />
            User Management
          </h1>
          <p className="text-gray-500 text-sm mt-1">{users.length} total users</p>
        </div>
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search users..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 pr-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
          />
        </div>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-gray-500 text-xs uppercase tracking-wider bg-gray-50/80">
              <th className="px-5 py-3.5 font-medium">User</th>
              <th className="px-5 py-3.5 font-medium">Email</th>
              <th className="px-5 py-3.5 font-medium">Phone</th>
              <th className="px-5 py-3.5 font-medium">Status</th>
              <th className="px-5 py-3.5 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {filtered.map((u) => (
              <tr key={u.id} className="hover:bg-gray-50/50 transition">
                <td className="px-5 py-3.5">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-blue-100 text-blue-700 rounded-full flex items-center justify-center font-semibold text-xs">
                      {(u.first_name || 'U')[0]}{(u.last_name || '')[0]}
                    </div>
                    <span className="font-medium text-gray-900">{u.first_name} {u.last_name}</span>
                  </div>
                </td>
                <td className="px-5 py-3.5 text-gray-600">{u.email}</td>
                <td className="px-5 py-3.5 text-gray-600">{u.phone}</td>
                <td className="px-5 py-3.5">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    u.block_status ? 'bg-red-50 text-red-700' : 'bg-green-50 text-green-700'
                  }`}>
                    {u.block_status ? 'Blocked' : 'Active'}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <div className="flex items-center justify-end gap-1">
                    {u.block_status ? (
                      <button onClick={() => unblockUser(u.id)} className="p-1.5 rounded-lg hover:bg-green-50 text-green-600 transition" title="Unblock">
                        <ShieldCheck size={16} />
                      </button>
                    ) : (
                      <button onClick={() => blockUser(u.id)} className="p-1.5 rounded-lg hover:bg-amber-50 text-amber-600 transition" title="Block">
                        <ShieldOff size={16} />
                      </button>
                    )}
                    <button onClick={() => deleteUser(u.id)} className="p-1.5 rounded-lg hover:bg-red-50 text-red-500 transition" title="Delete">
                      <Trash2 size={16} />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {filtered.length === 0 && <p className="text-center text-gray-400 py-10">No users found.</p>}
      </div>
    </div>
  )
}
