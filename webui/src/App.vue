<script setup>
import { ref, onMounted, onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './components/Navbar.vue'

const isLoggedIn = ref(false)
const searchValue = ref("")

const router = useRouter()

const handleLogout = () => {
	isLoggedIn.value = false
	router.replace("/login")
}

const updateLoginStatus = (status) => {
	isLoggedIn.value = status
}

const navigateToRoute = (route) => {
	router.replace(route)
}

const handleSearch = (queryParam) => {
	searchValue.value = queryParam
	router.replace("/search")
}


// Lifecycle hook creazione iniziale
// Controlla se è la prima volta che l'utente accede e cancella il localStorage se è così
onBeforeMount(() => {
	if (!localStorage.getItem('notFirstStart')) {
		localStorage.clear()
		localStorage.setItem('notFirstStart', true)
	}
})

// Lifecycle hook montaggio del componente
onMounted(() => {
	if (!localStorage.getItem('token')) {
		router.replace("/login") // Se non c'è token, manda al login
	} else {
		isLoggedIn.value = true
	}
})
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col p-0">
				<main >
					<Navbar v-if="isLoggedIn"
					@logoutNavbar="handleLogout"
					@requestUpdateView="navigateToRoute"
					@searchNavbar="handleSearch"/>

					<router-view v-slot="{ Component }">
						<transition name="fade" mode="out-in">
							<!-- rendering dinamico della componente v-bind:is -->
							<component :is="Component"
									   @updatedLoggedChild="updateLoginStatus"
									   @requestUpdateView="navigateToRoute"
									   :searchValue="searchValue"/>
						</transition>
					</router-view>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
	opacity: 0;
}
</style>
