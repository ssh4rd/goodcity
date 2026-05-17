<template>
  <div>
    <h1 style="margin-bottom:1.5rem">Лучшие практики городского благоустройства</h1>

    <PracticeFilters @filter="onFilter" />

    <div v-if="store.loading">Загрузка...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else>
      <div v-if="store.practices.length === 0" class="card" style="text-align:center;color:#6b7280">
        Практик не найдено
      </div>
      <PracticeCard v-for="p in store.practices" :key="p.id" :practice="p" />
      <div class="pagination">
        <button class="btn btn-secondary" :disabled="page <= 1" @click="changePage(page - 1)">← Назад</button>
        <span>Страница {{ page }} · Всего {{ store.total }}</span>
        <button class="btn btn-secondary" :disabled="page * perPage >= store.total" @click="changePage(page + 1)">Вперёд →</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { usePracticesStore } from '../stores/practices'
import PracticeCard from '../components/PracticeCard.vue'
import PracticeFilters from '../components/PracticeFilters.vue'

const store = usePracticesStore()
const page = ref(1)
const perPage = 10
const filters = ref({})

async function load() {
  await store.fetchPractices({ ...filters.value, page: page.value, per_page: perPage })
}

function onFilter(f) {
  filters.value = f
  page.value = 1
  load()
}

function changePage(p) {
  page.value = p
  load()
}

onMounted(load)
</script>

<style scoped>
.pagination { display: flex; align-items: center; gap: 1rem; justify-content: center; margin-top: 1.5rem; }
</style>
