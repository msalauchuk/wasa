<script setup>


import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAxios } from '@/services/useAxios'  // Importa il composable useAxios
import ErrorMsg from '@/components/ErrorMsg.vue'


const route = useRoute() // composable
const axios = useAxios()
const router = useRouter()


const errormsg = ref(null)
const nickname = ref('')


const isValidNickname = computed(() => {
	return nickname.value &&
		nickname.value.length >= 3 &&
		nickname.value.length <= 16 &&
		nickname.value.trim().length > 0
})

const modifyNickname = async () => {
	console.log('modifyNickname function called');

	if (!isValidNickname.value) {
		console.log('Invalid nickname');

		errormsg.value = "Invalid nickname. Please ensure it is between 3 and 16 characters long."
		return
	}

	try {
		console.log('Sending request to update nickname');

		const response = await axios.put(`/users/${route.params.id}`, {
			nickname: nickname.value
		})

		console.log('Response received:', response);
		if (response.status === 204) {
			nickname.value = ''
			errormsg.value = null


			console.log('Attempting to save to sessionStorage');

			// Usa sessionStorage invece di localStorage
			localStorage.setItem('nicknameUpdateSuccess', 'Nickname updated successfully!')
			console.log('Saved to sessionStorage:', localStorage.getItem('nicknameUpdateSuccess'))

			console.log('Attempting to redirect');
			// Reindirizza immediatamente
			router.push(`/users/${route.params.id}`)
		}
	} catch (e) {
		console.error('Error updating nickname:', e)
		errormsg.value = e.response?.data?.message || e.message || "An error occurred"
	}
}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col d-flex justify-content-center mb-2">
				<h1>{{route.params.id }}'s Settings</h1>
			</div>
		</div>

		<div class="row text-center">
			<div class="col-12">
				<p class="mt-2">Please note, you can only change the <span class="text-success">nickname</span> (the part before "@"). The <span class="fw-bold">identifier</span> (the part after "@") is fixed and cannot be altered.</p>
			</div>
		</div>

		<div class="row mt-2">
			<div class="col d-flex justify-content-center">
				<div class="input-group mb-3 w-50">
					<input
						type="text"
						class="form-control"
						placeholder="Enter your new nickname..."
						maxlength="16"
						minlength="3"
						v-model="nickname"
					/>
					<div class="input-group-append">
						<button class="btn btn-outline-secondary"
								@click="modifyNickname"
								:disabled="!isValidNickname">
							Modify
						</button>
					</div>
				</div>
			</div>
		</div>


		<!-- Display nickname preview -->
		<div class="row">
			<div v-if="nickname.trim().length > 0" class="col d-flex justify-content-center">
				Preview: {{ nickname }} @{{route.params.id }}
			</div>
		</div>


		<!-- Display error message -->
		<div class="row">
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>

<style scoped>
.container-fluid {
	padding-top: 20px;
}
</style>
