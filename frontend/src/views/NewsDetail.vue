<template>
  <div>
    <!-- Кнопка назад -->
    <v-btn
        @click="goBack"
        color="primary"
        class="ma-5"
        outlined
    >
      ← Назад к списку новостей
    </v-btn>

    <!-- Загрузка -->
    <div v-if="loading" class="text-center ma-5">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
      <p>Загрузка новости...</p>
    </div>

    <!-- Ошибка -->
    <v-alert
        v-if="error"
        type="error"
        dismissible
        @input="error = ''"
        class="mx-5 my-3"
    >
      {{ error }}
    </v-alert>

    <!-- Новость -->
    <div v-if="post && !loading">
      <v-card class="mx-5 my-5" elevation="5">
        <v-card-title>
          <h2>{{ post.title }}</h2>
        </v-card-title>
        <v-card-text>
          <p>{{ post.content }}</p>
          <v-card-subtitle class="mt-3">
            ID: {{ post.id }} | {{ post.pub_time }}
          </v-card-subtitle>

          <div class="mt-3" v-if="post.link">
            <v-btn
                color="primary"
                :href="post.link"
                target="_blank"
            >
              Перейти к источнику
            </v-btn>
          </div>
        </v-card-text>
      </v-card>

      <!-- Уведомление -->
      <v-snackbar
          v-model="snackbar"
          :timeout="4000"
          color="success"
          top
          right
          multi-line
      >
        {{ successMessage }}
        <template #action="{ attrs }">
          <v-btn
              color="white"
              text
              v-bind="attrs"
              @click="snackbar = false"
          >
            Закрыть
          </v-btn>
        </template>
      </v-snackbar>

      <!-- Секция комментариев -->
      <v-card class="mx-5 my-5" elevation="3">
        <v-card-title>
          Комментарии ({{ totalComments }})
        </v-card-title>

        <!-- Форма добавления нового комментария -->
        <v-card-text>
          <v-form @submit.prevent="addComment">
            <v-text-field
                v-model="newComment.username"
                label="Имя пользователя"
                required
                outlined
                dense
            ></v-text-field>
            <v-textarea
                v-model="newComment.content"
                label="Комментарий"
                required
                outlined
                rows="3"
            ></v-textarea>
            <v-btn
                type="submit"
                color="primary"
                :loading="submittingComment"
                :disabled="!newComment.username || !newComment.content"
            >
              Добавить комментарий
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>

      <!-- Дерево комментариев -->
      <div v-if="comments.length > 0" class="mx-5">
        <comment-tree
            v-for="comment in comments"
            :key="comment.id"
            :comment="comment"
            :news-id="post.id"
            @reply="handleReply"
        />
      </div>

      <div v-else-if="!loading" class="text-center ma-5">
        <p>Комментариев пока нет. Будьте первым!</p>
      </div>
    </div>
  </div>
</template>

<script>
import {computed, onMounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import CommentTree from './CommentTree.vue'

export default {
  name: "NewsDetail",
  components: {
    CommentTree
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    // Реактивные переменные
    const post = ref(null)
    const comments = ref([])
    const loading = ref(false)
    const error = ref('')
    const submittingComment = ref(false)
    const successMessage = ref('')
    const snackbar = ref(false)

    const newComment = reactive({
      username: '',
      content: '',
      parent_id: null
    })

    // Получение baseUrl из переменных окружения
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

    const newsId = computed(() => route.params.id)

    const totalComments = computed(() => {
      const countComments = (commentsList) => {
        let count = commentsList.length
        commentsList.forEach(comment => {
          if (comment.children && comment.children.length > 0) {
            count += countComments(comment.children)
          }
        })
        return count
      }
      return countComments(comments.value)
    })

    // Методы
    const fetchNewsDetail = async () => {
      loading.value = true
      error.value = ''

      try {
        const response = await fetch(`${baseUrl}/news/${newsId.value}`)

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }

        const data = await response.json()

        post.value = data.data.post
        comments.value = data.data.comments || []

      } catch (err) {
        error.value = `Ошибка при загрузке новости: ${err.message}`
      } finally {
        loading.value = false
      }
    }

    const addComment = async () => {
      if (!newComment.username || !newComment.content) {
        return
      }

      submittingComment.value = true
      error.value = ''
      successMessage.value = ''

      try {
        const response = await fetch(`${baseUrl}/comments`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            news_id: Number(newsId.value),
            username: newComment.username,
            content: newComment.content,
            parent_id:
                newComment.parent_id !== null
                    ? Number(newComment.parent_id)
                    : null
          })
        })

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }

        // Очищаем форму
        newComment.username = ''
        newComment.content = ''
        newComment.parent_id = null

        successMessage.value = 'Сообщение успешно создано. После модерации будет опубликовано.'
        snackbar.value = true

      } catch (err) {
        error.value = `Ошибка при добавлении комментария: ${err.message}`
      } finally {
        submittingComment.value = false
      }
    }

    const handleReply = ({parent_id, username, content}) => {
      // Обработка ответа на комментарий
      submitReply(parent_id, username, content)
    }

    const submitReply = async (parent_id, username, content) => {
      error.value = ''
      successMessage.value = ''
      try {
        const response = await fetch(`${baseUrl}/comments`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            news_id: Number(newsId.value),
            username: username,
            content: content,
            parent_id: parent_id !== null ? Number(parent_id) : null
          })
        })

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }

        successMessage.value = 'Сообщение успешно создано. После модерации будет опубликовано.'
        snackbar.value = true

      } catch (err) {
        error.value = `Ошибка при добавлении ответа: ${err.message}`
      }
    }

    const goBack = () => {
      router.push('/')
    }

    // Загружаем новость при монтировании компонента
    onMounted(() => {
      fetchNewsDetail()
    })

    return {
      post,
      comments,
      loading,
      error,
      submittingComment,
      newComment,
      totalComments,
      successMessage,
      snackbar,
      // Методы
      addComment,
      handleReply,
      goBack
    }
  }
}
</script>

<style scoped>
.v-card {
  margin-bottom: 20px;
}
</style>