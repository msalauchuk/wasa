<script setup>


import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAxios } from '@/services/useAxios'  // Importa il composable useAxios
import ErrorMsg from '@/components/ErrorMsg.vue'
import PageNotFound from '@/components/PageNotFound.vue'
import Photo from '@/components/Photo.vue'

const route = useRoute()
const router = useRouter()
const axios = useAxios()  // Usa il composable useAxios

const errormsg = ref(null)
const userState = reactive({
	exists: false,
	banned: false,
	currentUserBanned: false,
	nickname: "",
	followStatus: false,
	followerCount: 0,
	followingCount: 0,
	postCount: 0,
	photos: [],
	following: [],
	followers: [],
})

// cambio nickname
const successMsg = ref('')

function getToken() {
	return localStorage.getItem('token')
}

const currentPath = computed(() => route.params.id)
const sameUser = computed(() => currentPath.value === getToken())


watch(currentPath, (newId, oldId) => {
	if (newId !== oldId) {
		loadInfo()
	}
})

async function uploadFile(event) {
	const file = event.target.files[0]
	if (!file) return

	const reader = new FileReader()
	reader.onload = async () => {
		try {
			const response = await axios.post(`/users/${currentPath.value}/photos`, reader.result, {
				headers: { 'Content-Type': file.type }
			})
			userState.photos.unshift(response.data)
			userState.postCount++
		} catch (error) {
			errormsg.value = error.toString()
		}
	}
	reader.readAsArrayBuffer(file)
}
async function toggleFollow() {
	try {
		if (userState.followStatus) {
			await axios.delete(`/users/${currentPath.value}/followers/${getToken()}`)
			userState.followerCount--
		} else {
			await axios.put(`/users/${currentPath.value}/followers/${getToken()}`)
			userState.followerCount++
		}
		userState.followStatus = !userState.followStatus
	} catch (error) {
		errormsg.value = error.toString()
	}
}
async function toggleBan() {
	try {
		if (userState.banned) {
			await axios.delete(`/users/${getToken()}/banned_users/${currentPath.value}`)
			await loadInfo()
		} else {
			await axios.put(`/users/${getToken()}/banned_users/${currentPath.value}`)
			userState.followStatus = false
		}
		userState.banned = !userState.banned
	} catch (error) {
		errormsg.value = error.toString()
	}
}

async function loadInfo() {
	if (!currentPath.value) return

	try {
		const response = await axios.get(`/users/${currentPath.value}`)

		userState.banned = false
		userState.exists = true
		userState.currentUserBanned = false

		//partial content
		if (response.status === 206) {
			userState.banned = true
			return
		}

		// no content
		if (response.status === 204) {
			userState.exists = false
			return
		}

		const { nickname, followers, following, posts } = response.data
		userState.nickname = nickname
		userState.followerCount = followers?.length ?? 0
		userState.followingCount = following?.length ?? 0
		userState.postCount = posts?.length ?? 0
		userState.followStatus = followers?.some(f => f.user_id === getToken()) ?? false
		userState.photos = posts?.map(post => ({
			...post,
			upload_date: new Date(post.upload_date || post.created_at || new Date())
		})) ?? []
		userState.followers = followers ?? []
		userState.following = following ?? []

	} catch (error) {
		userState.currentUserBanned = true
		errormsg.value = error.toString()
	}
}

function goToSettings() {
	router.push(`${currentPath.value}/settings`)
}

function removePhotoFromList(photo_id) {
	userState.photos = userState.photos.filter(photo => photo.photo_id !== photo_id)
	userState.postCount--
}

onMounted(async() => {
	await loadInfo()

	const message = localStorage.getItem('nicknameUpdateSuccess')
	if (message) {
		successMsg.value = message
		localStorage.removeItem('nicknameUpdateSuccess')

		setTimeout(() => {
			successMsg.value = ''
		}, 5000)
	}
})

</script>

<template>

    <div class="container-fluid" v-if="!userState.currentUserBanned && userState.exists">
		<div v-if="successMsg" class="row mt-3">
			<div class="col-12 d-flex justify-content-center">
				<div class="alert alert-success">{{ successMsg }}</div>
			</div>
		</div>
        <div class="row">
            <div class="col-12 d-flex justify-content-center">
                <div class="card w-50 container-fluid">

                    <div class="row">
                        <div class="col">
                            <div class="card-body d-flex justify-content-between align-items-center">
                                <h5 class="card-title p-0 me-auto mt-auto">{{userState.nickname}} @{{currentPath}}</h5>
                                <button v-if="!sameUser && !userState.banned" @click="toggleFollow" class="btn btn-success ms-2">
                                    {{userState.followStatus ? "Unfollow" : "Follow"}}
                                </button>

                                <button v-if="!sameUser" @click="toggleBan" class="btn btn-danger ms-2">
                                    {{userState.banned ? "Unban" : "Ban"}}
                                </button>

                                <button v-else class="my-trnsp-btn ms-2" @click="goToSettings">
                                    <!--Settings  <font-awesome-icon icon="fa-solid fa-gear" />-->
                                    <i class="my-nav-icon-gear fa-solid fa-gear"></i>
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-if="!userState.banned" class="row mt-1 mb-1">
                        <div class="col-4 d-flex justify-content-start">
                            <h6 class="ms-3 p-0 ">Posts: {{userState.postCount}}</h6>
                        </div>

                        <div class="col-4 d-flex justify-content-center">
                            <h6 class=" p-0 ">Followers: {{userState.followerCount}}</h6>
                        </div>

                        <div class="col-4 d-flex justify-content-end">
                            <h6 class=" p-0 me-3">Following: {{userState.followingCount}}</h6>
                        </div>
                    </div>
                </div>
            </div>
        </div>


        <div class="row">

            <div class="container-fluid mt-3">

                <div class="row ">
                    <div class="col-12 d-flex justify-content-center">
                        <h2>Posts</h2>
                        <input id="fileUploader" type="file" class="profile-file-upload" @change="uploadFile" accept=".jpg, .png">
                        <label v-if="sameUser" class="btn my-btn-add-photo ms-2 d-flex align-items-center" for="fileUploader"> Add </label>
                    </div>
                </div>

                <div class="row ">
                    <div class="col-3"></div>
                    <div class="col-6">
                        <hr class="border border-dark">
                    </div>
                    <div class="col-3"></div>
                </div>
            </div>
        </div>

        <div class="row">
            <div class="col">

                <div v-if="!userState.banned && userState.postCount>0">
					<Photo v-for="photo in userState.photos"
						   :key="photo.photo_id"
						   :owner="currentPath"
						   v-bind="photo"
						   :isOwner="sameUser"
						   @removePhoto="removePhotoFromList"
					/>

                </div>

                <div v-else class="mt-5 ">
                    <h2 class="d-flex justify-content-center" style="color: white;">No posts yet</h2>
                </div>

            </div>
        </div>


    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
    <PageNotFound v-else class="h-25 "/>

</template>

<style>
.profile-file-upload{
    display: none;
}

.my-nav-icon-gear{
    color: grey;
}
.my-nav-icon-gear:hover{
    transform: scale(1.3);
}

.my-btn-add-photo{
    background-color: #fffefe;
    border-color: grey;
}
.my-btn-add-photo:hover{
    color: white;
    background-color: lightgrey;
    border-color: grey;
}

</style>
