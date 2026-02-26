import toast from 'react-hot-toast'
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

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition">
      <div className="h-48 bg-gray-100 flex items-center justify-center">
        {product.image1 ? (
          <img
            src={`/public/images/${product.image1}`}
            alt={product.name}
            className="h-full w-full object-contain p-2"
          />
        ) : (
          <span className="text-gray-400 text-4xl">GPU</span>
        )}
      </div>
      <div className="p-4">
        <h3 className="font-semibold text-lg truncate">{product.name}</h3>
        <p className="text-sm text-gray-500">{product.brand}</p>
        <div className="mt-2 flex items-center gap-2">
          <span className="text-green-600 font-bold text-xl">
            ₹{discountedPrice}
          </span>
          {product.discount_id && (
            <span className="text-gray-400 line-through text-sm">₹{product.price}</span>
          )}
        </div>
        <p className="text-xs text-gray-500 mt-1">
          Stock: {product.stock > 0 ? product.stock : 'Out of stock'}
        </p>
        {user && product.stock > 0 && (
          <button
            onClick={handleAddToCart}
            className="mt-3 w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded transition"
          >
            Add to Cart
          </button>
        )}
      </div>
    </div>
  )
}
