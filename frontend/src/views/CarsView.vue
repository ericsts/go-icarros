<template>
  <div class="page">
    <div class="page-header">
      <h1>Carros</h1>
      <button class="btn btn-primary" @click="openCreate">+ Cadastrar Carro</button>
    </div>

    <div v-if="loading" class="loading">Carregando carros...</div>
    <div v-else-if="error" class="alert alert-error">{{ error }}</div>

    <div v-else class="card">
      <div v-if="cars.length === 0" style="text-align:center;color:var(--text-light);padding:40px">
        Nenhum carro cadastrado ainda.
      </div>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Marca</th>
              <th>Modelo</th>
              <th>Ano</th>
              <th>Valor</th>
              <th>Dono</th>
              <th>Ações</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="car in cars" :key="car.id">
              <td style="color:var(--text-light)">#{{ car.id }}</td>
              <td><strong>{{ car.marca }}</strong></td>
              <td>{{ car.modelo }}</td>
              <td>{{ car.ano }}</td>
              <td>{{ formatCurrency(car.valor) }}</td>
              <td style="color:var(--text-light)">Usuário #{{ car.user_id }}</td>
              <td>
                <div class="actions" v-if="auth.isAdmin || car.user_id === auth.userId">
                  <button class="btn btn-ghost btn-sm" @click="openEdit(car)">Editar</button>
                  <button class="btn btn-danger  btn-sm" @click="handleDelete(car.id)">Excluir</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal criar/editar -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingId ? 'Editar Carro' : 'Cadastrar Carro' }}</h2>
          <button class="btn-close" @click="closeModal">×</button>
        </div>
        <div class="modal-body">
          <div v-if="formError" class="alert alert-error">{{ formError }}</div>
          <form @submit.prevent="save">
            <div class="form-row">
              <div class="form-group">
                <label>Marca</label>
                <input v-model="form.marca" required placeholder="Volkswagen" />
              </div>
              <div class="form-group">
                <label>Modelo</label>
                <input v-model="form.modelo" required placeholder="Golf GTI" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>Ano</label>
                <input v-model.number="form.ano" type="number" required min="1900" :max="new Date().getFullYear()+1" />
              </div>
              <div class="form-group">
                <label>Valor (R$)</label>
                <input v-model.number="form.valor" type="number" step="0.01" required />
              </div>
            </div>

            <template v-if="!editingId">
              <hr class="divider" />
              <p style="font-size:13px;color:var(--text-light);margin-bottom:12px">Configurações do leilão</p>
              <div class="form-row">
                <div class="form-group">
                  <label>Lance mínimo (R$)</label>
                  <input v-model.number="form.min_bid" type="number" step="0.01" required />
                </div>
                <div class="form-group">
                  <label>Encerramento do leilão</label>
                  <input v-model="form.auction_ends_at" type="datetime-local" required />
                </div>
              </div>
            </template>

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
import { getCars, createCar, updateCar, deleteCar } from '../api'
import { formatCurrency } from '../utils'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

const cars      = ref([])
const loading   = ref(true)
const error     = ref('')
const showModal = ref(false)
const editingId = ref(null)
const saving    = ref(false)
const formError = ref('')

const emptyForm = () => ({ marca: '', modelo: '', ano: new Date().getFullYear(), valor: '', min_bid: '', auction_ends_at: '' })
const form = reactive(emptyForm())

async function load() {
  try {
    const { data } = await getCars()
    cars.value = data ?? []
  } catch { error.value = 'Erro ao carregar carros' }
  finally   { loading.value = false }
}

function openCreate() {
  Object.assign(form, emptyForm())
  editingId.value = null
  formError.value = ''
  showModal.value = true
}

function openEdit(car) {
  Object.assign(form, { marca: car.marca, modelo: car.modelo, ano: car.ano, valor: car.valor })
  editingId.value = car.id
  formError.value = ''
  showModal.value = true
}

function closeModal() { showModal.value = false }

async function save() {
  formError.value = ''
  saving.value    = true
  try {
    if (editingId.value) {
      const { data } = await updateCar(editingId.value, { marca: form.marca, modelo: form.modelo, ano: form.ano, valor: form.valor })
      const i = cars.value.findIndex(c => c.id === editingId.value)
      if (i >= 0) cars.value[i] = data
    } else {
      const payload = {
        ...form,
        auction_ends_at: new Date(form.auction_ends_at).toISOString()
      }
      const { data } = await createCar(payload)
      cars.value.unshift(data)
    }
    closeModal()
  } catch (e) {
    formError.value = e.response?.data?.error ?? 'Erro ao salvar'
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  if (!confirm('Excluir este carro?')) return
  try {
    await deleteCar(id)
    cars.value = cars.value.filter(c => c.id !== id)
  } catch { alert('Erro ao excluir') }
}

onMounted(load)
</script>
