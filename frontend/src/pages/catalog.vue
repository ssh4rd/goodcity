<template>
  <div class="catalog-page">
    <div class="catalog-toolbar container">
      <div class="search-wrap">
        <span class="search-icon">
          <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
          </svg>
        </span>
        <input
          v-model="search"
          class="search-input"
          placeholder="Поиск практик..."
          @input="debouncedLoad"
        />
      </div>

      <div class="category-pills">
        <button
          v-for="cat in allCategories"
          :key="cat.value"
          class="pill"
          :class="{ active: activeCategory === cat.value }"
          @click="setCategory(cat.value)"
        >
          {{ cat.label }}
        </button>
      </div>

      <div class="sort-wrap">
        <select v-model="sortBy" class="form-input form-select sort-select" @change="load">
          <option value="">Сортировка: По рейтингу</option>
          <option value="date">По дате</option>
        </select>
      </div>
    </div>

    <div class="container">
      <p v-if="!loading" class="total-count">{{ store.total }} практик найдено</p>

      <div v-if="loading" class="practices-grid">
        <div class="skeleton-card" v-for="i in 6" :key="i"></div>
      </div>
      <div v-else-if="store.practices.length === 0" class="empty">
        <p>Практики не найдены. Попробуйте изменить фильтры.</p>
      </div>
      <div v-else class="practices-grid">
        <PracticeCard v-for="p in store.practices" :key="p.id" :practice="p" />
      </div>

      <div v-if="store.total > perPage" class="pagination">
        <button class="btn btn-outline btn-sm" :disabled="page <= 1" @click="changePage(page - 1)">← Назад</button>
        <span class="page-info">Страница {{ page }} из {{ totalPages }}</span>
        <button class="btn btn-outline btn-sm" :disabled="page >= totalPages" @click="changePage(page + 1)">Вперёд →</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePracticesStore } from '../stores/practices'
import PracticeCard from '../components/PracticeCard.vue'

const store = usePracticesStore()
const search = ref('')
const activeCategory = ref('')
const sortBy = ref('')
const page = ref(1)
const perPage = 9

const allCategories = [
  { value: '', label: 'Все категории' },
  { value: 'Транспорт', label: 'Транспорт' },
  { value: 'Экология', label: 'Экология' },
  { value: 'Инфраструктура', label: 'Инфраструктура' },
  { value: 'Благоустройство', label: 'Общественные пространства' },
  { value: 'ЖКХ', label: 'ЖКХ' },
]

const totalPages = computed(() => Math.ceil(store.total / perPage))

let debounceTimer = null
function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; load() }, 400)
}

function setCategory(val) {
  activeCategory.value = val
  page.value = 1
  load()
}

function changePage(p) {
  page.value = p
  load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function load() {
  await store.fetchPractices({
    city: '',
    category: activeCategory.value,
    search: search.value,
    page: page.value,
    per_page: perPage,
  })
}

onMounted(load)
</script>

<style scoped>
.catalog-page { padding-bottom: 64px; }

.catalog-toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 20px;
  padding-bottom: 20px;
  flex-wrap: wrap;
}

.search-wrap {
  position: relative;
  min-width: 240px;
}
.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--c-muted);
  display: flex;
}
.search-input {
  width: 100%;
  padding: 10px 14px 10px 38px;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  outline: none;
  background: var(--c-white);
  transition: border-color 0.15s;
}
.search-input:focus { border-color: var(--c-green); }

.category-pills {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  flex: 1;
}
.pill {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  background: var(--c-white);
  color: var(--c-muted);
  border: 1.5px solid var(--c-border);
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
}
.pill:hover { border-color: var(--c-green); color: var(--c-green); }
.pill.active {
  background: var(--c-green);
  color: white;
  border-color: var(--c-green);
}

.sort-select { width: 200px; font-size: 14px; padding: 9px 36px 9px 14px; }

.total-count {
  font-size: 14px;
  color: var(--c-muted);
  margin-bottom: 20px;
}

.practices-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.empty {
  text-align: center;
  padding: 64px 0;
  color: var(--c-muted);
}

.skeleton-card {
  height: 320px;
  background: linear-gradient(90deg, #F3F4F6 25%, #E5E7EB 50%, #F3F4F6 75%);
  background-size: 200% 100%;
  border-radius: var(--radius);
  animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 40px;
}
.page-info { font-size: 14px; color: var(--c-muted); }

@media (max-width: 768px) {
  .catalog-toolbar { flex-wrap: wrap; gap: 10px; padding-top: 14px; padding-bottom: 14px; }
  .search-wrap { min-width: 100%; order: -1; }
  .category-pills {
    flex-wrap: nowrap;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    padding-bottom: 4px;
    scrollbar-width: none;
  }
  .category-pills::-webkit-scrollbar { display: none; }
  .pill { flex-shrink: 0; }
  .sort-select { width: 100%; }
  .practices-grid { grid-template-columns: 1fr; }
  .skeleton-card { height: 280px; }
}

@media (max-width: 480px) {
  .practices-grid { grid-template-columns: 1fr; gap: 12px; }
}
</style>
