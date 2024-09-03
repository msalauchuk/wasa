<script>
export default {
	data: function () {
		return {
			errormsg: null,
			nickname: "",
		}
	},

	methods:{

		// Validation to ensure the nickname meets required criteria

		isValidNickname(nickname) {
			return nickname && nickname.length >= 3 && nickname.length <= 16 && nickname.trim().length > 0;
		},
		async modifyNickname(){
			if (!this.isValidNickname(this.nickname)) {
				this.errormsg = "Invalid nickname. Please ensure it is between 3 and 16 characters long.";
				return;
			}
			try{
				// Nickname put: /users/:id
				let response = await this.$axios.put("/users/"+this.$route.params.id,{
					nickname: this.nickname,
				});
				if (response.status === 200) {
					this.nickname = ""; // Clear the input after successful update
					alert("Nickname updated successfully!"); // Simple success alert
				}

			} catch (e){
				this.errormsg = e.toString();
			}
		},
	},

}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col d-flex justify-content-center mb-2">
				<h1>{{ this.$route.params.id }}'s Settings</h1>
			</div>
		</div>

		<div class="row">
			<div class="col-12 d-flex justify-content-center">
				<p class="me-1" style="color: var(--color-red-danger);">[Disclaimer] </p>
				<p> Your user handle is structured as <span class="ms-1 me-1" style="color: green;">nickname</span>@identifier. </p>
			</div>
			<div class="col-12 d-flex justify-content-center">
				<p>You can only modify the <span class="ms-1 me-1" style="color: green;">nickname</span> (before the "@") and not the identifier (after the "@").</p>
			</div>
			<div class="col-12 d-flex justify-content-center">
				<p>Note: 'Username' is synonymous with 'nickname' and refers to the same attribute.</p>
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
								:disabled="!isValidNickname(nickname)">
							Modify
						</button>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div v-if="nickname.trim().length > 0" class="col d-flex justify-content-center">
				Preview: {{ nickname }} @{{ this.$route.params.id }}
			</div>
		</div>

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
