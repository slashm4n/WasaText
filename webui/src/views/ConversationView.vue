<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['selectedMessageChanged', 'conversationUpdated'],
	props: ['session_token', 'user', 'selected_conversation_id', 'need_update_conversation'],
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
                    this.errormsg = "Unexpected response " + response.status;
                    return;
                }

				this.messages = response.data;
			
				// TO DO: scroll to end NON FUNZIONA, PROBABILMENTE DA CAMBIARE POSIZIONE
				this.$refs.convview.scrollIntoView();

			} catch (e) {
				this.errormsg = "Error: " + e;
			}
		},
		async onMessageClick(msg) {
			this.selected_message_id = msg.msg_id;
			this.$emit('selectedMessageChanged', this.selected_message_id);
		}
	},
	watch: {
    	selected_conversation_id(newValue, oldValue) {
      		this.loadConversation();
		},
    	need_update_conversation(newValue, oldValue) {
      		this.loadConversation();
			this.$emit('conversationUpdated');
		}
	},
}
</script>

<template>
	<div ref="convview" v-if="session_token != 0 && user != null && selected_conversation_id != 0" class="conversation-container">
		<div :class="{'message-box friend-message':msg.from_user_id!=user.id,'message-box my-message':msg.from_user_id==user.id}"
		      v-for="(msg, index) in messages" :key="msg.id" :tabindex="index"
		      @click="onMessageClick(msg)">
			<p>{{ msg.msg }}<span style="font-size: smaller;">{{ msg.sent_timestamp }}</span><span style="font-size: xx-small">{{ msg.reaction }}</span></p>
		</div>
	</div>

	<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
