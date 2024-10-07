<script setup>
import { ref, onMounted } from 'vue'
import { useAxios } from '@/services/useAxios'
import ErrorMsg from '@/components/ErrorMsg.vue'
import Photo from '@/components/Photo.vue'



const axios = useAxios()

const errormsg = ref(null)
const photos = ref([])
const loadStream = async () => {
	try {
		errormsg.value = null
		const token = localStorage.getItem('token')
		const response = await axios.get(`/users/${token}/home`)

		photos.value = (response.data || []).map(photo => {
			// Se upload_date non è disponibile, usiamo la data corrente come fallback
			const upload_date = photo.upload_date ? new Date(photo.upload_date) : new Date()
			console.log('Data di upload per foto:', photo.photo_id, upload_date)
			return {
				...photo,
				upload_date
			}
		})
	} catch (e) {
		errormsg.value = e.toString()
	}
}

onMounted(loadStream)

</script>

<template>
	<div class="container-fluid">
		<ErrorMsg v-if="errormsg" :msg="errormsg"/>

		<div class="row">
			<Photo
				v-for="photo in photos"
				:key=" photo.photo_id"
				:owner="photo.owner"
				:photo_id="photo.photo_id"
				:comments="photo.comments ?? []"
				:likes="photo.likes ?? []"
				:upload_date="photo.upload_date"
			/>
		</div>

		<div v-if="photos.length === 0" class="row ">
			<h1 class="d-flex justify-content-center mt-5" style="color: white;">There's no content yet, follow somebody!</h1>
		</div>
	</div>
</template>

<style>
</style>
