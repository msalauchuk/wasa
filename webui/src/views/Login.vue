<script setup>
import { ref, computed, onMounted, getCurrentInstance } from 'vue'
import { useRouter } from 'vue-router'

const app = getCurrentInstance()
const axios = app.appContext.config.globalProperties.$axios
const router = useRouter()
const emit = defineEmits(['updatedLoggedChild']);
const userId = ref('')
const error = ref(null)
const isLoading = ref(false)


const handeLogin = async() =>{
	error.value = null // Resetta gli errori precedenti
	isLoading.value = true  //il caricamento è in corso

	try {
		const response = await axios.post("/session",
			{user_id:userId.value.trim()})

		localStorage.setItem('token', response.data.user_id)
		router.replace('/home')
		emit('updatedLoggedChild', true) // Comunica al parent

	}catch (err){
		error.value = err.toString()
	}
}

const isButtonDisabled = computed(() => {
	const trimmedLength = userId.value.trim().length
	return !userId.value || trimmedLength < 3 || userId.value.length > 16
})

// Se c'è un token, redirection alla home
onMounted(()=>{
	if (localStorage.getItem('token')){
		router.replace('/home')
	}
})
</script>

<template>
	<div class="container-fluid h-100 m-0 p-0 login">
		<div class="row">
			<div class="col-12">
				<ErrorMsg v-if="error" :msg="error"></ErrorMsg>
			</div>
		</div>

		<div class="row h-100 align-items-center justify-content-center m-0">
			<div class="col-12 col-sm-8 col-md-4 col-lg-3 col-xl-2">
				<form @submit.prevent="handeLogin" class="d-flex flex-column align-items-center p-0">
					<h3 class="login-title mb-3">WASAphoto</h3>
					<input
						type="text"
						class="form-control mb-3"
						v-model="userId"
						maxlength="16"
						minlength="3"
						placeholder="username"/>
					<button class="btn btn-dark w-100"
							:disabled=isButtonDisabled>
						Log in
					</button>
				</form>
			</div>
		</div>
	</div>
</template>

<style>
.login {
    background-image: url("../assets/images/white.png");
    height: 100vh;
}

.login-title {
	font-family: 'title';
	color: black;
	font-size: 40px;
}


.login .btn-dark {
	background-color: #6299f6; /* Example: a deep blue */
	color: #ffffff; /* Ensures the text color is white */
	border: lightgrey;
}

.login .btn-dark:hover {
	background-color: #6299f6; /* Darker shade for hover state */
}
.form-control {
	width: 100%; /* Ensure the input takes full width of its parent */
}
</style>
