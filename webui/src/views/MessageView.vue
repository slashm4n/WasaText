<script setup>
//import { ref } from "vue";
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['messageModified'],
    props: ['session_token', 'user', 'selected_message_id'],
    data: function() {
		return {
			errormsg: '',
		}
	},
	methods: {
        async onDeleteMessage() {
            if (this.selected_message_id == 0) {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'delete',
					url: '/messages/' + this.selected_message_id,
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				this.$emit('messageModified', this.to_user_name);
			} catch (e) {
				this.errormsg = "Error: " + e;
			}
        },
		async onApplyReaction() {
			if (this.selected_message_id == 0) {
				return;
			}

			if (this.$refs.comment == null|| this.$refs.comment.value == '') {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'post',
					url: '/messages/:msg_id/comments',
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					},
                    data: {
                        "msg_id": this.selected_message_id,
                        "reaction": this.$refs.comment.value
                    }
				});

				if (response.status != 201) {
                    this.errormsg = "Unexpected response " + response.status;
                    return;
                }

				this.$emit('messageModified', this.to_user_name);
			} catch (e) {
				this.errormsg = "Error: " + e;
			}
		},
		async onDeleteReaction() {
            if (this.selected_message_id == 0) {
				return;
			}
			try {
				const response = await this.$axios({
					method: 'delete',
					url: '/messages/:' + this.selected_message_id + '/comments',
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					}
				});

				this.$emit('messageModified', this.to_user_name);
			} catch (e) {
				this.errormsg = "Error: " + e;
			}
		}
	},
	watch: {
        selected_message_id(newValue, oldValue) {
        }
	},
}
</script>

<template>
    <div v-if="session_token != 0" class="message-container">
        <div v-if="selected_message_id != 0" style="position: relative; top: 0.7em">
            <input ref="comment" type="text" autocomplete="off" placeholder="Reaction">
			<button @click="onApplyReaction">Apply</button>
			<button @click="onDeleteReaction">Remove</button>
			<button @click="onDeleteMessage">Delete message</button>
        </div>
    </div>
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
