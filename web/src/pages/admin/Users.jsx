import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import api from '../../api/axios'

export default function Users() {
  const [users, setUsers] = useState([])

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

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">User Management</h1>
      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full text-left">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Email</th>
              <th className="px-4 py-3">Phone</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id} className="border-b hover:bg-gray-50">
                <td className="px-4 py-3">{u.first_name} {u.last_name}</td>
                <td className="px-4 py-3">{u.email}</td>
                <td className="px-4 py-3">{u.phone}</td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs ${u.block_status ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                    {u.block_status ? 'Blocked' : 'Active'}
                  </span>
                </td>
                <td className="px-4 py-3 flex gap-2">
                  {u.block_status ? (
                    <button onClick={() => unblockUser(u.id)} className="text-green-600 hover:underline text-sm">Unblock</button>
                  ) : (
                    <button onClick={() => blockUser(u.id)} className="text-yellow-600 hover:underline text-sm">Block</button>
                  )}
                  <button onClick={() => deleteUser(u.id)} className="text-red-600 hover:underline text-sm">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {users.length === 0 && <p className="text-center text-gray-500 py-8">No users found.</p>}
      </div>
    </div>
  )
}
