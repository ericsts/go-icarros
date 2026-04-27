<template>
  <div class="page">
    <div class="page-header">
      <h1>Logs do sistema</h1>
    </div>

    <!-- Filtros -->
    <div class="card" style="margin-bottom:16px">
      <form class="filters" @submit.prevent="load">
        <div class="form-group" style="margin:0;flex:0 0 150px">
          <label>Level</label>
          <select v-model="filters.level">
            <option value="">Todos</option>
            <option value="info">Info</option>
            <option value="warn">Warn</option>
            <option value="error">Error</option>
          </select>
        </div>
        <div class="form-group" style="margin:0;flex:1">
          <label>Evento</label>
          <input v-model="filters.event" placeholder="ex: car.created" />
        </div>
        <div class="form-group" style="margin:0;flex:0 0 100px">
          <label>Limite</label>
          <input v-model.number="filters.limit" type="number" min="1" max="500" />
        </div>
        <div style="display:flex;align-items:flex-end">
          <button type="submit" class="btn btn-primary">Filtrar</button>
        </div>
      </form>
    </div>

    <div v-if="loading" class="loading">Carregando logs...</div>
    <div v-else-if="error" class="alert alert-error">{{ error }}</div>

    <div v-else class="card">
      <div v-if="logs.length === 0" style="text-align:center;color:var(--text-light);padding:40px">
        Nenhum log encontrado com esses filtros.
      </div>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Data/Hora</th>
              <th>Level</th>
              <th>Evento</th>
              <th>Mensagem</th>
              <th>Metadata</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id">
              <td style="white-space:nowrap;color:var(--text-light)">{{ formatDate(log.created_at) }}</td>
              <td>
                <span class="badge" :class="levelBadge(log.level)">{{ log.level }}</span>
              </td>
              <td style="font-family:monospace;font-size:12px">{{ log.event }}</td>
              <td>{{ log.message }}</td>
              <td>
                <details v-if="log.metadata && Object.keys(log.metadata).length">
                  <summary style="cursor:pointer;font-size:12px;color:var(--primary)">Ver</summary>
                  <pre class="meta-pre">{{ JSON.stringify(log.metadata, null, 2) }}</pre>
                </details>
                <span v-else style="color:var(--text-light)">—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div style="color:var(--text-light);font-size:12px;padding:12px 0 0;text-align:right">
        {{ logs.length }} registro(s)
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getLogs } from '../api'
import { formatDate } from '../utils'

const logs    = ref([])
const loading = ref(true)
const error   = ref('')
const filters = reactive({ level: '', event: '', limit: 100 })

function levelBadge(l) {
  return { info: 'badge-info', warn: 'badge-warning', error: 'badge-danger' }[l] ?? 'badge-gray'
}

async function load() {
  loading.value = true
  error.value   = ''
  try {
    const params = {}
    if (filters.level) params.level = filters.level
    if (filters.event) params.event = filters.event
    if (filters.limit) params.limit = filters.limit
    const { data } = await getLogs(params)
    logs.value = data ?? []
  } catch { error.value = 'Erro ao carregar logs' }
  finally   { loading.value = false }
}

onMounted(load)
</script>

<style scoped>
.filters { display: flex; gap: 16px; align-items: flex-start; flex-wrap: wrap; }
.meta-pre {
  margin-top: 6px;
  padding: 8px;
  background: #f8fafc;
  border: 1px solid var(--border);
  border-radius: 4px;
  font-size: 11px;
  max-width: 300px;
  overflow-x: auto;
  white-space: pre-wrap;
}
</style>
