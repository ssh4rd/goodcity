<template>
  <div class="practice-page">
    <div v-if="loading" class="container" style="padding-top:48px">
      <div class="skeleton-block" style="height:32px;width:300px;margin-bottom:32px"></div>
      <div class="skeleton-block" style="height:400px"></div>
    </div>

    <div v-else-if="error" class="container" style="padding-top:48px">
      <p class="error-state">{{ error }}</p>
    </div>

    <template v-else-if="practice">
      <div class="container practice-container">
        <!-- Breadcrumbs -->
        <nav class="breadcrumbs">
          <router-link to="/catalog">Каталог</router-link>
          <span class="sep">/</span>
          <span>{{ practice.category }}</span>
          <span class="sep">/</span>
          <span class="current">{{ practice.title }}</span>
        </nav>

        <h1 class="practice-title">{{ practice.title }}</h1>

        <!-- Tags -->
        <div class="tags">
          <span class="tag">{{ practice.category }}</span>
          <span class="tag">{{ practice.city }}</span>
          <span :class="`badge badge-${practice.status}`">{{ statusLabel }}</span>
        </div>

        <!-- Main content -->
        <div class="practice-body">
          <!-- Left: image -->
          <div class="practice-image-wrap">
            <div class="practice-image">
              <svg width="64" height="52" viewBox="0 0 64 52" fill="none">
                <rect x="2" y="2" width="60" height="48" rx="4" fill="#D1D5DB"/>
                <rect x="10" y="10" width="44" height="32" rx="2" fill="#E5E7EB"/>
                <text x="32" y="32" text-anchor="middle" fill="#9CA3AF" font-size="13" font-family="Inter,sans-serif">Фото практики</text>
              </svg>
            </div>
          </div>

          <!-- Right: info -->
          <div class="practice-info">
            <div class="info-card card">
              <!-- Rating header -->
              <div class="rating-header">
                <div class="rating-big">
                  <span class="rating-score">{{ displayScore }}</span>
                  <span class="rating-stars stars">★★★★★</span>
                  <span class="rating-votes">{{ practice.vote_count }} оценок</span>
                </div>
                <div class="rating-meta">
                  <div class="meta-item">
                    <span class="meta-label">{{ practice.vote_count }} комментариев</span>
                    <span class="meta-date">{{ formattedDate }}</span>
                  </div>
                  <div class="meta-item">
                    <span class="meta-label">{{ practice.city }}</span>
                    <span class="meta-sub">Россия</span>
                  </div>
                </div>
              </div>

              <!-- Criteria -->
              <div class="criteria-section">
                <h3 class="criteria-title">Критерии оценки</h3>
                <div class="criteria-bar-lg">
                  <div class="criteria-bar-lg-row">
                    <span class="criteria-bar-lg-label">Польза</span>
                    <div class="criteria-bar-track-lg"><div class="criteria-bar-fill blue" :style="{ width: utilityPct + '%' }"></div></div>
                    <span class="criteria-pct">{{ utilityPct }}%</span>
                  </div>
                  <div class="criteria-bar-lg-row">
                    <span class="criteria-bar-lg-label">Эстетика</span>
                    <div class="criteria-bar-track-lg"><div class="criteria-bar-fill green" :style="{ width: aestheticPct + '%' }"></div></div>
                    <span class="criteria-pct">{{ aestheticPct }}%</span>
                  </div>
                  <div class="criteria-bar-lg-row">
                    <span class="criteria-bar-lg-label">Экология</span>
                    <div class="criteria-bar-track-lg"><div class="criteria-bar-fill teal" :style="{ width: ecologyPct + '%' }"></div></div>
                    <span class="criteria-pct">{{ ecologyPct }}%</span>
                  </div>
                </div>
              </div>

              <!-- Description -->
              <div class="desc-section">
                <h3 class="criteria-title">Описание</h3>
                <p class="desc-text">{{ practice.description }}</p>
              </div>

              <!-- Author -->
              <div class="author-card">
                <div class="author-avatar">{{ authorInitials }}</div>
                <div class="author-info">
                  <div class="author-name">Автор #{{ practice.author_id }}</div>
                  <div class="author-meta">Эксперт · Репутация 4.2 · {{ practice.vote_count }} практик</div>
                </div>
                <button class="btn btn-outline btn-sm">Подписаться</button>
              </div>

              <!-- Actions -->
              <div class="practice-actions">
                <VoteButton :practice-id="practice.id" :initial-count="practice.vote_count" />
                <button class="btn btn-outline">↗ Поделиться</button>
              </div>
            </div>
          </div>
        </div>

        <!-- Comments -->
        <CommentList :practice-id="practice.id" />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../../api'
import VoteButton from '../../components/VoteButton.vue'
import CommentList from '../../components/CommentList.vue'

const route = useRoute()
const practice = ref(null)
const loading = ref(true)
const error = ref(null)

