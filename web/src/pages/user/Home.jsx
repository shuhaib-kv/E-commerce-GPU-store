import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import toast from 'react-hot-toast'
import api from '../../api/axios'
import ProductCard from '../../components/ProductCard'

export default function Home() {
  const [products, setProducts] = useState([])

  useEffect(() => {
    api.get('/user/viewproducts?pageSize=8&page=1')
      .then((res) => setProducts(res.data.data || []))
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load products'))
  }, [])

  return (
    <div>
      <section className="bg-gradient-to-r from-gray-900 to-gray-800 text-white py-20">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <h1 className="text-5xl font-bold mb-4">Premium GPUs</h1>
          <p className="text-xl text-gray-300 mb-8">
            Find the best graphics cards for gaming, AI, and professional workloads
          </p>
          <Link
            to="/products"
            className="bg-green-600 hover:bg-green-700 px-8 py-3 rounded-lg text-lg font-semibold transition"
          >
            Browse Products
          </Link>
        </div>
      </section>

      <section className="max-w-7xl mx-auto px-4 py-12">
        <h2 className="text-2xl font-bold mb-6">Featured Products</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {products.map((p) => (
            <ProductCard key={p.id} product={p} />
          ))}
        </div>
        {products.length === 0 && (
          <p className="text-center text-gray-500 py-8">No products available yet.</p>
        )}
      </section>
    </div>
  )
}
