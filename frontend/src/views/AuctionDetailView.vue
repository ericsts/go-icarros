<template>
  <div class="page">
    <div v-if="loading" class="loading">Carregando leilão...</div>
    <div v-else-if="error" class="alert alert-error">{{ error }}</div>

    <template v-else>
      <!-- Cabeçalho -->
      <div class="page-header">
        <div>
          <button class="btn btn-ghost btn-sm" style="margin-bottom:8px" @click="$router.back()">← Voltar</button>
          <h1>
            Leilão #{{ auction.id }}
            <span v-if="car">— {{ car.marca }} {{ car.modelo }} {{ car.ano }}</span>
          </h1>
        </div>
        <span class="badge" :class="auction.status === 'open' ? 'badge-success' : 'badge-gray'">
          {{ auction.status === 'open' ? 'Aberto' : 'Encerrado' }}
        </span>
      </div>

      <!-- Stats -->
      <div class="stats-row card" style="margin-bottom:20px">
        <div class="stat-box">
          <span class="label">Lance mínimo</span>
          <span class="value">{{ formatCurrency(auction.min_bid) }}</span>
        </div>
        <div class="stat-box highlight">
          <span class="label">Lance atual</span>
          <span class="value">{{ auction.current_bid ? formatCurrency(auction.current_bid) : '—' }}</span>
        </div>
        <div class="stat-box">
          <span class="label">Total de lances</span>
          <span class="value">{{ auction.total_bids ?? 0 }}</span>
        </div>
        <div class="stat-box" :class="auction.status === 'open' ? 'timer' : ''">
          <span class="label">{{ auction.status === 'open' ? 'Encerra em' : 'Encerrado em' }}</span>
          <span class="value">{{ auction.status === 'open' ? timeLeft(auction.ends_at) : formatDate(auction.ends_at) }}</span>
        </div>
        <div v-if="car" class="stat-box">
          <span class="label">Valor original</span>
          <span class="value">{{ formatCurrency(car.valor) }}</span>
        </div>
      </div>

      <div class="detail-grid">
        <!-- Coluna esquerda: dar lance + feed ao vivo -->
        <div>
          <!-- Formulário de lance -->
          <div v-if="isActive" class="card" style="margin-bottom:16px">
            <h2 style="font-size:16px;font-weight:600;margin-bottom:16px">Dar um lance</h2>
            <div v-if="bidError" class="alert alert-error">{{ bidError }}</div>
            <div v-if="bidSuccess" class="alert alert-success">Lance enviado com sucesso!</div>
            <form @submit.prevent="submitBid">
              <div class="form-group">
                <label>Valor do lance</label>
                <input
                  v-model.number="bidAmount"
                  type="number"
                  step="0.01"
                  :min="minBid + 0.01"
                  :placeholder="`Mínimo: ${formatCurrency(minBid)}`"
                  required
                />
              </div>
              <button class="btn btn-primary" style="width:100%;justify-content:center" :disabled="bidLoading">
                {{ bidLoading ? 'Enviando...' : 'Confirmar lance' }}
              </button>
            </form>
          </div>

          <!-- Feed ao vivo -->
          <div class="card">
            <div class="feed-header">
              <h2 style="font-size:16px;font-weight:600">Feed ao vivo</h2>
              <span class="live-dot" :class="wsConnected ? 'connected' : ''">
                {{ wsConnected ? '● Conectado' : '○ Desconectado' }}
              </span>
            </div>

            <div v-if="liveFeed.length === 0" style="color:var(--text-light);padding:20px 0;text-align:center;font-size:13px">
              Aguardando novos lances...
            </div>
            <div v-else class="feed-list">
              <div v-for="(bid, i) in liveFeed" :key="i" class="feed-item" :class="i === 0 ? 'new' : ''">
                <div class="feed-amount">{{ formatCurrency(bid.amount) }}</div>
                <div class="feed-meta">Usuário #{{ bid.user_id }} · {{ formatDate(bid.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Coluna direita: histórico completo -->
        <div class="card">
          <h2 style="font-size:16px;font-weight:600;margin-bottom:16px">Histórico de lances</h2>
          <div v-if="bids.length === 0" style="color:var(--text-light);text-align:center;padding:20px 0;font-size:13px">
            Nenhum lance ainda
          </div>
          <div v-else class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>#</th>
                  <th>Usuário</th>
                  <th>Valor</th>
                  <th>Data</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(bid, i) in bids" :key="bid.id" :class="i === 0 ? 'top-bid' : ''">
                  <td>{{ i === 0 ? '🏆' : i + 1 }}</td>
                  <td>#{{ bid.user_id }}</td>
                  <td><strong>{{ formatCurrency(bid.amount) }}</strong></td>
                  <td style="color:var(--text-light)">{{ formatDate(bid.created_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { getAuction, getBids, getCar, placeBid } from '../api'
import { formatCurrency, formatDate, timeLeft } from '../utils'

const route = useRoute()
const auth  = useAuthStore()
const id    = route.params.id

const auction    = ref(null)
const car        = ref(null)
const bids       = ref([])
const loading    = ref(true)
const error      = ref('')
const bidAmount  = ref('')
const bidLoading = ref(false)
const bidError   = ref('')
const bidSuccess = ref(false)
const liveFeed   = ref([])
const wsConnected= ref(false)
let   ws         = null

const minBid = computed(() => {
  return auction.value?.current_bid > 0
    ? auction.value.current_bid
    : auction.value?.min_bid ?? 0
})

const isActive = computed(() =>
  auction.value?.status === 'open' && new Date(auction.value.ends_at) > new Date()
)

async function loadData() {
  try {
    const [aRes, bRes] = await Promise.all([getAuction(id), getBids(id)])
    auction.value = aRes.data
    bids.value    = bRes.data ?? []

    const cRes = await getCar(auction.value.car_id)
    car.value = cRes.data
  } catch (e) {
    error.value = 'Erro ao carregar leilão'
  } finally {
    loading.value = false
  }
}

async function submitBid() {
  bidError.value   = ''
  bidSuccess.value = false
  bidLoading.value = true
  try {
    const { data } = await placeBid(id, bidAmount.value)
    bids.value.unshift(data)
    auction.value.current_bid = data.amount
    auction.value.total_bids = (auction.value.total_bids ?? 0) + 1
    bidAmount.value  = ''
    bidSuccess.value = true
    setTimeout(() => { bidSuccess.value = false }, 3000)
  } catch (e) {
    bidError.value = e.response?.data?.error ?? 'Erro ao enviar lance'
  } finally {
    bidLoading.value = false
  }
}

function connectWS() {
  if (!isActive.value) return
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${proto}//${location.host}/ws/auctions/${id}?token=${auth.token}`)
  ws.onopen  = () => { wsConnected.value = true }
  ws.onclose = () => { wsConnected.value = false }
  ws.onerror = () => { wsConnected.value = false }
  ws.onmessage = ({ data }) => {
    try {
      const bid = JSON.parse(data)
      liveFeed.value.unshift(bid)
      if (liveFeed.value.length > 20) liveFeed.value.pop()
      auction.value.current_bid = bid.amount
      auction.value.total_bids  = (auction.value.total_bids ?? 0) + 1
    } catch {}
  }
}

onMounted(async () => {
  await loadData()
  connectWS()
})

onUnmounted(() => ws?.close())
</script>

<style scoped>
.stats-row {
  display: flex;
  gap: 0;
  padding: 0;
  overflow: hidden;
}
.stat-box {
  flex: 1;
  padding: 16px 20px;
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.stat-box:last-child { border-right: none; }
.stat-box .label { font-size: 11px; color: var(--text-light); text-transform: uppercase; letter-spacing: .04em; }
.stat-box .value { font-size: 18px; font-weight: 700; }
.stat-box.highlight .value { color: var(--primary); }
.stat-box.timer .value { color: var(--warning); }

.detail-grid { display: grid; grid-template-columns: 380px 1fr; gap: 20px; }

.feed-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.live-dot { font-size: 12px; color: var(--text-light); }
.live-dot.connected { color: var(--success); }

.feed-list { display: flex; flex-direction: column; gap: 8px; max-height: 360px; overflow-y: auto; }
.feed-item {
  padding: 10px 14px;
  border-radius: 6px;
  background: #f8fafc;
  border: 1px solid var(--border);
  transition: background .3s;
}
.feed-item.new { background: #eff6ff; border-color: #bfdbfe; }
.feed-amount { font-size: 16px; font-weight: 700; color: var(--primary); }
.feed-meta { font-size: 12px; color: var(--text-light); margin-top: 2px; }

.top-bid td { background: #fffbeb; }

@media (max-width: 700px) {
  .stats-row { flex-wrap: wrap; }
  .stat-box { border-right: none; border-bottom: 1px solid var(--border); min-width: 50%; }
  .detail-grid { grid-template-columns: 1fr; }
}
</style>
