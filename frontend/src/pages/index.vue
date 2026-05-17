<template>
  <div class="home">
    <!-- Hero -->
    <section class="hero">
      <div class="hero-inner container">
        <div class="hero-text">
          <h1 class="hero-title">Лучшие практики городского благоустройства</h1>
          <p class="hero-subtitle">
            Краудсорс-платформа для сбора, оценки и распространения успешных решений в области городской среды
          </p>
          <div class="hero-actions">
            <router-link to="/catalog" class="btn-hero-primary">Смотреть практики</router-link>
            <router-link to="/submit" class="btn-hero-outline">Добавить свою</router-link>
          </div>
        </div>
        <div class="hero-visual">
          <div class="hero-circles">
            <div class="circle circle-1"></div>
            <div class="circle circle-2"></div>
          </div>
          <div class="hero-card">
            <div class="hero-card-top">
              <div class="hero-card-avatar"></div>
              <div>
                <div class="hero-card-score">4.9 ★</div>
                <div class="hero-card-name">Умные парковки</div>
              </div>
            </div>
            <div class="hero-card-divider"></div>
            <p class="hero-card-stats">1 240 практик · 342 города · 8 500+ экспертов</p>
            <div class="hero-card-bars">
              <div class="hbar-track"><div class="hbar-fill" style="width:80%"></div></div>
              <div class="hbar-track"><div class="hbar-fill" style="width:65%"></div></div>
              <div class="hbar-track"><div class="hbar-fill" style="width:72%"></div></div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats -->
    <section class="stats">
      <div class="container">
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-accent" style="background:var(--c-green)"></div>
            <div class="stat-value">1 240</div>
            <div class="stat-label">практик в базе</div>
          </div>
          <div class="stat-card">
            <div class="stat-accent" style="background:var(--c-blue)"></div>
            <div class="stat-value">342</div>
            <div class="stat-label">города России</div>
          </div>
          <div class="stat-card">
            <div class="stat-accent" style="background:#8B5CF6"></div>
            <div class="stat-value">8 500+</div>
            <div class="stat-label">экспертов</div>
          </div>
          <div class="stat-card">
            <div class="stat-accent" style="background:var(--c-amber)"></div>
            <div class="stat-value">4.6</div>
            <div class="stat-label">средний рейтинг</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Best practices -->
    <section class="best-section">
      <div class="container">
        <div class="section-header">
          <div>
            <h2 class="section-title">Лучшие практики</h2>
            <p class="section-sub">Топ по рейтингу за последний месяц</p>
          </div>
          <router-link to="/catalog" class="btn btn-outline btn-sm">Смотреть все →</router-link>
        </div>

        <div v-if="loading" class="loading-grid">
          <div class="skeleton-card" v-for="i in 3" :key="i"></div>
        </div>
        <div v-else class="practices-grid">
          <PracticeCard v-for="p in topPractices" :key="p.id" :practice="p" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { api } from '../api'
import PracticeCard from '../components/PracticeCard.vue'

const practices = ref([])
const loading = ref(true)

const topPractices = computed(() => practices.value.slice(0, 3))

onMounted(async () => {
  try {
    const data = await api.getPractices({ per_page: 3 })
    practices.value = data.data || []
  } catch {}
  finally { loading.value = false }
})
</script>

<style scoped>
/* Hero */
.hero {
  background: var(--c-navy);
  padding: 72px 0 80px;
}
.hero-inner {
  display: flex;
  align-items: center;
  gap: 64px;
}
.hero-text { flex: 1; max-width: 620px; }
.hero-title {
  font-size: 48px;
  font-weight: 700;
  color: white;
  line-height: 1.15;
  margin-bottom: 20px;
}
.hero-subtitle {
  font-size: 18px;
  color: rgba(255,255,255,0.7);
  line-height: 1.6;
  margin-bottom: 36px;
  max-width: 520px;
}
.hero-actions { display: flex; gap: 12px; flex-wrap: wrap; }
.btn-hero-primary {
  background: var(--c-green);
  color: white;
  font-size: 16px;
  font-weight: 600;
  padding: 14px 28px;
  border-radius: 10px;
  transition: background 0.15s;
}
.btn-hero-primary:hover { background: var(--c-green-dark); }
.btn-hero-outline {
  background: transparent;
  color: white;
  font-size: 16px;
  font-weight: 500;
  padding: 13px 28px;
  border-radius: 10px;
  border: 1.5px solid rgba(255,255,255,0.3);
  transition: border-color 0.15s;
}
.btn-hero-outline:hover { border-color: rgba(255,255,255,0.6); }

