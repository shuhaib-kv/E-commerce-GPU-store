import { useEffect, useState } from 'react'
import api from '../../api/axios'

export default function AdminProducts() {
  const [products, setProducts] = useState([])
  const [categories, setCategories] = useState([])
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState({
    name: '', price: '', model_no: '', stock: '', category_id: '', description: '', brand: '',
  })
  const [images, setImages] = useState({ image1: null, image2: null, image3: null })

  const fetchProducts = () => {
    api.get('/admin/product/view')
      .then((res) => setProducts(res.data.products || []))
      .catch(() => {})
  }

  useEffect(() => {
    fetchProducts()
    api.get('/admin/category/view')
      .then((res) => setCategories(res.data.categories || []))
      .catch(() => {})
  }, [])

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    const formData = new FormData()
    Object.entries(form).forEach(([k, v]) => formData.append(k, v))
    if (images.image1) formData.append('image1', images.image1)
    if (images.image2) formData.append('image2', images.image2)
    if (images.image3) formData.append('image3', images.image3)

    try {
      if (editing) {
        await api.patch(`/admin/product/edit/${editing}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      } else {
        await api.post('/admin/product/add', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      }
      setShowForm(false)
      setEditing(null)
      setForm({ name: '', price: '', model_no: '', stock: '', category_id: '', description: '', brand: '' })
      setImages({ image1: null, image2: null, image3: null })
      fetchProducts()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to save product')
    }
  }

  const deleteProduct = async (id) => {
    if (!confirm('Delete this product?')) return
    await api.delete(`/admin/product/delete/${id}`)
    fetchProducts()
  }

  const startEdit = (p) => {
    setForm({ name: p.name, price: p.price, model_no: p.model_no, stock: p.stock, category_id: p.category_id, description: p.description, brand: p.brand })
    setEditing(p._id)
    setShowForm(true)
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Products</h1>
        <button onClick={() => { setShowForm(!showForm); setEditing(null) }} className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded transition">
          {showForm ? 'Cancel' : 'Add Product'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow mb-8 grid grid-cols-1 md:grid-cols-2 gap-4">
          <input name="name" placeholder="Product Name" value={form.name} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="brand" placeholder="Brand" value={form.brand} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="price" type="number" placeholder="Price" value={form.price} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="model_no" type="number" placeholder="Model No" value={form.model_no} onChange={handleChange} className="border rounded px-3 py-2" required />
          <input name="stock" type="number" placeholder="Stock" value={form.stock} onChange={handleChange} className="border rounded px-3 py-2" required />
          <select name="category_id" value={form.category_id} onChange={handleChange} className="border rounded px-3 py-2" required>
            <option value="">Select Category</option>
            {categories.map((c) => <option key={c._id} value={c._id}>{c.name}</option>)}
          </select>
          <textarea name="description" placeholder="Description" value={form.description} onChange={handleChange} className="border rounded px-3 py-2 md:col-span-2" rows={3} />
          <div>
            <label className="block text-sm mb-1">Image 1</label>
            <input type="file" accept="image/*" onChange={(e) => setImages({ ...images, image1: e.target.files[0] })} className="text-sm" />
          </div>
          <div>
            <label className="block text-sm mb-1">Image 2</label>
            <input type="file" accept="image/*" onChange={(e) => setImages({ ...images, image2: e.target.files[0] })} className="text-sm" />
          </div>
          <div>
            <label className="block text-sm mb-1">Image 3</label>
            <input type="file" accept="image/*" onChange={(e) => setImages({ ...images, image3: e.target.files[0] })} className="text-sm" />
          </div>
          <button className="bg-green-600 hover:bg-green-700 text-white py-2 rounded transition font-semibold md:col-span-2">
            {editing ? 'Update Product' : 'Add Product'}
          </button>
        </form>
      )}

      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full text-left">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Brand</th>
              <th className="px-4 py-3">Price</th>
              <th className="px-4 py-3">Stock</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((p) => (
              <tr key={p._id} className="border-b hover:bg-gray-50">
                <td className="px-4 py-3">{p.name}</td>
                <td className="px-4 py-3">{p.brand}</td>
                <td className="px-4 py-3">₹{p.price}</td>
                <td className="px-4 py-3">{p.stock}</td>
                <td className="px-4 py-3 flex gap-2">
                  <button onClick={() => startEdit(p)} className="text-blue-600 hover:underline text-sm">Edit</button>
                  <button onClick={() => deleteProduct(p._id)} className="text-red-600 hover:underline text-sm">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {products.length === 0 && <p className="text-center text-gray-500 py-8">No products found.</p>}
      </div>
    </div>
  )
}
