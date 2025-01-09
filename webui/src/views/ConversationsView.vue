<script setup>
import ConversationView from './ConversationView.vue';
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	props: ['session_token', 'user'],
	data() {
		return {
			errormsg: '',
			conversations: [],
			selected_conversation_id: 0
		}
	},
	methods: {
    	async loadConversationsList() {
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
				alert("Error: " + e);
			}
		},
		async onConversationClick(conv) {
			console.log(`Clicked on ${conv.user_or_group_name}`);
			this.selected_conversation_id = conv.conversation_id
		}
  	},
  	watch: {
    	session_token(newToken, oldToken) {
      		this.loadConversationsList()
		}
	}
}
</script>

<template>
	<div class="avolio-conversations-container">
		<div class="avolio-chat-list">
			<div class="avolio-chat-box" v-for="c in conversations" :key="c.id"  @click="onConversationClick(c)">
				<div class="img-box"></div>
				<div class="avolio-chat-details">{{ c.user_or_group_name }}</div>
			</div>

		</div>
	</div>
	<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
    <ConversationView :session_token="session_token" :user="user" :conversation_id="selected_conversation_id"></ConversationView>
</template>
