<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-logo">🚗 iCarros</div>
      <h1>Entrar na plataforma</h1>

      <div v-if="registered" class="alert alert-success">Conta criada com sucesso! Faça login.</div>
      <div v-if="error" class="alert alert-error">{{ error }}</div>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>E-mail</label>
          <input v-model="form.email" type="email" placeholder="seu@email.com" required autofocus />
        </div>
        <div class="form-group">
          <label>Senha</label>
          <input v-model="form.password" type="password" placeholder="••••••••" required />
        </div>
        <button class="btn btn-primary" style="width:100%;justify-content:center" :disabled="loading">
          {{ loading ? 'Entrando...' : 'Entrar' }}
        </button>
      </form>

      <p style="text-align:center;margin-top:20px;font-size:13px;color:#94a3b8">
        Não tem conta? <router-link to="/register" style="color:#2563eb">Criar conta</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter, useRoute } from 'vue-router'

const auth   = useAuthStore()
const router = useRouter()
const route  = useRoute()

const form       = reactive({ email: '', password: '' })
const loading    = ref(false)
const error      = ref('')
const registered = computed(() => route.query.registered === '1')

async function handleLogin() {
  error.value   = ''
  loading.value = true
  try {
    await auth.login(form.email, form.password)
    router.push('/auctions')
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Credenciais inválidas'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
  padding: 24px;
}
.login-card {
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  width: 100%;
  max-width: 380px;
  box-shadow: 0 24px 80px rgba(0,0,0,.3);
}
.alert-success {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
  border-radius: 6px;
  padding: 10px 14px;
  font-size: 13px;
  margin-bottom: 12px;
}
.login-logo {
  font-size: 28px;
  font-weight: 800;
  text-align: center;
  margin-bottom: 8px;
}
h1 {
  text-align: center;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-light);
  margin-bottom: 28px;
}
</style>
