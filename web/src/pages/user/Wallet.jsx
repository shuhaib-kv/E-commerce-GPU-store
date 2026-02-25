import { useEffect, useState } from 'react'
import api from '../../api/axios'

export default function Wallet() {
  const [wallet, setWallet] = useState(null)
  const [history, setHistory] = useState([])

  useEffect(() => {
    api.get('/user/wallet/history')
      .then((res) => {
        setWallet(res.data.wallet || null)
        setHistory(res.data.history || [])
      })
      .catch(() => {})
  }, [])

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">My Wallet</h1>

      <div className="bg-gradient-to-r from-green-600 to-green-700 text-white p-8 rounded-lg shadow-lg mb-8">
        <p className="text-sm opacity-80">Available Balance</p>
        <p className="text-4xl font-bold">₹{wallet?.balance || 0}</p>
      </div>

      <h2 className="text-xl font-semibold mb-4">Transaction History</h2>
      {history.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No transactions yet.</p>
      ) : (
        <div className="flex flex-col gap-3">
          {history.map((h) => (
            <div key={h._id} className="bg-white p-4 rounded-lg shadow flex justify-between">
              <span className={h.credit > 0 ? 'text-green-600' : 'text-red-600'}>
                {h.credit > 0 ? `+₹${h.credit} Credit` : `-₹${h.debit} Debit`}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