const statusLabels = { pending: 'На модерации', approved: 'Одобрено', rejected: 'Отклонено' }
const statusLabel = computed(() => statusLabels[practice.value?.status] || '')
const formattedDate = computed(() =>
  practice.value ? new Date(practice.value.created_at).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }) : ''
)
const displayScore = computed(() => {
  if (!practice.value) return '0.0'
  return (3.5 + Math.min(1.4, (practice.value.vote_count || 0) * 0.02)).toFixed(1)
})
const authorInitials = computed(() => `#${practice.value?.author_id || '?'}`)

const utilityPct = computed(() => practice.value ? 50 + ((practice.value.id * 17 + 3) % 45) : 0)
const aestheticPct = computed(() => practice.value ? 40 + ((practice.value.id * 23 + 7) % 50) : 0)
const ecologyPct = computed(() => practice.value ? 35 + ((practice.value.id * 31 + 11) % 55) : 0)

onMounted(async () => {
  try {
    practice.value = await api.getPractice(route.params.id)
  } catch {
    error.value = 'Практика не найдена'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.practice-page { padding-bottom: 64px; }
.practice-container { padding-top: 24px; }

/* Breadcrumbs */
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--c-muted);
  margin-bottom: 16px;
}
.breadcrumbs a { color: var(--c-muted); }
.breadcrumbs a:hover { color: var(--c-text); }
.sep { color: var(--c-border); }
.current { color: var(--c-text); }

.practice-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 16px;
  line-height: 1.2;
}

.tags { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 28px; }
.tag {
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  background: #F3F4F6;
  color: var(--c-muted);
  border: 1px solid var(--c-border);
}

/* Body layout */
.practice-body {
  display: grid;
  grid-template-columns: 1fr 1.1fr;
  gap: 24px;
  margin-bottom: 40px;
}

.practice-image-wrap {}
.practice-image {
  background: #F3F4F6;
  border-radius: var(--radius);
  border: 1px solid var(--c-border);
  aspect-ratio: 4/3;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Info card */
.info-card { padding: 24px; display: flex; flex-direction: column; gap: 20px; }

.rating-header {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--c-border);
}
.rating-big { display: flex; flex-direction: column; gap: 4px; }
.rating-score { font-size: 40px; font-weight: 700; line-height: 1; color: var(--c-text); }
.rating-stars { font-size: 18px; }
.rating-votes { font-size: 13px; color: var(--c-muted); }

.rating-meta { display: flex; flex-direction: column; gap: 12px; padding-top: 4px; }
.meta-item { display: flex; flex-direction: column; gap: 2px; }
.meta-label { font-size: 14px; font-weight: 500; }
.meta-date, .meta-sub { font-size: 13px; color: var(--c-muted); }

.criteria-title { font-size: 15px; font-weight: 600; margin-bottom: 12px; }
.criteria-bar-lg { display: flex; flex-direction: column; gap: 8px; }
.criteria-bar-lg-row {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: var(--c-muted);
}
.criteria-bar-lg-label { width: 64px; flex-shrink: 0; }
.criteria-bar-track-lg {
  flex: 1;
  height: 8px;
  background: #F3F4F6;
  border-radius: 4px;
  overflow: hidden;
}
.criteria-pct { width: 36px; text-align: right; font-size: 12px; }

.desc-section {}
.desc-text {
  font-size: 14px;
  line-height: 1.7;
  color: var(--c-text);
}

.author-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #F8FAFB;
  border-radius: var(--radius-sm);
  border: 1px solid var(--c-border);
}
.author-avatar {
  width: 42px; height: 42px;
  border-radius: 50%;
  background: var(--c-green);
  color: white;
  font-size: 11px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.author-name { font-size: 14px; font-weight: 600; }
.author-meta { font-size: 12px; color: var(--c-muted); }

.practice-actions { display: flex; gap: 10px; }

/* Skeleton */
.skeleton-block {
  background: linear-gradient(90deg, #F3F4F6 25%, #E5E7EB 50%, #F3F4F6 75%);
  background-size: 200% 100%;
  border-radius: var(--radius);
  animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

.error-state { color: var(--c-red-dark); font-size: 16px; }

@media (max-width: 768px) {
  .practice-title { font-size: 22px; }
  .practice-body { grid-template-columns: 1fr; }
  .practice-image { aspect-ratio: 16/9; }
  .rating-header { grid-template-columns: 1fr; gap: 12px; }
  .rating-big { flex-direction: row; align-items: center; gap: 12px; }
  .rating-score { font-size: 32px; }
  .rating-meta { flex-direction: row; flex-wrap: wrap; gap: 16px; }
  .author-card { flex-wrap: wrap; }
  .practice-actions { flex-direction: column; }
  .practice-actions > * { width: 100%; justify-content: center; }
}

@media (max-width: 480px) {
  .tags { flex-wrap: wrap; }
  .info-card { padding: 16px; }
}
</style>
