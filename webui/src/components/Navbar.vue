<script setup>

import {ref, computed} from 'vue'

const emit = defineEmits(['logoutNavbar', 'requestUpdateView', 'searchNavbar'])


const searchText = ref('')
const icons = ref({
	home: 'fa-regular',
	profile: 'fa-regular'
})

const userId = computed(() => localStorage.getItem('token'))


const handleLogout = () => {
	localStorage.removeItem('token')
	emit('logoutNavbar', false)
}

const navigateTo = (route) => {
	emit('requestUpdateView', route)
}
const handleSearch = () => {
	if (!searchText.value.trim()) return
	emit('searchNavbar', searchText.value)
	searchText.value = ''
}

const handleIconHover = (icon, isHovered) => {
	icons.value[icon] = isHovered ? 'fa-solid' : 'fa-regular'
}

const navigationItems = computed(() => [
	{
		icon:' fa-house',
		action: () => navigateTo('/home'),
		onHover: (isHovered) => handleIconHover('home', isHovered),
		class: 'home-icon',
		noHover: true
	},
	{
		icon: icons.value.profile + ' fa-user',
		action: () => navigateTo(`/users/${userId.value}`),
		onHover: (isHovered) => handleIconHover('profile', isHovered),
		class: 'my-nav-icon-profile'
	},
	{
		icon: 'fa-solid fa-right-from-bracket',
		action: handleLogout,
		class: 'my-nav-icon-quit',
		noHover: true
	}
])
</script>

<template>

	<!-- Navbar component -->
	<nav class="navbar navbar-expand-lg sticky-top mb-3 bg-transparent">
		<div class="container-fluid">
			<!-- First column: Home icon -->
			<div class="col-4">
				<a
					v-for="(item, index) in navigationItems"
					:key="index"
					class="navbar-brand ms-2"
					@click="item.action"
					v-show="index === 0"
				>
					<i
						:class="['me-1 w-100 h-100', icons.home, item.icon, item.class]"
						@mouseover="() => !item.noHover && item.onHover(true)"
						@mouseout="() => !item.noHover && item.onHover(false)"
					/>
				</a>
			</div>

			<!-- Second column: Search bar -->
			<div class="col-4">
				<form @submit.prevent="handleSearch" class="d-flex justify-content-center">
					<input
						v-model="searchText"
						class="form-control w-50"
						type="search"
						placeholder="Search users"
					>
				</form>
			</div>

			<!-- Third column: Profile and Logout icons -->
			<div class="col-4 d-flex justify-content-end">
				<button
					v-for="(item, index) in navigationItems.slice(1)"
					:key="index"
					class="my-trnsp-btn me-2"
					@click="item.action"
				>
					<i
						:class="['me-1 w-100 h-100', item.icon, item.class]"
						@mouseover="() => !item.noHover && item.onHover(true)"
						@mouseout="() => !item.noHover && item.onHover(false)"
					/>
				</button>
			</div>
		</div>
	</nav>
</template>

<style scoped>
.navbar:hover {
	cursor: pointer;
}

.home-icon {
	font-size: 24px;
	color: white;
	transition: color 0.2s ease;
}

.home-icon:hover {
	color: black;
}

.my-nav-icon-profile,
.my-nav-icon-quit {
	color: black;
	transition: all 0.2s ease;
}

.my-nav-icon-quit:hover {
	color: var(--color-red-danger);
	transform: scale(1.2);
}

.my-trnsp-btn {
	background: none;
	border: none;
	padding: 0;
}
</style>