/* Hero visual */
.hero-visual { position: relative; width: 400px; height: 280px; flex-shrink: 0; }
.hero-circles { position: absolute; inset: 0; }
.circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.25;
}
.circle-1 {
  width: 240px; height: 240px;
  background: var(--c-green);
  top: 0; left: 20px;
}
.circle-2 {
  width: 260px; height: 200px;
  background: #4B9EFF;
  bottom: 0; right: 0;
}
.hero-card {
  position: relative;
  z-index: 1;
  background: rgba(255,255,255,0.12);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 16px;
  padding: 20px;
  margin: 20px;
}
.hero-card-top { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.hero-card-avatar {
  width: 44px; height: 44px;
  background: var(--c-green);
  border-radius: 50%;
  flex-shrink: 0;
}
.hero-card-score { font-size: 18px; font-weight: 700; color: white; }
.hero-card-name { font-size: 13px; color: rgba(255,255,255,0.7); }
.hero-card-divider { height: 1px; background: rgba(255,255,255,0.15); margin-bottom: 10px; }
.hero-card-stats { font-size: 12px; color: rgba(255,255,255,0.6); margin-bottom: 14px; }
.hero-card-bars { display: flex; flex-direction: column; gap: 6px; }
.hbar-track { height: 6px; background: rgba(255,255,255,0.15); border-radius: 3px; overflow: hidden; }
.hbar-fill { height: 100%; background: var(--c-green); border-radius: 3px; }

/* Stats */
.stats { padding: 32px 0; }
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.stat-card {
  background: var(--c-white);
  border: 1px solid var(--c-border);
  border-radius: var(--radius);
  padding: 20px 20px 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  position: relative;
  overflow: hidden;
}
.stat-accent {
  position: absolute;
  left: 0;
  top: 16px;
  bottom: 16px;
  width: 4px;
  border-radius: 0 2px 2px 0;
}
.stat-value { font-size: 32px; font-weight: 700; color: var(--c-text); line-height: 1.2; }
.stat-label { font-size: 14px; color: var(--c-muted); }

/* Best section */
.best-section { padding: 40px 0 64px; }
.section-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 28px;
}
.section-title { font-size: 28px; font-weight: 700; margin-bottom: 4px; }
.section-sub { font-size: 14px; color: var(--c-muted); }

.practices-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}
.loading-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}
.skeleton-card {
  height: 320px;
  background: linear-gradient(90deg, #F3F4F6 25%, #E5E7EB 50%, #F3F4F6 75%);
  background-size: 200% 100%;
  border-radius: var(--radius);
  animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

@media (max-width: 768px) {
  .hero { padding: 40px 0 48px; }
  .hero-inner { flex-direction: column; gap: 32px; }
  .hero-title { font-size: 30px; }
  .hero-subtitle { font-size: 15px; margin-bottom: 24px; }
  .hero-visual { display: none; }
  .btn-hero-primary, .btn-hero-outline { font-size: 14px; padding: 12px 20px; }

  .stats-grid { grid-template-columns: 1fr 1fr; gap: 10px; }
  .stat-value { font-size: 24px; }

  .section-header { flex-direction: column; align-items: flex-start; gap: 12px; }
  .practices-grid { grid-template-columns: 1fr; }
  .loading-grid { grid-template-columns: 1fr; }
}

@media (max-width: 480px) {
  .hero-title { font-size: 24px; }
  .hero-actions { flex-direction: column; }
  .btn-hero-primary, .btn-hero-outline { width: 100%; text-align: center; }
}
</style>
