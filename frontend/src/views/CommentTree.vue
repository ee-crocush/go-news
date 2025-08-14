<template>
	<v-card class="my-3" elevation="2" :class="{ 'ml-8': isReply }">
		<v-card-text>
			<div class="d-flex justify-space-between align-start">
				<div>
					<strong>{{ comment.username }}</strong>
					<v-chip class="ml-2" x-small color="grey lighten-3">
						{{ comment.pub_time }}
					</v-chip>
				</div>
				<v-btn @click="toggleReplyForm" color="primary" text small>
					Ответить
				</v-btn>
			</div>

			<p class="mt-2 mb-0">{{ comment.content }}</p>

			<!-- Форма ответа -->
			<div v-if="showReplyForm" class="mt-3">
				<v-form @submit.prevent="submitReply">
					<v-text-field
						v-model="replyForm.username"
						label="Имя пользователя"
						required
						outlined
						dense
						hide-details></v-text-field>
					<v-textarea
						v-model="replyForm.content"
						label="Ответ"
						required
						outlined
						rows="2"
						class="mt-2"
						hide-details></v-textarea>
					<div class="mt-2">
						<v-btn
							type="submit"
							color="primary"
							small
							:loading="submitting"
							:disabled="!replyForm.username || !replyForm.content">
							Ответить
						</v-btn>
						<v-btn @click="cancelReply" color="grey" text small class="ml-2">
							Отмена
						</v-btn>
					</div>
				</v-form>
			</div>
		</v-card-text>

		<!-- Дочерние комментарии -->
		<div v-if="comment.children && comment.children.length > 0">
			<comment-tree
				v-for="childComment in comment.children"
				:key="childComment.id"
				:comment="childComment"
				:news-id="newsId"
				:is-reply="true"
				@reply="$emit('reply', $event)" />
		</div>
	</v-card>
</template>

<script>
import { ref, reactive } from 'vue'

export default {
	name: 'CommentTree',
	props: {
		comment: {
			type: Object,
			required: true,
		},
		newsId: {
			type: [String, Number],
			required: true,
		},
		isReply: {
			type: Boolean,
			default: false,
		},
	},
	emits: ['reply'],
	setup(props, { emit }) {
		const showReplyForm = ref(false)
		const submitting = ref(false)
		const replyForm = reactive({
			username: '',
			content: '',
		})

		const resetForm = () => {
			replyForm.username = ''
			replyForm.content = ''
		}

		const toggleReplyForm = () => {
			showReplyForm.value = !showReplyForm.value
			if (showReplyForm.value) {
				resetForm()
			}
		}

		const cancelReply = () => {
			showReplyForm.value = false
			resetForm()
		}

		const submitReply = () => {
			if (!replyForm.username || !replyForm.content) {
				return
			}

			submitting.value = true
			emit('reply', {
				news_id: props.newsId,
				parent_id: props.comment.id,
				username: replyForm.username,
				content: replyForm.content,
			})

			showReplyForm.value = false
			resetForm()
			submitting.value = false
		}

		return {
			showReplyForm,
			submitting,
			replyForm,
			toggleReplyForm,
			cancelReply,
			submitReply,
		}
	},
}
</script>

<style scoped>
.ml-8 {
	margin-left: 32px;
}
</style>
