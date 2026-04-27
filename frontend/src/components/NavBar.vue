<template>
  <nav class="navbar">
    <div class="nav-inner">
      <div class="nav-brand">
        <router-link to="/auctions">🚗 iCarros</router-link>
      </div>

      <div class="nav-links">
        <router-link to="/auctions" active-class="active">Leilões</router-link>
        <router-link to="/cars"     active-class="active">Carros</router-link>
        <template v-if="auth.isAdmin">
          <span class="nav-sep">|</span>
          <router-link to="/admin/users" active-class="active">Usuários</router-link>
          <router-link to="/admin/logs"  active-class="active">Logs</router-link>
        </template>
      </div>

      <div class="nav-right">
        <button class="nav-user" @click="openProfile" :title="auth.userEmail">
          <span class="nav-avatar">{{ initials }}</span>
          <span class="nav-name">{{ auth.userName }}</span>
        </button>
        <button class="btn btn-ghost btn-sm" @click="handleLogout">Sair</button>
      </div>
    </div>
  </nav>

  <!-- Modal de perfil -->
  <div v-if="showProfile" class="modal-overlay" @click.self="closeProfile">
    <div class="modal modal-sm">
      <div class="modal-header">
        <h2>Meu Perfil</h2>
        <button class="btn-close" @click="closeProfile">×</button>
      </div>
      <div class="modal-body">
        <div v-if="profileError" class="alert alert-error">{{ profileError }}</div>
        <div v-if="profileSuccess" class="alert alert-success">{{ profileSuccess }}</div>
        <form @submit.prevent="saveProfile">
          <div class="form-group">
            <label>Nome</label>
            <input v-model="form.name" required placeholder="Seu nome" />
          </div>
          <div class="form-group">
            <label>E-mail</label>
            <input v-model="form.email" type="email" required placeholder="seu@email.com" />
          </div>
          <div class="form-group">
            <label>Nova senha <span class="label-hint">(deixe em branco para não alterar)</span></label>
            <input v-model="form.password" type="password" placeholder="••••••••" autocomplete="new-password" />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-ghost" @click="closeProfile">Cancelar</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? 'Salvando...' : 'Salvar' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { updateMe } from '../api'

const auth   = useAuthStore()
const router = useRouter()

const showProfile  = ref(false)
const saving       = ref(false)
const profileError = ref('')
const profileSuccess = ref('')

const form = reactive({ name: '', email: '', password: '' })

const initials = computed(() => {
  const n = auth.userName || ''
  return n.split(' ').map(p => p[0]).slice(0, 2).join('').toUpperCase() || '?'
})

function openProfile() {
  form.name     = auth.userName
  form.email    = auth.userEmail
  form.password = ''
  profileError.value   = ''
  profileSuccess.value = ''
  showProfile.value = true
}

function closeProfile() {
  showProfile.value = false
}

async function saveProfile() {
  profileError.value   = ''
  profileSuccess.value = ''
  saving.value = true
  try {
    const payload = { name: form.name, email: form.email }
    if (form.password) payload.password = form.password
    await updateMe(payload)
    auth.setProfile({ name: form.name, email: form.email })
    profileSuccess.value = 'Perfil atualizado com sucesso!'
    form.password = ''
  } catch (e) {
    profileError.value = e.response?.data?.error ?? 'Erro ao salvar perfil'
  } finally {
    saving.value = false
  }
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  background: #1e293b;
  color: #e2e8f0;
  position: sticky;
  top: 0;
  z-index: 50;
  border-bottom: 1px solid #334155;
}
.nav-inner {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0 16px;
  height: 52px;
  display: flex;
  align-items: center;
  gap: 24px;
}
.nav-brand a {
  font-weight: 700;
  font-size: 16px;
  color: #fff;
  white-space: nowrap;
}
.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}
.nav-links a {
  color: #94a3b8;
  font-size: 14px;
  font-weight: 500;
  padding: 4px 0;
  border-bottom: 2px solid transparent;
  transition: color .15s;
}
.nav-links a:hover, .nav-links a.active {
  color: #fff;
  border-bottom-color: var(--primary);
}
.nav-sep { color: #334155; }
.nav-right { display: flex; align-items: center; gap: 10px; margin-left: auto; }

.nav-user {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background .15s;
  color: #e2e8f0;
}
.nav-user:hover { background: #334155; }

.nav-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--primary, #2563eb);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.nav-name {
  font-size: 13px;
  font-weight: 500;
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.label-hint {
  font-size: 11px;
  font-weight: 400;
  color: var(--text-light);
}

.modal-sm { max-width: 420px; }

.alert-success {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
  border-radius: 6px;
  padding: 10px 14px;
  font-size: 13px;
  margin-bottom: 12px;
}
</style>
