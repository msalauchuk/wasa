<script>
export default {
	data: function() {
		return {
			errormsg: null,
			identifier: "",
			disabled: true,
		}
	},
	methods: {
		async login() {
			// Clear previous error message
			this.errormsg = null;
			try {
				// Make a POST request to the /session endpoint with user_id
				let response = await this.$axios.post("/session",{
					user_id: this.identifier.trim()
				});

				localStorage.setItem('token',response.data.user_id);
				this.$router.replace("/home")
				this.$emit('updatedLoggedChild',true)

			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
		},

	},
	mounted(){
		if (localStorage.getItem('token')){
			this.$router.replace("/home")
		}
	},

}
</script>

<template>
	<div class="container-fluid h-100 m-0 p-0 login">
		<div class="row">
			<div class="col-12">
				<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			</div>
		</div>

		<div class="row h-100 align-items-center justify-content-center m-0">
			<div class="col-12 col-sm-8 col-md-4 col-lg-3 col-xl-2">
				<form @submit.prevent="login" class="d-flex flex-column align-items-center p-0">
					<h3 class="login-title mb-3">WASAphoto</h3>
					<input
						type="text"
						class="form-control mb-3"
						v-model="identifier"
						maxlength="16"
						minlength="3"
						placeholder="username" />
					<button class="btn btn-dark w-100"
							:disabled="identifier == null || identifier.length >16 || identifier.length <3 || identifier.trim().length<3">
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
