import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getMe } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token     = ref(localStorage.getItem('token')     || '')
  const userId    = ref(Number(localStorage.getItem('userId') || '0'))
  const userName  = ref(localStorage.getItem('userName')  || '')
  const userEmail = ref(localStorage.getItem('userEmail') || '')
  const userRole  = ref(localStorage.getItem('userRole')  || '')

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin         = computed(() => userRole.value === 'admin')

  async function login(email, password) {
    const { data } = await loginApi(email, password)
    token.value     = data.token
    userId.value    = data.id
    userName.value  = data.name
    userEmail.value = data.email
    userRole.value  = data.role
    localStorage.setItem('token',     data.token)
    localStorage.setItem('userId',    data.id)
    localStorage.setItem('userName',  data.name)
    localStorage.setItem('userEmail', data.email)
    localStorage.setItem('userRole',  data.role)
  }

  async function fetchMe() {
    try {
      const { data } = await getMe()
      userId.value    = data.id
      userName.value  = data.name
      userEmail.value = data.email
      userRole.value  = data.role
      localStorage.setItem('userId',    data.id)
      localStorage.setItem('userName',  data.name)
      localStorage.setItem('userEmail', data.email)
      localStorage.setItem('userRole',  data.role)
    } catch { /* token inválido — o interceptor do axios já redireciona */ }
  }

  function setProfile({ name, email }) {
    userName.value  = name
    userEmail.value = email
    localStorage.setItem('userName',  name)
    localStorage.setItem('userEmail', email)
  }

  function logout() {
    token.value     = ''
    userId.value    = 0
    userName.value  = ''
    userEmail.value = ''
    userRole.value  = ''
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('userName')
    localStorage.removeItem('userEmail')
    localStorage.removeItem('userRole')
  }

  return { token, userId, isAuthenticated, isAdmin, userName, userEmail, login, logout, setProfile, fetchMe }
})
