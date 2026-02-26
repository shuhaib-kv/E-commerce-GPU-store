import axios from 'axios'
import toast from 'react-hot-toast'

const api = axios.create({
  baseURL: '',
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (!err.response) {
      toast.error('Network error. Please check your connection.')
      return Promise.reject(err)
    }

    if (err.response.status === 401) {
      const isLoggedIn = localStorage.getItem('user') || localStorage.getItem('admin')
      if (isLoggedIn) {
        localStorage.removeItem('user')
        localStorage.removeItem('admin')
        if (!window.location.pathname.includes('/login')) {
          window.location.href = '/login'
        }
      }
    } else if (err.response.status >= 500) {
      toast.error('Server error. Please try again later.')
    }

    return Promise.reject(err)
  }
)

export default api
