<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['selectedMessageChanged'],
	props: ['session_token', 'user', 'selected_conversation_id'],
	data: function() {
		return {
			errormsg: '',
			messages: [],
			selected_message_id: 0
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
		async onMessageClick(msg) {
			this.selected_message_id = msg.msg_id;
			this.$emit('selectedMessageChanged', this.selected_message_id);}
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
		      v-for="msg in messages" :key="msg.id"
		      @click="onMessageClick(msg)">
			<p><div>{{ msg.msg }}</div><div style="font-size: smaller;">{{ msg.sent_timestamp }}</div></p>
		</div>
	</div>

	<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
