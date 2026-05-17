<template>
  <div class="filters card">
    <div class="filters-row">
      <div class="form-group">
        <label>Город</label>
        <input v-model="localCity" placeholder="Все города" @change="applyFilter" />
      </div>
      <div class="form-group">
        <label>Категория</label>
        <select v-model="localCategory" @change="applyFilter">
          <option value="">Все категории</option>
          <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
        </select>
      </div>
      <button class="btn btn-secondary" @click="reset">Сбросить</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['filter'])

const categories = ['Транспорт', 'Благоустройство', 'ЖКХ', 'Экология', 'Безопасность', 'Культура', 'Образование', 'Здоровье']

const localCity = ref('')
const localCategory = ref('')

function applyFilter() {
  emit('filter', { city: localCity.value, category: localCategory.value })
}

function reset() {
  localCity.value = ''
  localCategory.value = ''
  emit('filter', {})
}
</script>

<style scoped>
.filters { margin-bottom: 1.5rem; }
.filters-row { display: flex; gap: 1rem; align-items: flex-end; flex-wrap: wrap; }
.filters-row .form-group { flex: 1; min-width: 150px; margin-bottom: 0; }
</style>
