<template>
  <div class="page">
    <div class="page-header">
      <h1>Leilões</h1>
    </div>

    <div v-if="loading" class="loading">Carregando leilões...</div>
    <div v-else-if="error" class="alert alert-error">{{ error }}</div>

    <div v-else-if="auctions.length === 0" class="card" style="text-align:center;color:var(--text-light);padding:40px">
      Nenhum leilão encontrado. Cadastre um carro para iniciar um leilão.
    </div>

    <div v-else class="cards-grid">
      <div
        v-for="a in auctions"
        :key="a.id"
        class="auction-card card"
        @click="$router.push(`/auctions/${a.id}`)"
      >
        <div class="auction-card-top">
          <span class="badge" :class="statusBadge(a.status)">{{ a.status === 'open' ? 'Aberto' : 'Encerrado' }}</span>
          <span class="auction-id">#{{ a.id }}</span>
        </div>

        <div class="auction-car">🚗 Carro #{{ a.car_id }}</div>

        <div class="auction-stats">
          <div class="stat">
            <span class="stat-label">Lance mínimo</span>
            <span class="stat-value">{{ formatCurrency(a.min_bid) }}</span>
          </div>
          <div class="stat">
            <span class="stat-label">Lance atual</span>
            <span class="stat-value highlight">{{ a.current_bid ? formatCurrency(a.current_bid) : '—' }}</span>
          </div>
          <div class="stat">
            <span class="stat-label">Total de lances</span>
            <span class="stat-value">{{ a.total_bids ?? 0 }}</span>
          </div>
          <div class="stat">
            <span class="stat-label">{{ a.status === 'open' ? 'Encerra em' : 'Encerrou em' }}</span>
            <span class="stat-value" :class="a.status === 'open' ? 'timer' : ''">
              {{ a.status === 'open' ? timeLeft(a.ends_at) : formatDate(a.ends_at) }}
            </span>
          </div>
        </div>

        <button class="btn btn-primary" style="width:100%;justify-content:center;margin-top:16px">
          {{ a.status === 'open' ? 'Dar lance' : 'Ver detalhes' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAuctions } from '../api'
import { formatCurrency, formatDate, timeLeft } from '../utils'

const auctions = ref([])
const loading  = ref(true)
const error    = ref('')

function statusBadge(s) { return s === 'open' ? 'badge-success' : 'badge-gray' }

onMounted(async () => {
  try {
    const { data } = await getAuctions()
    auctions.value = data ?? []
  } catch (e) {
    error.value = 'Erro ao carregar leilões'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.auction-card {
  cursor: pointer;
  transition: transform .15s, box-shadow .15s;
}
.auction-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,.12); }
.auction-card-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.auction-id { font-size: 12px; color: var(--text-light); }
.auction-car { font-weight: 600; font-size: 15px; margin-bottom: 16px; }
.auction-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.stat { display: flex; flex-direction: column; gap: 2px; }
.stat-label { font-size: 11px; color: var(--text-light); text-transform: uppercase; letter-spacing: .04em; }
.stat-value { font-size: 14px; font-weight: 600; }
.stat-value.highlight { color: var(--primary); }
.stat-value.timer { color: var(--warning); }
</style>
