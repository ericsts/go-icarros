<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-logo">🚗</div>
      <h1>Criar conta</h1>
      <p class="auth-sub">Junte-se ao iCarros e comece a dar lances</p>

      <div v-if="error" class="alert alert-error">{{ error }}</div>

      <form @submit.prevent="submit">
        <div class="form-group">
          <label>Nome completo</label>
          <input v-model="form.name" required placeholder="Seu nome" autocomplete="name" />
        </div>
        <div class="form-group">
          <label>E-mail</label>
          <input v-model="form.email" type="email" required placeholder="seu@email.com" autocomplete="email" />
        </div>
        <div class="form-group">
          <label>Senha</label>
          <input v-model="form.password" type="password" required minlength="6" placeholder="Mínimo 6 caracteres" autocomplete="new-password" />
        </div>

        <button type="submit" class="btn btn-primary" style="width:100%;justify-content:center" :disabled="loading">
          {{ loading ? 'Criando conta...' : 'Criar conta' }}
        </button>
      </form>

      <p class="auth-footer">
        Já tem conta? <router-link to="/login">Entrar</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api'

const router  = useRouter()
const loading = ref(false)
const error   = ref('')
const form    = reactive({ name: '', email: '', password: '' })

async function submit() {
  error.value   = ''
  loading.value = true
  try {
    await register(form)
    router.push({ path: '/login', query: { registered: '1' } })
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Erro ao criar conta'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0f172a;
  padding: 24px;
}
.auth-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 40px 36px;
  width: 100%;
  max-width: 400px;
}
.auth-logo {
  font-size: 36px;
  text-align: center;
  margin-bottom: 8px;
}
h1 {
  font-size: 22px;
  font-weight: 700;
  color: #f1f5f9;
  text-align: center;
  margin: 0 0 4px;
}
.auth-sub {
  text-align: center;
  color: #64748b;
  font-size: 13px;
  margin: 0 0 24px;
}
.form-group { margin-bottom: 16px; }
.form-group label { display: block; font-size: 13px; color: #94a3b8; margin-bottom: 6px; }
.form-group input {
  width: 100%;
  background: #0f172a;
  border: 1px solid #334155;
  color: #f1f5f9;
  border-radius: 6px;
  padding: 9px 12px;
  font-size: 14px;
  box-sizing: border-box;
  transition: border-color .15s;
}
.form-group input:focus { outline: none; border-color: var(--primary, #2563eb); }
.auth-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 13px;
  color: #64748b;
}
.auth-footer a { color: var(--primary, #2563eb); }
</style>
