import { Routes, Route, Outlet } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import ErrorBoundary from './components/ErrorBoundary'
import Navbar from './components/Navbar'
import AdminSidebar from './components/AdminSidebar'
import ProtectedRoute from './components/ProtectedRoute'
import AdminRoute from './components/AdminRoute'

import Home from './pages/user/Home'
import Login from './pages/user/Login'
import Signup from './pages/user/Signup'
import Products from './pages/user/Products'
import Cart from './pages/user/Cart'
import Orders from './pages/user/Orders'
import Address from './pages/user/Address'
import Wallet from './pages/user/Wallet'
import ProductDetail from './pages/user/ProductDetail'

import AdminLogin from './pages/admin/Login'
import Dashboard from './pages/admin/Dashboard'
import Users from './pages/admin/Users'
import AdminProducts from './pages/admin/Products'
import Categories from './pages/admin/Categories'
import AdminOrders from './pages/admin/Orders'
import Coupons from './pages/admin/Coupons'
import Discounts from './pages/admin/Discounts'

function UserLayout() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <Outlet />
    </div>
  )
}

function AdminLayout() {
  return (
    <div className="flex min-h-screen bg-gray-50">
      <AdminSidebar />
      <main className="flex-1 p-8 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}

export default function App() {
  return (
    <ErrorBoundary>
      <Toaster position="top-right" toastOptions={{ duration: 4000 }} />
      <Routes>
        {/* User routes */}
        <Route element={<UserLayout />}>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/products" element={<Products />} />
          <Route path="/products/:id" element={<ProductDetail />} />
          <Route path="/cart" element={<ProtectedRoute><Cart /></ProtectedRoute>} />
          <Route path="/orders" element={<ProtectedRoute><Orders /></ProtectedRoute>} />
          <Route path="/address" element={<ProtectedRoute><Address /></ProtectedRoute>} />
          <Route path="/wallet" element={<ProtectedRoute><Wallet /></ProtectedRoute>} />
        </Route>

        {/* Admin routes */}
        <Route path="/admin/login" element={<AdminLogin />} />
        <Route path="/admin" element={<AdminRoute><AdminLayout /></AdminRoute>}>
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="users" element={<Users />} />
          <Route path="products" element={<AdminProducts />} />
          <Route path="categories" element={<Categories />} />
          <Route path="orders" element={<AdminOrders />} />
          <Route path="coupons" element={<Coupons />} />
          <Route path="discounts" element={<Discounts />} />
        </Route>
      </Routes>
    </ErrorBoundary>
  )
}
