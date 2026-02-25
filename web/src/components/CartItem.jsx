export default function CartItem({ item }) {
  return (
    <div className="flex items-center justify-between bg-white p-4 rounded-lg shadow">
      <div>
        <h3 className="font-semibold">{item.product_name}</h3>
        <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
      </div>
      <span className="font-bold text-green-600">₹{item.product_price * item.quantity}</span>
    </div>
  )
}
