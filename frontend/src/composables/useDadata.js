import { ref } from 'vue'

const TOKEN = __DADATA_TOKEN__ || ''

const DADATA_URL = 'https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address'

async function suggest(query, params = {}) {
  if (!query || query.length < 2) return []
  try {
    const res = await fetch(DADATA_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Token ${TOKEN}`,
      },
      body: JSON.stringify({ query, count: 10, ...params }),
    })
    const json = await res.json()
    return json.suggestions ?? []
  } catch {
    return []
  }
}

export function useCitySearch() {
  const suggestions = ref([])
  const loading = ref(false)

  async function search(query) {
    loading.value = true
    suggestions.value = await suggest(query, {
      from_bound: { value: 'city' },
      to_bound: { value: 'city' },
      locations: [{ country: 'Россия' }],
    })
    loading.value = false
  }

  function clear() { suggestions.value = [] }

  return { suggestions, loading, search, clear }
}

export function useDistrictSearch() {
  const suggestions = ref([])
  const loading = ref(false)

  async function search(query, cityName) {
    if (!cityName) return
    loading.value = true
    suggestions.value = await suggest(query, {
      from_bound: { value: 'city_district' },
      to_bound: { value: 'city_district' },
      locations: [{ city: cityName }],
    })
    loading.value = false
  }

  function clear() { suggestions.value = [] }

  return { suggestions, loading, search, clear }
}
