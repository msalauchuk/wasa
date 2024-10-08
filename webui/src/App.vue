<script setup>
import { ref, onMounted, onBeforeMount, onUnmounted } from 'vue'
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


// prevenire lo spostamento dovuto alla barra di scorrimento
const preventScrollShift = () => {
	const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth
	document.documentElement.style.setProperty('--scrollbar-width', `${scrollbarWidth}px`)
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
	// Previeni lo spostamento dovuto alla barra di scorrimento
	preventScrollShift()
	window.addEventListener('resize', preventScrollShift)
})

// Lifecycle hook smontaggio del componente
onUnmounted(() => {
	window.removeEventListener('resize', preventScrollShift)
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
:root {
	--scrollbar-width: 0px;
}

body {
	overflow-y: scroll;
	padding-right: var(--scrollbar-width);
}

</style>
