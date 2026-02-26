import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import api from '../../api/axios'
import ProductCard from '../../components/ProductCard'

export default function Products() {
  const [products, setProducts] = useState([])
  const [filters, setFilters] = useState({ name: '', brand: '', minPrice: '', maxPrice: '' })
  const [categories, setCategories] = useState([])
  const [page, setPage] = useState(1)
  const pageSize = 12

  const fetchProducts = () => {
    const params = new URLSearchParams()
    params.set('pageSize', pageSize)
    params.set('page', page)
    if (filters.name) params.set('name', filters.name)
    if (filters.brand) params.set('brand', filters.brand)
    if (filters.minPrice) params.set('minPrice', filters.minPrice)
    if (filters.maxPrice) params.set('maxPrice', filters.maxPrice)

    api.get(`/user/viewproducts?${params}`)
      .then((res) => setProducts(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load products'))
  }

  useEffect(() => {
    api.get('/admin/category/view')
      .then((res) => setCategories(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load categories'))
  }, [])

  useEffect(() => { fetchProducts() }, [page])

  const handleFilter = (e) => {
    e.preventDefault()
    setPage(1)
    fetchProducts()
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Products</h1>

      <form onSubmit={handleFilter} className="bg-white p-4 rounded-lg shadow mb-8 flex flex-wrap gap-4 items-end">
        <input
          placeholder="Search name..."
          value={filters.name}
          onChange={(e) => setFilters({ ...filters, name: e.target.value })}
          className="border rounded px-3 py-2 flex-1 min-w-[150px]"
        />
        <input
          placeholder="Brand"
          value={filters.brand}
          onChange={(e) => setFilters({ ...filters, brand: e.target.value })}
          className="border rounded px-3 py-2 w-32"
        />
        <input
          type="number"
          placeholder="Min price"
          value={filters.minPrice}
          onChange={(e) => setFilters({ ...filters, minPrice: e.target.value })}
          className="border rounded px-3 py-2 w-28"
        />
        <input
          type="number"
          placeholder="Max price"
          value={filters.maxPrice}
          onChange={(e) => setFilters({ ...filters, maxPrice: e.target.value })}
          className="border rounded px-3 py-2 w-28"
        />
        <button className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded transition">
          Filter
        </button>
      </form>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((p) => (
          <ProductCard key={p.id} product={p} />
        ))}
      </div>

      {products.length === 0 && (
        <p className="text-center text-gray-500 py-12">No products found.</p>
      )}

      <div className="flex justify-center gap-4 mt-8">
        <button
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page === 1}
          className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300 disabled:opacity-50"
        >
          Previous
        </button>
        <span className="px-4 py-2">Page {page}</span>
        <button
          onClick={() => setPage((p) => p + 1)}
          disabled={products.length < pageSize}
          className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300 disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  )
}
