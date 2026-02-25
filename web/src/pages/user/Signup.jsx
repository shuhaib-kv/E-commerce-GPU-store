import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../../api/axios'
import { useAuth } from '../../context/AuthContext'

export default function Signup() {
  const [form, setForm] = useState({
    first_name: '', last_name: '', user_name: '', email: '', password: '', phone: '',
  })
  const [error, setError] = useState('')
  const { loginUser } = useAuth()
  const navigate = useNavigate()

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      const res = await api.post('/user/signup', form)
      loginUser(res.data.user || { email: form.email })
      navigate('/')
    } catch (err) {
      setError(err.response?.data?.error || 'Signup failed')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h1 className="text-2xl font-bold mb-6 text-center">Create Account</h1>
        {error && <p className="text-red-500 text-sm mb-4 text-center">{error}</p>}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-2 gap-4">
            <input name="first_name" placeholder="First Name" value={form.first_name} onChange={handleChange} className="border rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500" required />
            <input name="last_name" placeholder="Last Name" value={form.last_name} onChange={handleChange} className="border rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500" required />
          </div>
          <input name="user_name" placeholder="Username" value={form.user_name} onChange={handleChange} className="border rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="email" type="email" placeholder="Email" value={form.email} onChange={handleChange} className="border rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="phone" placeholder="Phone" value={form.phone} onChange={handleChange} className="border rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <input name="password" type="password" placeholder="Password" value={form.password} onChange={handleChange} className="border rounded px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500" required />
          <button className="bg-green-600 hover:bg-green-700 text-white py-2 rounded transition font-semibold">
            Sign Up
          </button>
        </form>
        <p className="mt-4 text-center text-sm text-gray-600">
          Already have an account? <Link to="/login" className="text-green-600 hover:underline">Login</Link>
        </p>
      </div>
    </div>
  )
}
