<template>
  <div class="landing">
    <div class="hero">
      <div class="hero-icon">🚗</div>
      <h1 class="hero-title">iCarros</h1>
      <p class="hero-subtitle">
        Olá! Você está no <strong>iCarros</strong> — o melhor site de leilões fictícios do mundo.
      </p>
      <p class="hero-desc">
        Cadastre seus carros, crie leilões em tempo real e faça seus lances ao vivo.
      </p>
      <div class="hero-actions">
        <router-link to="/login"    class="btn btn-primary   btn-lg">Entrar</router-link>
        <router-link to="/register" class="btn btn-outline   btn-lg">Criar conta</router-link>
      </div>
    </div>

    <div class="features">
      <div class="feature-card">
        <span class="feat-icon">🏎️</span>
        <h3>Cadastre seus carros</h3>
        <p>Adicione marca, modelo, ano e valor — e coloque em leilão quando quiser.</p>
      </div>
      <div class="feature-card">
        <span class="feat-icon">⚡</span>
        <h3>Lances em tempo real</h3>
        <p>WebSocket mantém todos os participantes sincronizados ao vivo, sem recarregar.</p>
      </div>
      <div class="feature-card">
        <span class="feat-icon">🏆</span>
        <h3>Encerramento automático</h3>
        <p>Leilões encerram sozinhos no horário definido. O maior lance vence.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.landing {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 60%, #0f172a 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 16px 60px;
}

.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 80px 0 60px;
  max-width: 560px;
}

.hero-icon {
  font-size: 64px;
  margin-bottom: 16px;
  filter: drop-shadow(0 0 24px rgba(59,130,246,.5));
}

.hero-title {
  font-size: 52px;
  font-weight: 800;
  color: #fff;
  letter-spacing: -.02em;
  margin: 0 0 20px;
}

.hero-subtitle {
  font-size: 20px;
  color: #94a3b8;
  margin: 0 0 12px;
  line-height: 1.5;
}
.hero-subtitle strong { color: #e2e8f0; }

.hero-desc {
  font-size: 15px;
  color: #64748b;
  margin: 0 0 40px;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  justify-content: center;
}

.btn-lg {
  padding: 12px 32px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
}

.btn-outline {
  background: transparent;
  color: #e2e8f0;
  border: 1.5px solid #334155;
  transition: border-color .15s, color .15s;
}
.btn-outline:hover {
  border-color: #94a3b8;
  color: #fff;
}

.features {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  max-width: 900px;
  width: 100%;
}

.feature-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 28px 24px;
  text-align: center;
  transition: border-color .2s;
}
.feature-card:hover { border-color: #475569; }

.feat-icon {
  font-size: 36px;
  display: block;
  margin-bottom: 14px;
}

.feature-card h3 {
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0 0 8px;
}

.feature-card p {
  font-size: 13px;
  color: #64748b;
  line-height: 1.6;
  margin: 0;
}

@media (max-width: 640px) {
  .hero-title { font-size: 36px; }
  .features   { grid-template-columns: 1fr; }
}
</style>
