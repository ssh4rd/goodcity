<template>
  <router-link :to="`/practices/${practice.id}`" class="practice-card">
    <div class="card-image">
      <span class="status-dot" :class="practice.status"></span>
      <div class="image-placeholder">
        <svg width="40" height="32" viewBox="0 0 40 32" fill="none">
          <rect x="2" y="2" width="36" height="28" rx="2" fill="#D1D5DB"/>
          <rect x="8" y="8" width="24" height="16" rx="1" fill="#9CA3AF"/>
          <rect x="10" y="10" width="20" height="12" rx="1" fill="#E5E7EB"/>
        </svg>
      </div>
    </div>
    <div class="card-body">
      <h3 class="card-title">{{ practice.title }}</h3>
      <div class="card-rating">
        <span class="stars">{{ starStr }}</span>
        <span class="score">{{ displayScore }}</span>
      </div>
      <div class="card-bars">
        <div class="criteria-bar">
          <span class="criteria-bar-label">Польза</span>
          <div class="criteria-bar-track"><div class="criteria-bar-fill blue" :style="{ width: utilityPct + '%' }"></div></div>
        </div>
        <div class="criteria-bar">
          <span class="criteria-bar-label">Эстетика</span>
          <div class="criteria-bar-track"><div class="criteria-bar-fill green" :style="{ width: aestheticPct + '%' }"></div></div>
        </div>
        <div class="criteria-bar">
          <span class="criteria-bar-label">Экология</span>
          <div class="criteria-bar-track"><div class="criteria-bar-fill teal" :style="{ width: ecologyPct + '%' }"></div></div>
        </div>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  practice: { type: Object, required: true }
})

const displayScore = computed(() => {
  const base = 3.5 + Math.min(1.4, (props.practice.vote_count || 0) * 0.02)
  return base.toFixed(1)
})

const starStr = computed(() => {
  const s = Math.round(parseFloat(displayScore.value))
  return '★'.repeat(s) + '☆'.repeat(5 - s)
})

// Deterministic percentages based on id
const utilityPct = computed(() => 50 + ((props.practice.id * 17 + 3) % 45))
const aestheticPct = computed(() => 40 + ((props.practice.id * 23 + 7) % 50))
const ecologyPct = computed(() => 35 + ((props.practice.id * 31 + 11) % 55))
</script>

<style scoped>
.practice-card {
  display: flex;
  flex-direction: column;
  background: var(--c-white);
  border-radius: var(--radius);
  border: 1px solid var(--c-border);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  transition: box-shadow 0.2s, transform 0.2s;
  text-decoration: none;
  color: inherit;
}
.practice-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.card-image {
  position: relative;
  background: #F3F4F6;
  height: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.image-placeholder { opacity: 0.6; }

.status-dot {
  position: absolute;
  top: 12px;
  left: 12px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--c-green);
}
.status-dot.pending { background: var(--c-amber); }
.status-dot.rejected { background: var(--c-red); }

.card-body { padding: 16px; flex: 1; display: flex; flex-direction: column; gap: 10px; }

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--c-text);
  line-height: 1.4;
}

.card-rating {
  display: flex;
  align-items: center;
  gap: 2px;
}
.stars { color: var(--c-amber); font-size: 14px; }
.score { font-size: 13px; font-weight: 600; color: var(--c-text); margin-left: 6px; }

.card-bars { display: flex; flex-direction: column; gap: 5px; }

@media (max-width: 768px) {
  .card-image { height: 140px; }
  .card-title { font-size: 14px; }
}
</style>
