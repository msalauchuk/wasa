<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAxios } from '@/services/useAxios'  // Importa il composable useAxios



const props = defineProps({
	owner: String,
	likes: Array,
	comments: Array,
	upload_date: Date,
	photo_id: Number,
	isOwner: Boolean
})

console.log('Props ricevute in Photo:', props)


const emit = defineEmits(['removePhoto'])

const axios = useAxios()
const router = useRouter()

const photoURL = ref('')
const liked = ref(false)
const allComments = ref([])
const allLikes = ref([])

const loadPhoto = () => {
	photoURL.value = `${__API_URL__}/users/${props.owner}/photos/${props.photo_id}`
}

const deletePhoto = async () => {
	try {
		await axios.delete(`/users/${props.owner}/photos/${props.photo_id}`)
		emit('removePhoto', props.photo_id)
	} catch (e) {
		console.error('Errore durante l\'eliminazione della foto:', e)
	}
}

const photoOwnerClick = () => {
	router.replace(`/users/${props.owner}`)
}

const toggleLike = async () => {
	if (props.isOwner) return

	const bearer = localStorage.getItem('token')
	const endpoint = `/users/${props.owner}/photos/${props.photo_id}/likes/${bearer}`

	try {
		if (!liked.value) {
			await axios.put(endpoint)
			allLikes.value.push({ user_id: bearer, nickname: bearer })
		} else {
			await axios.delete(endpoint)
			allLikes.value = allLikes.value.filter(like => like.user_id !== bearer)
		}
		liked.value = !liked.value
	} catch (e) {
		console.error('Errore durante il toggle del like:', e)
	}
}

const removeCommentFromList = (commentId) => {
	allComments.value = allComments.value.filter(comment => comment.comment_id !== commentId)
}

const addCommentToList = (comment) => {
	allComments.value.push(comment)
}

onMounted(async () => {
	await loadPhoto()
	if (props.likes) {
		allLikes.value = props.likes
		liked.value = allLikes.value.some(like => like.user_id === localStorage.getItem('token'))
	}
	if (props.comments) {
		allComments.value = props.comments
	}
})


const formattedUploadDate = computed(() => {
	if (!props.upload_date) return 'Data non disponibile'
	const date = new Date(props.upload_date)
	return date.toLocaleDateString('it-IT', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	})
})
const likesCount = computed(() => allLikes.value.length)
const commentsCount = computed(() => allComments.value.length)

</script>

<template>
	<div class="container-fluid mt-3 mb-5">
		<LikeModal :modal_id="`like_modal${photo_id}`" :likes="allLikes" />
		<CommentModal
			:modal_id="`comment_modal${photo_id}`"
			:comments_list="allComments"
			:photo_owner="owner"
			:photo_id="photo_id"
			@eliminateComment="removeCommentFromList"
			@addComment="addCommentToList"
		/>


		<div class="d-flex flex-row justify-content-center">
			<div class="card my-card">
				<div class="d-flex justify-content-end">
					<button v-if="isOwner" class="my-trnsp-btn my-dlt-btn me-2" @click="deletePhoto">
						<i class="fa-solid fa-trash w-100 h-100"></i>
					</button>
				</div>
				<div class="d-flex justify-content-center photo-background-color">
					<img :src="photoURL" class="card-img-top img-fluid" alt="Photo">
				</div>
				<div class="card-body">
					<div class="container">
						<div class="d-flex flex-row justify-content-end align-items-center mb-2">
							<button class="my-trnsp-btn m-0 p-1 me-auto" @click="photoOwnerClick">
								<i>From {{ owner }}</i>
							</button>

							<button class="my-trnsp-btn m-0 p-1 d-flex justify-content-center align-items-center">
								<i @click="toggleLike" :class="['me-1', 'my-heart-color', 'w-100', 'h-100', 'fa', liked ? 'fa-heart' : 'fa-heart-o']"></i>
								<i data-bs-toggle="modal" :data-bs-target="`#like_modal${photo_id}`" class="my-comment-color">
									{{ likesCount }}
								</i>
							</button>

							<button
								class="my-trnsp-btn m-0 p-1 d-flex justify-content-center align-items-center"
								data-bs-toggle="modal"
								:data-bs-target="`#comment_modal${photo_id}`"
							>
								<i class="my-comment-color fa-regular fa-comment me-1"></i>
								<i class="my-comment-color-2">{{ commentsCount }}</i>
							</button>
                        </div>

						<div class="d-flex flex-row justify-content-start align-items-center mt-2">
							<p class="upload-date-text mb-0">
								<i class="fas fa-calendar-alt me-1"></i>
								Uploaded on {{ formattedUploadDate }}
							</p>
						</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.photo-background-color{
	background-color: grey;
}

.my-card{
	width: 27rem;
	border-color: black;
	border-width: thin;
}

.my-heart-color{
	color: grey;
}
.my-heart-color:hover{
	color: red;
}

.my-comment-color {
	color: grey;
}
.my-comment-color:hover{
	color: black;
}

.my-comment-color-2{
	color:grey
}

.my-dlt-btn{
	font-size: 19px;
}
.my-dlt-btn:hover{
	font-size: 19px;
	color: var(--color-red-danger);
}
.upload-date-text {
	font-size: 0.7rem;  /* Ridotto ulteriormente */
	color: #9a9a9a;  /* Colore più chiaro */
	font-style: italic;
	font-weight: 300;  /* Font weight più leggero */
	letter-spacing: 0.02em;  /* Leggero spacing tra le lettere */
	opacity: 0.8;  /* Leggera trasparenza */
}

.upload-date-text i {
	font-size: 0.7rem;  /* Icona leggermente più piccola del testo */
}
</style>
