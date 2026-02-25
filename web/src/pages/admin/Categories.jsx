import { useEffect, useState } from 'react'
import api from '../../api/axios'

export default function Categories() {
  const [categories, setCategories] = useState([])
  const [name, setName] = useState('')
  const [editing, setEditing] = useState(null)

  const fetchCategories = () => {
    api.get('/admin/category/view')
      .then((res) => setCategories(res.data.categories || []))
      .catch(() => {})
  }

  useEffect(() => { fetchCategories() }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      if (editing) {
        await api.patch('/admin/category/edit', { category_id: editing, name })
      } else {
        await api.post('/admin/category/add', { name })
      }
      setName('')
      setEditing(null)
      fetchCategories()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to save category')
    }
  }

  const deleteCategory = async (id) => {
    if (!confirm('Delete this category?')) return
    await api.delete('/admin/category/delete', { data: { category_id: id } })
    fetchCategories()
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Categories</h1>

      <form onSubmit={handleSubmit} className="bg-white p-4 rounded-lg shadow mb-8 flex gap-4">
        <input
          placeholder="Category name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="border rounded px-3 py-2 flex-1"
          required
        />
        <button className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded transition">
          {editing ? 'Update' : 'Add'}
        </button>
        {editing && (
          <button type="button" onClick={() => { setEditing(null); setName('') }} className="text-gray-500 hover:underline">
            Cancel
          </button>
        )}
      </form>

      <div className="bg-white rounded-lg shadow">
        {categories.map((c) => (
          <div key={c._id} className="flex justify-between items-center px-4 py-3 border-b hover:bg-gray-50">
            <span className="font-medium">{c.name}</span>
            <div className="flex gap-3">
              <button onClick={() => { setEditing(c._id); setName(c.name) }} className="text-blue-600 hover:underline text-sm">Edit</button>
              <button onClick={() => deleteCategory(c._id)} className="text-red-600 hover:underline text-sm">Delete</button>
            </div>
          </div>
        ))}
        {categories.length === 0 && <p className="text-center text-gray-500 py-8">No categories yet.</p>}
      </div>
    </div>
  )
}
