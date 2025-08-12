<template>
  <div>
    <h2 class="mx-5 my-5">GoNews - агрегатор новостей.</h2>

    <!-- Панель управления -->
    <v-card class="mx-5 my-5 pa-4" elevation="5">
      <v-card-title>Управление новостями</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="4">
            <v-btn
                @click="getAllNews"
                color="primary"
                class="ma-2"
                :loading="loading"
            >
              Получить все новости
            </v-btn>

            <v-btn
                @click="getLastNews"
                color="secondary"
                class="ma-2"
                :loading="loading"
            >
              Получить последнюю новость
            </v-btn>
          </v-col>

          <v-col cols="12" md="4">
            <v-text-field
                v-model="search"
                label="Поиск по новостям"
                dense
                outlined
                class="ma-2"
                clearable
            ></v-text-field>
            <v-btn
                @click="searchNews"
                color="info"
                class="ma-2"
                :loading="loading"
                :disabled="!search"
            >
              Найти новости
            </v-btn>
          </v-col>

          <v-col cols="12" md="4">
            <v-text-field
                v-model="limitCount"
                label="Новостей на странице"
                type="number"
                min="1"
                max="100"
                dense
                outlined
                class="ma-2"
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Сообщения об ошибках -->
    <v-alert
        v-if="error"
        type="error"
        dismissible
        @input="error = ''"
        class="mx-5 my-3"
    >
      {{ error }}
    </v-alert>

    <!-- Информационное сообщение -->
    <v-alert
        v-if="info"
        type="info"
        dismissible
        @input="info = ''"
        class="mx-5 my-3"
    >
      {{ info }}
    </v-alert>

    <!-- Отображение новостей -->
    <div v-if="loading" class="text-center ma-5">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
      <p>Загрузка новостей...</p>
    </div>

    <div v-else-if="news.length === 0 && !loading" class="text-center ma-5">
      <p>Новости не загружены. Выберите один из вариантов выше.</p>
    </div>

    <div v-else>
      <v-card class="mx-5 my-3" elevation="2">
        <v-card-text class="d-flex justify-space-between align-center">
          <h3>Найдено новостей: {{ totalNews }}</h3>
          <div v-if="totalPages > 1">
            Страница {{ currentPage }} из {{ totalPages }}
          </div>
        </v-card-text>
      </v-card>

      <!-- Пагинация сверху -->
      <div v-if="totalPages > 1" class="text-center ma-3">
        <v-pagination
            v-model="currentPage"
            :length="totalPages"
            :total-visible="7"
            color="primary"
            @update:modelValue="loadPage"
        />
      </div>

      <!-- Список новостей -->
      <div v-for="post in news" :key="post.id">
        <v-card elevation="10" outlined class="mx-5 my-5">
          <v-card-title>
            {{ post.title }}
          </v-card-title>
          <v-card-text>
            {{ post.content }}
            <v-card-subtitle>
              ID: {{ post.id }} | {{ post.pub_time }}
            </v-card-subtitle>

            <div class="mt-3">
              <v-btn
                  color="primary"
                  class="mr-2"
                  @click="goToNewsDetail(post.id)"
              >
                Перейти
              </v-btn>
              <v-btn
                  color="secondary"
                  :href="post.link"
                  target="_blank"
                  v-if="post.link"
              >
                Перейти к источнику
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </div>

      <!-- Пагинация снизу -->
      <div v-if="totalPages > 1" class="text-center ma-3">
        <v-pagination
            v-model="currentPage"
            :length="totalPages"
            :total-visible="7"
            color="primary"
            @update:modelValue="loadPage"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: "News",
  setup() {
    const router = useRouter()

    // Реактивные переменные
    const news = ref([])
    const loading = ref(false)
    const error = ref('')
    const info = ref('')
    const limitCount = ref(10)
    const search = ref('')
    const currentPage = ref(1)
    const totalPages = ref(1)
    const totalNews = ref(0)

    // Получение baseUrl из переменных окружения
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

    // Методы
    const fetchNews = async (endpoint, params = {}, successMessage = '') => {
      loading.value = true
      error.value = ''
      info.value = ''

      try {
        const url = new URL(baseUrl + endpoint)

        // Добавляем параметры запроса
        Object.keys(params).forEach(key => {
          if (params[key] !== null && params[key] !== undefined && params[key] !== '') {
            url.searchParams.append(key, params[key])
          }
        })

        const response = await fetch(url)

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }

        const data = await response.json()

        // API возвращает данные в формате {posts: [...], total: number}
        if (data && data.data.news && Array.isArray(data.data.news)) {
          news.value = data.data.news
          totalNews.value = data.data.total || data.data.news.length

          // Вычисляем общее количество страниц
          totalPages.value = Math.ceil(totalNews.value / limitCount.value)
        } else if (Array.isArray(data.data)) {
          news.value = data
          totalNews.value = data.length
          totalPages.value = 1
        } else if (data.data && data.data.post && typeof data.data.post === 'object') {
          // Если получили один объект, оборачиваем в массив
          news.value = [data.data.post]
          totalNews.value = 1
          totalPages.value = 1
        } else {
          news.value = []
          totalNews.value = 0
          totalPages.value = 1
        }

        if (successMessage) {
          info.value = successMessage
        }

      } catch (err) {
        error.value = `Ошибка при загрузке новостей: ${err.message}`
        news.value = []
        totalNews.value = 0
        totalPages.value = 1
      } finally {
        loading.value = false
      }
    }

    const getAllNews = () => {
      const params = {
        page: 1,
        limit: limitCount.value
      }
      currentPage.value = 1
      fetchNews('/news', params, 'Загружены все новости')
    }

    const getLastNews = () => {
      fetchNews('/news/last', {}, 'Загружена последняя новость')
    }

    const searchNews = () => {
      if (!search.value) {
        error.value = 'Введите поисковый запрос'
        return
      }
      const params = {
        page: 1,
        limit: limitCount.value,
        search: search.value
      }
      currentPage.value = 1
      fetchNews('/news', params, `Найдены новости по запросу: ${search.value}`)
    }

    const loadPage = (page) => {
      currentPage.value = page
      const params = {
        page: page,
        limit: limitCount.value
      }

      if (search.value) {
        params.search = search.value
      }

      fetchNews('/news', params)
    }

    const goToNewsDetail = (newsId) => {
      router.push(`/news/${newsId}`)
    }

    watch(limitCount, (newVal, oldVal) => {
      if (newVal !== oldVal && newVal > 0) {
        currentPage.value = 1
        loadPage(1)
      }
    })

    return {
      // Реактивные переменные
      news,
      loading,
      error,
      info,
      limitCount,
      search,
      currentPage,
      totalPages,
      totalNews,

      // Методы
      getAllNews,
      getLastNews,
      searchNews,
      loadPage,
      goToNewsDetail
    }
  }
}
</script>

<style scoped>
.v-card-title a {
  text-decoration: none;
  color: inherit;
}

.v-card-title a:hover {
  text-decoration: underline;
}
</style>