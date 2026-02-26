import { Cpu } from 'lucide-react'

export default function CartItem({ item }) {
  return (
    <div className="flex items-center gap-4 bg-white p-4 rounded-xl shadow-sm border border-gray-100">
      <div className="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center flex-shrink-0">
        <Cpu size={28} className="text-gray-400" strokeWidth={1.5} />
      </div>
      <div className="flex-1 min-w-0">
        <h3 className="font-semibold text-gray-900 truncate">{item.product_name}</h3>
        <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
      </div>
      <span className="font-bold text-green-600 text-lg whitespace-nowrap">
        ₹{(item.product_price * item.quantity).toLocaleString()}
      </span>
    </div>
  )
}
