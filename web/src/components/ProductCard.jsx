import toast from 'react-hot-toast'
import { ShoppingCart, Cpu } from 'lucide-react'
import api from '../api/axios'
import { useAuth } from '../context/AuthContext'

export default function ProductCard({ product, onAddToCart }) {
  const { user } = useAuth()

  const handleAddToCart = async () => {
    try {
      await api.post('/cart/add', { productid: product.id, quantity: 1 })
      if (onAddToCart) onAddToCart()
      toast.success('Added to cart!')
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to add to cart')
    }
  }

  const discountedPrice = product.discount_id
    ? Math.round(product.price * (1 - (product.discount_percentage || 0) / 100))
    : product.price

  const imgSrc = product.image1
    ? product.image1.startsWith('http')
      ? product.image1
      : `/public/images/${product.image1}`
    : null

  return (
    <div className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-xl transition group border border-gray-100">
      <div className="h-48 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center relative overflow-hidden">
        {imgSrc ? (
          <img
            src={imgSrc}
            alt={product.name}
            className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="flex flex-col items-center text-gray-400">
            <Cpu size={48} strokeWidth={1} />
            <span className="text-xs mt-1">GPU</span>
          </div>
        )}
        {product.discount_id && product.discount_percentage > 0 && (
          <span className="absolute top-2 right-2 bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded-full">
            -{product.discount_percentage}%
          </span>
        )}
      </div>
      <div className="p-4">
        <h3 className="font-semibold text-gray-900 truncate">{product.name}</h3>
        <p className="text-sm text-gray-500 mt-0.5">{product.brand}</p>
        <div className="mt-3 flex items-baseline gap-2">
          <span className="text-green-600 font-bold text-xl">₹{discountedPrice.toLocaleString()}</span>
          {product.discount_id && (
            <span className="text-gray-400 line-through text-sm">₹{product.price.toLocaleString()}</span>
          )}
        </div>
        <p className={`text-xs mt-1 ${product.stock > 0 ? 'text-gray-500' : 'text-red-500 font-medium'}`}>
          {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
        </p>
        {user && product.stock > 0 && (
          <button
            onClick={handleAddToCart}
            className="mt-3 w-full flex items-center justify-center gap-2 bg-green-600 hover:bg-green-700 text-white py-2 rounded-lg transition font-medium text-sm"
          >
            <ShoppingCart size={16} />
            Add to Cart
          </button>
        )}
      </div>
    </div>
  )
}
