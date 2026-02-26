import { useEffect, useState } from 'react'
import { Search, SlidersHorizontal, ChevronLeft, ChevronRight } from 'lucide-react'
import toast from 'react-hot-toast'
import api from '../../api/axios'
import ProductCard from '../../components/ProductCard'

export default function Products() {
  const [products, setProducts] = useState([])
  const [filters, setFilters] = useState({ name: '', brand: '', minPrice: '', maxPrice: '' })
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

  useEffect(() => { fetchProducts() }, [page])

  const handleFilter = (e) => {
    e.preventDefault()
    setPage(1)
    fetchProducts()
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Products</h1>

      <form onSubmit={handleFilter} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 mb-8">
        <div className="flex items-center gap-2 mb-3">
          <SlidersHorizontal size={18} className="text-gray-500" />
          <span className="text-sm font-medium text-gray-700">Filters</span>
        </div>
        <div className="flex flex-wrap gap-3 items-end">
          <div className="relative flex-1 min-w-[180px]">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              placeholder="Search name..."
              value={filters.name}
              onChange={(e) => setFilters({ ...filters, name: e.target.value })}
              className="border border-gray-200 rounded-lg pl-9 pr-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent text-sm"
            />
          </div>
          <input
            placeholder="Brand"
            value={filters.brand}
            onChange={(e) => setFilters({ ...filters, brand: e.target.value })}
            className="border border-gray-200 rounded-lg px-3 py-2 w-32 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent text-sm"
          />
          <input
            type="number"
            placeholder="Min ₹"
            value={filters.minPrice}
            onChange={(e) => setFilters({ ...filters, minPrice: e.target.value })}
            className="border border-gray-200 rounded-lg px-3 py-2 w-28 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent text-sm"
          />
          <input
            type="number"
            placeholder="Max ₹"
            value={filters.maxPrice}
            onChange={(e) => setFilters({ ...filters, maxPrice: e.target.value })}
            className="border border-gray-200 rounded-lg px-3 py-2 w-28 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent text-sm"
          />
          <button className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg transition text-sm font-medium">
            Apply
          </button>
        </div>
      </form>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((p) => (
          <ProductCard key={p.id} product={p} />
        ))}
      </div>

      {products.length === 0 && (
        <p className="text-center text-gray-500 py-16">No products found.</p>
      )}

      <div className="flex justify-center items-center gap-3 mt-10">
        <button
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page === 1}
          className="flex items-center gap-1 px-4 py-2 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed text-sm font-medium transition"
        >
          <ChevronLeft size={16} /> Previous
        </button>
        <span className="px-4 py-2 bg-green-50 text-green-700 rounded-lg text-sm font-semibold">Page {page}</span>
        <button
          onClick={() => setPage((p) => p + 1)}
          disabled={products.length < pageSize}
          className="flex items-center gap-1 px-4 py-2 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed text-sm font-medium transition"
        >
          Next <ChevronRight size={16} />
        </button>
      </div>
    </div>
  )
}
