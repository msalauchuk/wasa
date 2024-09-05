<script>
export default {
	data: function () {
		return {
			errormsg: null,
			nickname: "",
			successmsg: null,  //success message

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
				this.successmsg = null;
				return;
			}
			try{
				// Nickname put: /users/:id
				let response = await this.$axios.put("/users/"+this.$route.params.id,{
					nickname: this.nickname,
				});
				if (response.status === 200) {
					this.nickname = ""; // Clear the input after successful update
					this.successmsg = "Nickname updated successfully!";
					this.errormsg = null;
				}

			} catch (e){
				this.errormsg = e.toString();
				this.successmsg = null;
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

		<div class="row text-center">
			<div class="col-12">
				<p class="mt-2">Please note, you can only change the <span class="text-success">nickname</span> (the part before "@"). The <span class="fw-bold">identifier</span> (the part after "@") is fixed and cannot be altered.</p>
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

		<!-- Display success message -->
		<div class="row" v-if="successmsg">
			<div class="col d-flex justify-content-center">
				<div class="alert alert-success">{{ successmsg }}</div>
			</div>
		</div>

		<!-- Display nickname preview -->
		<div class="row">
			<div v-if="nickname.trim().length > 0" class="col d-flex justify-content-center">
				Preview: {{ nickname }} @{{ this.$route.params.id }}
			</div>
		</div>


		<!-- Display error message -->
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
