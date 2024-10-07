<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAxios } from '@/services/useAxios'
import UserMiniCard from '@/components/UserMiniCard.vue'
import ErrorMsg from '@/components/ErrorMsg.vue'


const props = defineProps({
	searchValue: String
})

const router = useRouter()
const axios = useAxios()

const users = ref([])
const errormsg = ref(null)
const loadSearchedUsers = async () => {
	errormsg.value = null
	users.value = []

	if (!props.searchValue || props.searchValue.includes("?") || props.searchValue.includes("_")) {
		return
	}

	try {
		const response = await axios.get("/users", {
			params: { id: props.searchValue }
		})
		users.value = response.data
	} catch (e) {
		errormsg.value = e.toString()
	}
}
const goToProfile = (profileId) => {
	router.push(`/users/${profileId}`)
}
watch(() => props.searchValue, loadSearchedUsers)

onMounted(async () => {
	if (!localStorage.getItem('token')) {
		await router.replace("/login")
		return
	}
	await loadSearchedUsers()
})

</script>

<template>
	<div class="container-fluid h-100 ">
		<UserMiniCard
			v-for="user in users"
			:key="user.user_id"
			:identifier="user.user_id"
			:nickname="user.nickname"
			@clicked-user="goToProfile"
		/>

		<p v-if="users.length === 0" class="no-result-text d-flex justify-content-center"> No users found.</p>
		<ErrorMsg v-if="errormsg" :msg="errormsg"/>
	</div>
</template>

<style scoped>
.no-result-text {
	color: white;
	font-style: italic;
	margin-top: 1rem;
}
</style>
