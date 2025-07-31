<template>
  <div>
    <h2 class="mx-5 my-5">GoNews - агрегатор новостей.</h2>

    <!-- Панель управления -->
    <v-card class="mx-5 my-5 pa-4" elevation="5">
      <v-card-title>Управление новостями</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
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

          <v-col cols="12" md="6">
            <v-text-field
                v-model="limitCount"
                label="Количество новостей"
                type="number"
                min="1"
                dense
                outlined
                class="ma-2"
            ></v-text-field>
            <v-btn
                @click="getLatestNews"
                color="success"
                class="ma-2"
                :loading="loading"
                :disabled="!limitCount || limitCount < 1"
            >
              Получить последние {{ limitCount || 'N' }} новостей
            </v-btn>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
                v-model="newsId"
                label="ID новости"
                type="number"
                min="1"
                dense
                outlined
                class="ma-2"
            ></v-text-field>
            <v-btn
                @click="getNewsById"
                color="warning"
                class="ma-2"
                :loading="loading"
                :disabled="!newsId || newsId < 1"
            >
              Получить новость по ID {{ newsId || 'X' }}
            </v-btn>
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
          <h3>Найдено новостей: {{ news.length }}</h3>
          <div v-if="news.length > itemsPerPage">
            Страница {{ currentPage }} из {{ totalPages }}
            (показано {{ paginatedNews.length }} из {{ news.length }})
          </div>
        </v-card-text>
      </v-card>

      <!-- Пагинация сверху (если много новостей) -->
      <div v-if="totalPages > 1" class="text-center ma-3">
        <v-pagination
            v-model="currentPage"
            :length="totalPages"
            :total-visible="7"
            color="primary"
        ></v-pagination>
      </div>

      <!-- Список новостей -->
      <div v-for="post in paginatedNews" :key="post.id">
        <v-card elevation="10" outlined class="mx-5 my-5">
          <v-card-title>
            <a :href="post.link" target="_blank"> {{ post.title }} </a>
          </v-card-title>
          <v-card-text>
            {{ post.content }}
            <v-card-subtitle>
              ID: {{ post.id }} | {{ post.pub_time }}
            </v-card-subtitle>
          </v-card-text>
        </v-card>
      </div>

      <!-- Пагинация снизу (если много новостей) -->
      <div v-if="totalPages > 1" class="text-center ma-3">
        <v-pagination
            v-model="currentPage"
            :length="totalPages"
            :total-visible="7"
            color="primary"
        ></v-pagination>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'

export default {
  name: "News",
  setup() {
    // Реактивные переменные
    const news = ref([])
    const loading = ref(false)
    const error = ref('')
    const info = ref('')
    const limitCount = ref(null)
    const newsId = ref(null)
    const currentPage = ref(1)
    const itemsPerPage = ref(10)

    // Получение baseUrl из переменных окружения
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

    // Вычисляемые свойства
    const totalPages = computed(() => {
      return Math.ceil(news.value.length / itemsPerPage.value)
    })

    const paginatedNews = computed(() => {
      const start = (currentPage.value - 1) * itemsPerPage.value
      const end = start + itemsPerPage.value
      return news.value.slice(start, end)
    })

    // Методы
    const fetchNews = async (endpoint, successMessage = '', expectWrapped = false) => {
      loading.value = true
      error.value = ''
      info.value = ''

      try {
        const url = baseUrl + endpoint
        const response = await fetch(url)

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }

        const data = await response.json()

        if (expectWrapped) {
          // Для /news и /news/latest?limit= - данные в формате {posts: [...]}
          if (data && data.posts && Array.isArray(data.posts)) {
            news.value = data.posts
          } else {
            news.value = []
          }
        } else {
          // Для /news/last и /news/id - возвращается напрямую объект новости
          if (Array.isArray(data)) {
            news.value = data
          } else if (data && typeof data === 'object') {
            // Если получили один объект, оборачиваем в массив
            news.value = [data]
          } else {
            news.value = []
          }
        }

        if (successMessage) {
          info.value = successMessage
        }

        // Сбрасываем на первую страницу при новой загрузке
        currentPage.value = 1

      } catch (err) {
        error.value = `Ошибка при загрузке новостей: ${err.message}`
        news.value = []
      } finally {
        loading.value = false
      }
    }

    const getAllNews = () => {
      fetchNews('/news', 'Загружены все новости', true)
    }

    const getLastNews = () => {
      fetchNews('/news/last', 'Загружена последняя новость', false)
    }

    const getLatestNews = () => {
      if (!limitCount.value || limitCount.value < 1) {
        error.value = 'Укажите корректное количество новостей'
        return
      }
      fetchNews(
          `/news/latest?limit=${limitCount.value}`,
          `Загружены последние ${limitCount.value} новостей`,
          true
      )
    }

    const getNewsById = () => {
      if (!newsId.value || newsId.value < 1) {
        error.value = 'Укажите корректный ID новости'
        return
      }
      fetchNews(
          `/news/${newsId.value}`,
          `Загружена новость с ID ${newsId.value}`,
          false
      )
    }

    return {
      // Реактивные переменные
      news,
      loading,
      error,
      info,
      limitCount,
      newsId,
      currentPage,
      itemsPerPage,

      // Вычисляемые свойства
      totalPages,
      paginatedNews,

      // Методы
      getAllNews,
      getLastNews,
      getLatestNews,
      getNewsById
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