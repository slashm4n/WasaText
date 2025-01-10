<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['selectedConversationChanged'],
	props: ['session_token', 'user', 'sent_message_to_user_name'],
	data() {
		return {
			errormsg: '',
			conversations: [],
			selected_conversation_id: 0
		}
	},
	methods: {
    	async loadConversationsList() {
			if (this.session_token == 0) {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'get',
					url: '/conversations',
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				this.conversations = response.data;
			} catch (e) {
				this.errormsg = "Error: " + e;
			}
			this.selected_conversation_id = 0;
		},
		async onConversationClick(conv) {
			console.log(`Clicked on ${conv.user_or_group_name}`);
			this.selected_conversation_id = conv.conversation_id;
			this.$emit('selectedConversationChanged', this.selected_conversation_id);
		}
  	},
  	watch: {
    	session_token(newToken, oldToken) {
      		this.loadConversationsList();
		},
    	sent_message_to_user_name(old_to_user, new_to_user) {
      		this.loadConversationsList();
			// TO DO: selezionare la conversazione
		}
	}
}
</script>

<template>
	<div class="conversations-container" v-if="session_token != 0">
		<div class="conv-box" v-for="c in conversations" :key="c.id"  @click="onConversationClick(c)">
			<img v-if="c.user_or_group_photo != ''" class="photo-box" v-bind:src="c.user_or_group_photo"></img>
			<img v-if="c.user_or_group_photo == ''" class="photo-box" src="../assets/profile.png"></img>
			<div class="name-box">{{ c.user_or_group_name }}</div>
		</div>
	</div>
	<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
