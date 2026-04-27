<template>
  <div class="page">
    <div class="page-header">
      <h1>Usuários</h1>
      <button class="btn btn-primary" @click="openCreate">+ Novo Usuário</button>
    </div>

    <div v-if="loading" class="loading">Carregando usuários...</div>
    <div v-else-if="error" class="alert alert-error">{{ error }}</div>

    <div v-else class="card">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Nome</th>
              <th>E-mail</th>
              <th>Perfil</th>
              <th>Ações</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td style="color:var(--text-light)">#{{ u.id }}</td>
              <td><strong>{{ u.name }}</strong></td>
              <td>{{ u.email }}</td>
              <td>
                <span class="badge" :class="u.role === 'admin' ? 'badge-info' : 'badge-gray'">
                  {{ u.role }}
                </span>
              </td>
              <td>
                <div class="actions">
                  <button class="btn btn-ghost  btn-sm" @click="openEdit(u)">Editar</button>
                  <button class="btn btn-danger btn-sm" @click="handleDelete(u.id)">Excluir</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingId ? 'Editar Usuário' : 'Novo Usuário' }}</h2>
          <button class="btn-close" @click="closeModal">×</button>
        </div>
        <div class="modal-body">
          <div v-if="formError" class="alert alert-error">{{ formError }}</div>
          <form @submit.prevent="save">
            <div class="form-group">
              <label>Nome</label>
              <input v-model="form.name" required placeholder="João Silva" />
            </div>
            <div class="form-group">
              <label>E-mail</label>
              <input v-model="form.email" type="email" required placeholder="joao@email.com" />
            </div>
            <div class="form-group" v-if="!editingId">
              <label>Senha</label>
              <input v-model="form.password" type="password" required placeholder="••••••••" />
            </div>
            <div class="form-group">
              <label>Perfil</label>
              <select v-model="form.role" required>
                <option value="user">Usuário</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-ghost" @click="closeModal">Cancelar</button>
              <button type="submit" class="btn btn-primary" :disabled="saving">
                {{ saving ? 'Salvando...' : 'Salvar' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getUsers, createUser, updateUser, deleteUser } from '../api'

const users     = ref([])
const loading   = ref(true)
const error     = ref('')
const showModal = ref(false)
const editingId = ref(null)
const saving    = ref(false)
const formError = ref('')
const form      = reactive({ name: '', email: '', password: '', role: 'user' })

async function load() {
  try {
    const { data } = await getUsers()
    users.value = data ?? []
  } catch { error.value = 'Erro ao carregar usuários' }
  finally   { loading.value = false }
}

function openCreate() {
  Object.assign(form, { name: '', email: '', password: '', role: 'user' })
  editingId.value = null
  formError.value = ''
  showModal.value = true
}

function openEdit(u) {
  Object.assign(form, { name: u.name, email: u.email, password: '', role: u.role })
  editingId.value = u.id
  formError.value = ''
  showModal.value = true
}

function closeModal() { showModal.value = false }

async function save() {
  formError.value = ''
  saving.value    = true
  try {
    if (editingId.value) {
      const { data } = await updateUser(editingId.value, { name: form.name, email: form.email, role: form.role })
      const i = users.value.findIndex(u => u.id === editingId.value)
      if (i >= 0) users.value[i] = data
    } else {
      const { data } = await createUser({ name: form.name, email: form.email, password: form.password, role: form.role })
      users.value.unshift(data)
    }
    closeModal()
  } catch (e) {
    formError.value = e.response?.data?.error ?? 'Erro ao salvar'
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  if (!confirm('Excluir este usuário?')) return
  try {
    await deleteUser(id)
    users.value = users.value.filter(u => u.id !== id)
  } catch { alert('Erro ao excluir') }
}

onMounted(load)
</script>
