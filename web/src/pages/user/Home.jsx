import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ArrowRight, Cpu, Zap, Shield, Truck } from 'lucide-react'
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
      {/* Hero */}
      <section className="bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white py-24 relative overflow-hidden">
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-10 left-10"><Cpu size={120} /></div>
          <div className="absolute bottom-10 right-10"><Cpu size={160} /></div>
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"><Cpu size={200} /></div>
        </div>
        <div className="max-w-7xl mx-auto px-4 text-center relative z-10">
          <div className="inline-flex items-center gap-2 bg-green-600/20 text-green-400 px-4 py-1.5 rounded-full text-sm font-medium mb-6">
            <Zap size={14} />
            Premium Graphics Cards
          </div>
          <h1 className="text-5xl md:text-6xl font-bold mb-4 leading-tight">
            Power Your <span className="text-green-400">Vision</span>
          </h1>
          <p className="text-xl text-gray-400 mb-8 max-w-2xl mx-auto">
            Find the best GPUs for gaming, AI, and professional workloads at unbeatable prices.
          </p>
          <Link
            to="/products"
            className="inline-flex items-center gap-2 bg-green-600 hover:bg-green-700 px-8 py-3 rounded-xl text-lg font-semibold transition shadow-lg shadow-green-600/25"
          >
            Browse Products
            <ArrowRight size={20} />
          </Link>
        </div>
      </section>

      {/* Features */}
      <section className="max-w-7xl mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 -mt-16 relative z-20">
          {[
            { icon: Truck, title: 'Fast Delivery', desc: 'Free shipping on orders above ₹50,000' },
            { icon: Shield, title: 'Warranty', desc: 'All products come with manufacturer warranty' },
            { icon: Zap, title: 'Best Prices', desc: 'Competitive pricing with regular discounts' },
          ].map((f) => (
            <div key={f.title} className="bg-white rounded-xl shadow-lg p-6 flex items-start gap-4 border border-gray-100">
              <div className="bg-green-50 text-green-600 p-3 rounded-xl">
                <f.icon size={24} />
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">{f.title}</h3>
                <p className="text-sm text-gray-500 mt-1">{f.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Featured Products */}
      <section className="max-w-7xl mx-auto px-4 pb-16">
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-3">
            <div className="bg-green-50 text-green-600 p-2 rounded-lg">
              <Cpu size={22} />
            </div>
            <h2 className="text-2xl font-bold text-gray-900">Featured Products</h2>
          </div>
          <Link to="/products" className="text-green-600 hover:text-green-700 font-medium text-sm flex items-center gap-1">
            View all <ArrowRight size={16} />
          </Link>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {products.map((p) => (
            <ProductCard key={p.id} product={p} />
          ))}
        </div>
        {products.length === 0 && (
          <p className="text-center text-gray-500 py-12">No products available yet.</p>
        )}
      </section>
    </div>
  )
}
