<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>

export default {
	props: ['session_token', 'user', 'conversation_id'],
	data: function() {
		return {
			errormsg: '',
			messages: [],
		}
	},
	methods: {
		async loadConversation() {
			try {
				const response = await this.$axios({
					method: 'get',
					url: '/conversations/' + this.conversation_id,
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				this.messages = response.data;
			} catch (e) {
				alert("Error: " + e);
			}
		},

		handleClick(conversation) {
    		//console.log(`Clicked on ${conversation.user_or_group_name}`);
  		}
	},
	watch: {
    	conversation_id(newConversationId, oldConversationId) {
      		this.loadConversation()
		}
	},
}
</script>

<template>
	<div v-if="user != null" class="avolio-conversation-container">
		<!--div class="avolio-message-box avolio-friend-message" v-for="msg in messages" :key="msg.id">
			<p>{{ msg.msg }} {{ msg.from_user_id }}</p>
		</div-->
		{{ user.id }}
		<div :class="{'avolio-message-box avolio-friend-message':msg.from_user_id==1,
			'avolio-message-box avolio-my-message':msg.from_user_id!=1}"
		      v-for="msg in messages" :key="msg.id">
			<p>{{ msg.msg }} {{ msg.from_user_id }}</p>
		</div>
	</div>

	<ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>

<!--
		a @click.prevent="handleClick(c)" href="#"
		{{ msg.msg_id }}
		{{ msg.conversation_id }}
		{{ msg.from_user_id }}
		{{ msg.sent_timestamp }}
		{{ msg.msg }}
		<template v-if="msg.forwarded_from_msg_id != 0">[forwarded from {{ msg.forwarded_from_msg_id }}]</template>
		<template v-if="msg.seen == 0">
			✓
		</template>
		<template v-if="msg.seen == 1">
			✓✓
		</template>
		<template>
			{{ msg.reaction }}
		</template>
-->