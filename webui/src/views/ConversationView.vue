<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	props: ['session_token', 'user', 'selected_conversation_id'],
	data: function() {
		return {
			errormsg: '',
			messages: [],
		}
	},
	methods: {
		async loadConversation() {
			try {
				if (this.selected_conversation_id == 0)
				{
					return;
				}

				const response = await this.$axios({
					method: 'get',
					url: '/conversations/' + this.selected_conversation_id,
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				if (response.status != 200) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

				this.messages = response.data;
			} catch (e) {
				this.errormsg = "Error: " + e;
			}
		},
	},
	watch: {
    	selected_conversation_id(newConversationId, oldConversationId) {
      		this.loadConversation()
		}
	},
}
</script>

<template>
	<div v-if="session_token != 0 && user != null && selected_conversation_id != 0" class="conversation-container">
		<div :class="{'message-box friend-message':msg.from_user_id!=user.id,'message-box my-message':msg.from_user_id==user.id}"
		      v-for="msg in messages" :key="msg.id">
			<p>{{ msg.msg }} {{ msg.sent_timestamp }}</p>
		</div>
	</div>

	<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
