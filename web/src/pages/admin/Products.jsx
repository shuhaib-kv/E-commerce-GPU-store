import { useEffect, useState, useRef } from 'react'
import toast from 'react-hot-toast'
import { Package, Plus, X, Pencil, Trash2, Upload, Image as ImageIcon } from 'lucide-react'
import api from '../../api/axios'

export default function AdminProducts() {
  const [products, setProducts] = useState([])
  const [categories, setCategories] = useState([])
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState({
    name: '', price: '', modelno: '', stock: '', category_id: '', description: '', brand: '',
  })
  const [images, setImages] = useState({ image1: null, image2: null, image3: null })
  const [previews, setPreviews] = useState({ image1: null, image2: null, image3: null })
  const fileRefs = { image1: useRef(), image2: useRef(), image3: useRef() }

  const fetchProducts = () => {
    api.get('/admin/product/view')
      .then((res) => setProducts(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load products'))
  }

  useEffect(() => {
    fetchProducts()
    api.get('/admin/category/view')
      .then((res) => setCategories(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load categories'))
  }, [])

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleImageChange = (key, file) => {
    setImages((prev) => ({ ...prev, [key]: file }))
    if (file) {
      const url = URL.createObjectURL(file)
      setPreviews((prev) => {
        if (prev[key]) URL.revokeObjectURL(prev[key])
        return { ...prev, [key]: url }
      })
    } else {
      setPreviews((prev) => {
        if (prev[key]) URL.revokeObjectURL(prev[key])
        return { ...prev, [key]: null }
      })
    }
  }

  const removeImage = (key) => {
    setImages((prev) => ({ ...prev, [key]: null }))
    setPreviews((prev) => {
      if (prev[key]) URL.revokeObjectURL(prev[key])
      return { ...prev, [key]: null }
    })
    if (fileRefs[key].current) fileRefs[key].current.value = ''
  }

  const resetForm = () => {
    setForm({ name: '', price: '', modelno: '', stock: '', category_id: '', description: '', brand: '' })
    setImages({ image1: null, image2: null, image3: null })
    Object.values(previews).forEach((url) => { if (url) URL.revokeObjectURL(url) })
    setPreviews({ image1: null, image2: null, image3: null })
    setEditing(null)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const formData = new FormData()
    Object.entries(form).forEach(([k, v]) => formData.append(k, v))
    if (images.image1) formData.append('image1', images.image1)
    if (images.image2) formData.append('image2', images.image2)
    if (images.image3) formData.append('image3', images.image3)

    try {
      if (editing) {
        const hasImages = images.image1 || images.image2 || images.image3
        if (hasImages) {
          await api.patch(`/admin/product/edit/${editing}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
        } else {
          const jsonBody = {}
          Object.entries(form).forEach(([k, v]) => { if (v !== '') jsonBody[k] = v })
          if (jsonBody.price) jsonBody.price = Number(jsonBody.price)
          if (jsonBody.modelno) jsonBody.modelno = Number(jsonBody.modelno)
          if (jsonBody.stock) jsonBody.stock = Number(jsonBody.stock)
          await api.patch(`/admin/product/edit/${editing}`, jsonBody)
        }
        toast.success('Product updated')
      } else {
        await api.post('/admin/product/add', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
        toast.success('Product added')
      }
      setShowForm(false)
      resetForm()
      fetchProducts()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to save product')
    }
  }

  const deleteProduct = async (id) => {
    if (!confirm('Delete this product?')) return
    try {
      await api.delete(`/admin/product/delete/${id}`)
      toast.success('Product deleted')
      fetchProducts()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to delete product')
    }
  }

  const startEdit = (p) => {
    setForm({
      name: p.name || '',
      price: p.price != null ? String(p.price) : '',
      modelno: p.modelno != null ? String(p.modelno) : '',
      stock: p.stock != null ? String(p.stock) : '',
      category_id: p.category_id || '',
      description: p.description || '',
      brand: p.brand || '',
    })
    setImages({ image1: null, image2: null, image3: null })
    Object.values(previews).forEach((url) => { if (url) URL.revokeObjectURL(url) })
    setPreviews({ image1: null, image2: null, image3: null })
    setEditing(p.id)
    setShowForm(true)
  }

  const imgLabel = ['Primary', 'Secondary', 'Additional']

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <Package size={24} className="text-emerald-600" />
            Products
          </h1>
          <p className="text-gray-500 text-sm mt-1">{products.length} total products</p>
        </div>
        <button
          onClick={() => { if (showForm) { setShowForm(false); resetForm() } else { resetForm(); setShowForm(true) } }}
          className={`flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition ${
            showForm ? 'bg-gray-100 text-gray-700 hover:bg-gray-200' : 'bg-green-600 text-white hover:bg-green-700 shadow-sm'
          }`}
        >
          {showForm ? <><X size={16} /> Cancel</> : <><Plus size={16} /> Add Product</>}
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl border border-gray-200 shadow-sm mb-6 overflow-hidden">
          {/* Form Header */}
          <div className="px-6 py-4 border-b border-gray-100 bg-gray-50/50">
            <h3 className="font-semibold text-gray-900 flex items-center gap-2">
              {editing ? <><Pencil size={16} className="text-blue-600" /> Edit Product</> : <><Plus size={16} className="text-green-600" /> New Product</>}
            </h3>
            <p className="text-xs text-gray-500 mt-0.5">{editing ? 'Update product details below' : 'Fill in the details to add a new product'}</p>
          </div>

          <form onSubmit={handleSubmit} className="p-6">
            {/* Basic Info Section */}
            <div className="mb-6">
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Basic Information</p>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Product Name</label>
                  <input name="name" value={form.name} onChange={handleChange} placeholder="e.g. NVIDIA GeForce RTX 4090" className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent placeholder:text-gray-300" required />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Brand</label>
                  <input name="brand" value={form.brand} onChange={handleChange} placeholder="e.g. NVIDIA" className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent placeholder:text-gray-300" required />
                </div>
              </div>
            </div>

            {/* Pricing & Inventory */}
            <div className="mb-6">
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Pricing & Inventory</p>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Price (₹)</label>
                  <input name="price" type="number" value={form.price} onChange={handleChange} placeholder="0" className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent placeholder:text-gray-300" required />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Model No</label>
                  <input name="modelno" type="number" value={form.modelno} onChange={handleChange} placeholder="0" className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent placeholder:text-gray-300" required />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Stock</label>
                  <input name="stock" type="number" value={form.stock} onChange={handleChange} placeholder="0" className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent placeholder:text-gray-300" required />
                </div>
              </div>
            </div>

            {/* Category & Description */}
            <div className="mb-6">
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Details</p>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Category</label>
                  <select name="category_id" value={form.category_id} onChange={handleChange} className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent" required>
                    <option value="">Select Category</option>
                    {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                  <textarea name="description" value={form.description} onChange={handleChange} placeholder="Product description..." className="w-full border border-gray-200 rounded-lg px-3.5 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent placeholder:text-gray-300" rows={1} />
                </div>
              </div>
            </div>

            {/* Images Section */}
              <div className="mb-6">
                <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Product Images {editing && <span className="normal-case font-normal text-gray-300">— upload to replace existing</span>}</p>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {['image1', 'image2', 'image3'].map((key, i) => (
                    <div key={key}>
                      <label className="block text-sm font-medium text-gray-700 mb-1.5">{imgLabel[i]} Image</label>
                      {previews[key] ? (
                        <div className="relative group">
                          <div className="w-full h-40 rounded-lg overflow-hidden border border-gray-200 bg-gray-50">
                            <img src={previews[key]} alt={`Preview ${i + 1}`} className="w-full h-full object-cover" />
                          </div>
                          <button
                            type="button"
                            onClick={() => removeImage(key)}
                            className="absolute top-2 right-2 w-7 h-7 bg-red-500 hover:bg-red-600 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition shadow-md"
                          >
                            <X size={14} />
                          </button>
                          <p className="text-xs text-gray-500 mt-1.5 truncate">{images[key]?.name}</p>
                        </div>
                      ) : (
                        <button
                          type="button"
                          onClick={() => fileRefs[key].current?.click()}
                          className="w-full h-40 rounded-lg border-2 border-dashed border-gray-200 hover:border-green-400 bg-gray-50 hover:bg-green-50/30 flex flex-col items-center justify-center gap-2 transition cursor-pointer group"
                        >
                          <div className="w-10 h-10 rounded-full bg-gray-100 group-hover:bg-green-100 flex items-center justify-center transition">
                            <Upload size={18} className="text-gray-400 group-hover:text-green-600 transition" />
                          </div>
                          <span className="text-xs text-gray-400 group-hover:text-green-600 font-medium transition">Click to upload</span>
                          <span className="text-[10px] text-gray-300">PNG, JPG up to 5MB</span>
                        </button>
                      )}
                      <input
                        ref={fileRefs[key]}
                        type="file"
                        accept="image/*"
                        onChange={(e) => handleImageChange(key, e.target.files[0])}
                        className="hidden"
                      />
                    </div>
                  ))}
                </div>
              </div>

            {/* Submit */}
            <div className="flex items-center gap-3 pt-2">
              <button type="submit" className="flex-1 flex items-center justify-center gap-2 bg-green-600 hover:bg-green-700 text-white py-2.5 rounded-lg transition font-semibold text-sm shadow-sm">
                {editing ? <><Pencil size={15} /> Update Product</> : <><Plus size={15} /> Add Product</>}
              </button>
              <button
                type="button"
                onClick={() => { setShowForm(false); resetForm() }}
                className="px-6 py-2.5 rounded-lg border border-gray-200 text-gray-600 hover:bg-gray-50 transition text-sm font-medium"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Product Table */}
      <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-gray-500 text-xs uppercase tracking-wider bg-gray-50/80">
              <th className="px-5 py-3.5 font-medium">Product</th>
              <th className="px-5 py-3.5 font-medium">Brand</th>
              <th className="px-5 py-3.5 font-medium">Price</th>
              <th className="px-5 py-3.5 font-medium">Stock</th>
              <th className="px-5 py-3.5 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {products.map((p) => {
              const imgSrc = p.image1
                ? p.image1.startsWith('http') ? p.image1 : `/public/images/${p.image1}`
                : null
              return (
                <tr key={p.id} className="hover:bg-gray-50/50 transition">
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-gray-100 overflow-hidden flex-shrink-0 flex items-center justify-center">
                        {imgSrc ? (
                          <img src={imgSrc} alt={p.name} className="w-full h-full object-cover" />
                        ) : (
                          <ImageIcon size={18} className="text-gray-400" />
                        )}
                      </div>
                      <span className="font-medium text-gray-900">{p.name}</span>
                    </div>
                  </td>
                  <td className="px-5 py-3 text-gray-600">{p.brand}</td>
                  <td className="px-5 py-3 font-semibold text-gray-900">₹{(p.price || 0).toLocaleString()}</td>
                  <td className="px-5 py-3">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      (p.stock || 0) < 10 ? 'bg-red-50 text-red-700' : 'bg-green-50 text-green-700'
                    }`}>
                      {p.stock || 0}
                    </span>
                  </td>
                  <td className="px-5 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button onClick={() => startEdit(p)} className="p-1.5 rounded-lg hover:bg-blue-50 text-blue-600 transition" title="Edit">
                        <Pencil size={16} />
                      </button>
                      <button onClick={() => deleteProduct(p.id)} className="p-1.5 rounded-lg hover:bg-red-50 text-red-500 transition" title="Delete">
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {products.length === 0 && <p className="text-center text-gray-400 py-10">No products found.</p>}
      </div>
    </div>
  )
}
