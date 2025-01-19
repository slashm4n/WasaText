<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['selectedConversationChanged', 'conversationsListUpdated'],
	props: ['session_token', 'user', 'need_update_conversations_list'],
	data() {
		return {
			errormsg: '',
			conversations: [],
			selected_conversation_id: 0
		}
	},
	methods: {
    	async doUpdateConversationsList() {
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
				
				this.$emit('conversationsListUpdated')

			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
			this.selected_conversation_id = 0;
			this.errormsg = "";
		},
		async onConversationClick(conv) {
			this.selected_conversation_id = conv.conversation_id;
			this.$emit('selectedConversationChanged', this.selected_conversation_id);
		},
		async onConversationFocusOut() {
			this.selected_conversation_id = 0;
			// this.$emit('selectedConversationChanged', this.selected_conversation_id);
		}
  	},
  	watch: {
    	session_token(newValue, oldValue) {
			if (newValue == 0) {
				this.errormsg = '',
				this.conversations = [],
				this.selected_conversation_id = 0
			}
			else
				this.doUpdateConversationsList();
		},
    	need_update_conversations_list(newValue, oldValue) {
			if (newValue)
	      		this.doUpdateConversationsList();
			// TO DO: riselezionare l'ultima conversazione
		}
	},
	beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
            localStorage.setItem('conversations',  JSON.stringify(this.conversations));
            localStorage.setItem('selected_conversation_id', JSON.stringify(this.selected_conversation_id));
        });

        try {
            if (localStorage.getItem('conversations') != null) this.conversations = JSON.parse(localStorage.getItem('conversations'));
            if (localStorage.getItem('selected_conversation_id') != null) this.selected_conversation_id = JSON.parse(localStorage.getItem('selected_conversation_id'));
        } catch {
        }
    }
}
</script>

<template>
	<div v-if="session_token != 0">
		<div class="conversations-container" v-if="session_token != 0">
			<div class="conv-box" v-for="(c, index) in conversations" :key="c.id" :tabindex="index" @click="onConversationClick(c)" @focusout = "onConversationFocusOut">
				<img v-if="c.user_or_group_photo != ''" class="photo-box" v-bind:src="c.user_or_group_photo">
				<img v-if="c.user_or_group_photo == ''" class="photo-box" src="../assets/profile.png">
				<span class="name-box">{{ c.user_or_group_name }}{{ c.is_group ? ' (G)' : '' }}</span>
				<span class="snippet-box">{{ c.last_timestamp }} {{ c.last_msg.substring(0, 11) == "data:image/" ? "&#x1F4F7;" : c.last_msg.substring(0, 8) }}</span>
			</div>
		</div>
		<ErrorMsg :errormsg="errormsg" @error-dismissed="this.errormsg = '';"></ErrorMsg>
	</div>
</template>
