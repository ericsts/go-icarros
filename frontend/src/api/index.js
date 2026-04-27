import axios from 'axios'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

// Auth
export const login = (email, password) =>
  api.post('/login', { email, password })

// Profile (authenticated user)
export const getMe    = ()     => api.get('/me')
export const updateMe = data   => api.put('/me', data)

// Users
export const getUsers  = ()         => api.get('/users')
export const createUser = data      => api.post('/users', data)
export const updateUser = (id, data)=> api.put(`/users/${id}`, data)
export const deleteUser = id        => api.delete(`/users/${id}`)

// Cars
export const getCars    = ()        => api.get('/cars')
export const getMyCars  = ()        => api.get('/cars/my')
export const getCar     = id        => api.get(`/cars/${id}`)
export const createCar  = data      => api.post('/cars', data)
export const updateCar  = (id,data) => api.put(`/cars/${id}`, data)
export const deleteCar  = id        => api.delete(`/cars/${id}`)

// Auctions
export const getAuctions = ()       => api.get('/auctions')
export const getAuction  = id       => api.get(`/auctions/${id}`)
export const placeBid    = (id,amt) => api.post(`/auctions/${id}/bids`, { amount: amt })
export const getBids     = id       => api.get(`/auctions/${id}/bids`)

// Logs
export const getLogs = params => api.get('/logs', { params })
