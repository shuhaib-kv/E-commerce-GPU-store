import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { FolderTree, Plus, Pencil, Trash2, X } from 'lucide-react'
import api from '../../api/axios'

export default function Categories() {
  const [categories, setCategories] = useState([])
  const [name, setName] = useState('')
  const [editing, setEditing] = useState(null)

  const fetchCategories = () => {
    api.get('/admin/category/view')
      .then((res) => setCategories(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load categories'))
  }

  useEffect(() => { fetchCategories() }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      if (editing) {
        await api.patch('/admin/category/edit', { id: editing, name })
        toast.success('Category updated')
      } else {
        await api.post('/admin/category/add', { name })
        toast.success('Category added')
      }
      setName('')
      setEditing(null)
      fetchCategories()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to save category')
    }
  }

  const deleteCategory = async (id) => {
    if (!confirm('Delete this category?')) return
    try {
      await api.delete('/admin/category/delete', { data: { id } })
      toast.success('Category deleted')
      fetchCategories()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to delete category')
    }
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <FolderTree size={24} className="text-indigo-600" />
          Categories
        </h1>
        <p className="text-gray-500 text-sm mt-1">{categories.length} categories</p>
      </div>

      <form onSubmit={handleSubmit} className="bg-white rounded-xl border border-gray-200 shadow-sm p-4 mb-6 flex gap-3 items-end">
        <div className="flex-1">
          <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider mb-1.5">Category Name</label>
          <input
            placeholder="Enter category name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent"
            required
          />
        </div>
        <button className="flex items-center gap-2 bg-green-600 hover:bg-green-700 text-white px-5 py-2 rounded-lg transition text-sm font-medium shadow-sm">
          <Plus size={16} />
          {editing ? 'Update' : 'Add'}
        </button>
        {editing && (
          <button type="button" onClick={() => { setEditing(null); setName('') }} className="flex items-center gap-1 text-gray-500 hover:text-gray-700 px-3 py-2 text-sm">
            <X size={16} /> Cancel
          </button>
        )}
      </form>

      <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        {categories.length === 0 ? (
          <p className="text-center text-gray-400 py-10">No categories yet.</p>
        ) : (
          <div className="divide-y divide-gray-100">
            {categories.map((c) => (
              <div key={c.id} className="flex justify-between items-center px-5 py-3.5 hover:bg-gray-50/50 transition">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 bg-indigo-100 text-indigo-700 rounded-lg flex items-center justify-center">
                    <FolderTree size={16} />
                  </div>
                  <span className="font-medium text-gray-900">{c.name}</span>
                </div>
                <div className="flex items-center gap-1">
                  <button onClick={() => { setEditing(c.id); setName(c.name) }} className="p-1.5 rounded-lg hover:bg-blue-50 text-blue-600 transition" title="Edit">
                    <Pencil size={16} />
                  </button>
                  <button onClick={() => deleteCategory(c.id)} className="p-1.5 rounded-lg hover:bg-red-50 text-red-500 transition" title="Delete">
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
