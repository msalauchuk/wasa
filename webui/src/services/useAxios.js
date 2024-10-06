import { inject } from 'vue'

export function useAxios() {
	const axios = inject('axios')
	if (!axios) {
		throw new Error('Axios non è stato fornito')
	}
	return axios
}
