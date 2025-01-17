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
		async doLoadConversation() {
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

				// TO DO: scroll to end non funziona
				// UPDATE: non più necessario conversazione in reverse chronological order
				// this.$refs.convview.scrollIntoView();
				this.errormsg = "";
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
		async onMessageClick(msg) {
			this.selected_message_id = msg.msg_id;
			this.$emit('selectedMessageChanged', this.selected_message_id);
		}
	},
	watch: {
    	session_token(newValue, oldValue) {
			if (newValue) {
				this.errormsg = '';
			}
		},
    	selected_conversation_id(newValue, oldValue) {
      		this.doLoadConversation();
		},
    	need_update_conversation(newValue, oldValue) {
      		this.doLoadConversation();
			this.$emit('conversationUpdated');
		}
	},
}
</script>

<template>
	<div v-if="session_token != 0">
		<div ref="convview" v-if="session_token != 0 && user != null && selected_conversation_id != 0" class="conversation-container">
			<div :class="{'message-box friend-message':msg.from_user_id!=user.id,'message-box my-message':msg.from_user_id==user.id}"
				v-for="(msg, index) in messages" :key="msg.id" :tabindex="index"
				@click="onMessageClick(msg)">
				<p v-if=!msg.is_photo>{{ msg.msg }}<span style="font-size: smaller;">{{ msg.sent_timestamp }}</span><span style="font-size: medium">{{ msg.reaction }}</span></p>
				<p v-if=msg.is_photo><img class="photo-box-big" v-bind:src=msg.msg></img><span style="font-size: smaller;">{{ msg.sent_timestamp }}</span><span style="font-size: medium">{{ msg.reaction }}</span></p>
			</div>
		</div>
		
		<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
	</div>
</template>
