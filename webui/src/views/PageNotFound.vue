<script setup>

import {ref} from 'vue'

const message = ref("Oops! The requested " +
	"page doesn't exist. The link might be broken or content might've been deleted.")

</script>

<template>
	<div class="page-not-found">
		<p class="page-not-found__message">{{ message }}</p>
	</div>
</template>

<style>

.page-not-found {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100vh;
	background-color: #f8f9fa;
}

.page-not-found__message {
	color: #6c757d;
	font-size: 1.25rem;
	text-align: center;
	max-width: 600px;
	padding: 1rem;
	background-color: white;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
