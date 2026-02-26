import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import toast from 'react-hot-toast'
import { ShoppingCart, Cpu, Star, ChevronLeft, ChevronRight } from 'lucide-react'
import api from '../../api/axios'
import { useAuth } from '../../context/AuthContext'

export default function ProductDetail() {
  const { id } = useParams()
  const { user } = useAuth()
  const [product, setProduct] = useState(null)
  const [loading, setLoading] = useState(true)
  const [activeImg, setActiveImg] = useState(0)
  const [reviews, setReviews] = useState([])
  const [avgRating, setAvgRating] = useState(0)
  const [reviewForm, setReviewForm] = useState({ rating: 5, comment: '' })
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    fetchProduct()
  }, [id])

  const fetchProduct = async () => {
    try {
      const res = await api.get(`/user/viewproducts?id=${id}`)
      if (res.data.data && res.data.data.length > 0) {
        setProduct(res.data.data[0])
        fetchReviews()
      }
    } catch {
      toast.error('Failed to load product')
    } finally {
      setLoading(false)
    }
  }

  const fetchReviews = async () => {
    try {
      const res = await api.get(`/reviews/${id}`)
      setReviews(res.data.reviews || [])
      setAvgRating(res.data.average_rating || 0)
    } catch {
      // reviews not critical
    }
  }

  const handleAddToCart = async () => {
    try {
      await api.post('/cart/add', { productid: product.id, quantity: 1 })
      toast.success('Added to cart!')
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to add to cart')
    }
  }

  const handleSubmitReview = async (e) => {
    e.preventDefault()
    setSubmitting(true)
    try {
      const userName = user?.user_name || user?.email || 'Anonymous'
      await api.post('/reviews', {
        product_id: id,
        rating: reviewForm.rating,
        comment: reviewForm.comment,
        user_name: userName,
      })
      toast.success('Review submitted!')
      setReviewForm({ rating: 5, comment: '' })
      fetchReviews()
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to submit review')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[60vh]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600" />
      </div>
    )
  }

  if (!product) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-16 text-center">
        <h2 className="text-2xl font-bold text-gray-700">Product not found</h2>
      </div>
    )
  }

  const images = [product.image1, product.image2, product.image3].filter(Boolean)
  const getImgSrc = (img) => img?.startsWith('http') ? img : `/public/images/${img}`

  const renderStars = (rating, size = 16) => (
    <div className="flex gap-0.5">
      {[1, 2, 3, 4, 5].map((i) => (
        <Star
          key={i}
          size={size}
          className={i <= rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}
        />
      ))}
    </div>
  )

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      {/* Product Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
        {/* Image Gallery */}
        <div>
          <div className="bg-white rounded-xl shadow-md overflow-hidden aspect-square flex items-center justify-center relative">
            {images.length > 0 ? (
              <>
                <img
                  src={getImgSrc(images[activeImg])}
                  alt={product.name}
                  className="h-full w-full object-cover"
                />
                {images.length > 1 && (
                  <>
                    <button
                      onClick={() => setActiveImg((prev) => (prev - 1 + images.length) % images.length)}
                      className="absolute left-2 top-1/2 -translate-y-1/2 bg-black/40 hover:bg-black/60 text-white rounded-full p-1"
                    >
                      <ChevronLeft size={20} />
                    </button>
                    <button
                      onClick={() => setActiveImg((prev) => (prev + 1) % images.length)}
                      className="absolute right-2 top-1/2 -translate-y-1/2 bg-black/40 hover:bg-black/60 text-white rounded-full p-1"
                    >
                      <ChevronRight size={20} />
                    </button>
                  </>
                )}
              </>
            ) : (
              <div className="flex flex-col items-center text-gray-400">
                <Cpu size={80} strokeWidth={1} />
                <span className="text-sm mt-2">No image available</span>
              </div>
            )}
          </div>
          {images.length > 1 && (
            <div className="flex gap-2 mt-3">
              {images.map((img, i) => (
                <button
                  key={i}
                  onClick={() => setActiveImg(i)}
                  className={`w-20 h-20 rounded-lg overflow-hidden border-2 ${activeImg === i ? 'border-green-500' : 'border-gray-200'}`}
                >
                  <img src={getImgSrc(img)} alt="" className="w-full h-full object-cover" />
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Product Info */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{product.name}</h1>
          <p className="text-gray-500 mt-1">{product.brand} &middot; Model #{product.model_no}</p>

          {avgRating > 0 && (
            <div className="flex items-center gap-2 mt-3">
              {renderStars(Math.round(avgRating))}
              <span className="text-sm text-gray-500">({avgRating.toFixed(1)}) &middot; {reviews.length} review{reviews.length !== 1 ? 's' : ''}</span>
            </div>
          )}

          <div className="mt-6 flex items-baseline gap-3">
            <span className="text-green-600 font-bold text-3xl">
              ₹{(product.price || 0).toLocaleString()}
            </span>
            {product.discount_amount > 0 && (
              <>
                <span className="text-gray-400 line-through text-lg">
                  ₹{(product.original_price || 0).toLocaleString()}
                </span>
                <span className="bg-red-100 text-red-600 text-sm font-semibold px-2 py-0.5 rounded">
                  -{product.discount_percentage}%
                </span>
              </>
            )}
          </div>

          <p className={`mt-3 text-sm font-medium ${product.stock > 0 ? 'text-green-600' : 'text-red-500'}`}>
            {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
          </p>

          {product.description && (
            <p className="mt-4 text-gray-600 leading-relaxed">{product.description}</p>
          )}

          {user && product.stock > 0 && (
            <button
              onClick={handleAddToCart}
              className="mt-6 flex items-center justify-center gap-2 bg-green-600 hover:bg-green-700 text-white py-3 px-8 rounded-lg transition font-semibold text-lg"
            >
              <ShoppingCart size={20} />
              Add to Cart
            </button>
          )}

          {/* Specifications */}
          {product.specifications && Object.keys(product.specifications).length > 0 && (
            <div className="mt-8">
              <h3 className="text-lg font-semibold text-gray-900 mb-3">Specifications</h3>
              <div className="bg-gray-50 rounded-lg overflow-hidden">
                <table className="w-full text-sm">
                  <tbody>
                    {Object.entries(product.specifications).map(([key, value], i) => (
                      <tr key={key} className={i % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>
                        <td className="py-2.5 px-4 font-medium text-gray-600 capitalize w-1/3">
                          {key.replace(/_/g, ' ')}
                        </td>
                        <td className="py-2.5 px-4 text-gray-900">{String(value)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Reviews Section */}
      <div className="mt-12">
        <h2 className="text-2xl font-bold text-gray-900 mb-6">Customer Reviews</h2>

        {/* Add Review Form */}
        {user && (
          <form onSubmit={handleSubmitReview} className="bg-white rounded-xl shadow-md p-6 mb-8">
            <h3 className="font-semibold text-gray-900 mb-4">Write a Review</h3>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">Rating</label>
              <div className="flex gap-1">
                {[1, 2, 3, 4, 5].map((i) => (
                  <button
                    key={i}
                    type="button"
                    onClick={() => setReviewForm((f) => ({ ...f, rating: i }))}
                  >
                    <Star
                      size={24}
                      className={`cursor-pointer ${i <= reviewForm.rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}`}
                    />
                  </button>
                ))}
              </div>
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">Comment</label>
              <textarea
                value={reviewForm.comment}
                onChange={(e) => setReviewForm((f) => ({ ...f, comment: e.target.value }))}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                rows={3}
                placeholder="Share your experience with this product..."
              />
            </div>
            <button
              type="submit"
              disabled={submitting}
              className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-medium text-sm disabled:opacity-50"
            >
              {submitting ? 'Submitting...' : 'Submit Review'}
            </button>
          </form>
        )}

        {/* Reviews List */}
        {reviews.length === 0 ? (
          <p className="text-gray-500">No reviews yet. Be the first to review this product!</p>
        ) : (
          <div className="space-y-4">
            {reviews.map((review) => (
              <div key={review.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-5">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-green-100 text-green-700 rounded-full flex items-center justify-center font-semibold text-sm">
                      {(review.user_name || 'U')[0].toUpperCase()}
                    </div>
                    <span className="font-medium text-gray-900">{review.user_name || 'Anonymous'}</span>
                  </div>
                  {renderStars(review.rating)}
                </div>
                {review.comment && (
                  <p className="text-gray-600 text-sm mt-2">{review.comment}</p>
                )}
                <p className="text-xs text-gray-400 mt-2">
                  {new Date(review.created_at).toLocaleDateString('en-US', {
                    year: 'numeric', month: 'long', day: 'numeric',
                  })}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
