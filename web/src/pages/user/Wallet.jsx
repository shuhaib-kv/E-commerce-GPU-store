import { useEffect, useState } from 'react'
import { Wallet as WalletIcon, ArrowUpCircle, ArrowDownCircle } from 'lucide-react'
import toast from 'react-hot-toast'
import api from '../../api/axios'

export default function Wallet() {
  const [balance, setBalance] = useState(0)
  const [history, setHistory] = useState([])

  useEffect(() => {
    api.get('/user/wallet/history')
      .then((res) => {
        setBalance(res.data.balance || 0)
        setHistory(res.data.history || [])
      })
      .catch((err) => toast.error(err.response?.data?.message || 'Failed to load wallet'))
  }, [])

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <div className="bg-green-50 text-green-600 p-2 rounded-lg">
          <WalletIcon size={22} />
        </div>
        <h1 className="text-3xl font-bold text-gray-900">My Wallet</h1>
      </div>

      <div className="bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white p-8 rounded-2xl shadow-xl mb-8 relative overflow-hidden">
        <div className="absolute top-4 right-4 opacity-10">
          <WalletIcon size={100} />
        </div>
        <p className="text-sm text-gray-400">Available Balance</p>
        <p className="text-5xl font-bold mt-2">₹{balance.toLocaleString()}</p>
      </div>

      <h2 className="text-xl font-semibold text-gray-900 mb-4">Transaction History</h2>
      {history.length === 0 ? (
        <div className="text-center py-16">
          <WalletIcon size={64} className="text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 text-lg">No transactions yet.</p>
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {history.map((h) => (
            <div key={h.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex items-center justify-between">
              <div className="flex items-center gap-3">
                {h.credit > 0 ? (
                  <div className="bg-green-50 text-green-600 p-2 rounded-lg">
                    <ArrowUpCircle size={20} />
                  </div>
                ) : (
                  <div className="bg-red-50 text-red-500 p-2 rounded-lg">
                    <ArrowDownCircle size={20} />
                  </div>
                )}
                <span className="font-medium text-gray-900">
                  {h.credit > 0 ? 'Credit' : 'Debit'}
                </span>
              </div>
              <span className={`font-bold text-lg ${h.credit > 0 ? 'text-green-600' : 'text-red-500'}`}>
                {h.credit > 0 ? `+₹${h.credit.toLocaleString()}` : `-₹${h.debit.toLocaleString()}`}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
