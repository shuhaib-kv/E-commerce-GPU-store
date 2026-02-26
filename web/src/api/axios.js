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
      const isAdmin = localStorage.getItem('admin')
      const isUser = localStorage.getItem('user')
      if (isAdmin || isUser) {
        localStorage.removeItem('user')
        localStorage.removeItem('admin')
        if (!window.location.pathname.includes('/login')) {
          window.location.href = isAdmin && window.location.pathname.startsWith('/admin')
            ? '/admin/login'
            : '/login'
        }
      }
    } else if (err.response.status >= 500) {
      toast.error('Server error. Please try again later.')
    }

    return Promise.reject(err)
  }
)

export default api
