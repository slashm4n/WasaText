<script setup>
//import { ref } from "vue";
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
	emits: ['messageModified', 'messageForwarded'],
    props: ['session_token', 'user', 'all_users', 'selected_message_id'],
    data: function() {
		return {
			errormsg: '',
			to_user_name_or_group_name: ''
		}
	},
	methods: {
		async doForwardMessage() {
            if (this.selected_message_id == 0) {
				return;
			}
			if (this.to_user_name_or_group_name == '') {
				this.errormsg = "Receiver user not set";
				return;
			}

			const is_group = this.to_user_name_or_group_name.substring(this.to_user_name_or_group_name.length - 8) == " (group)";
			const name = is_group ? this.to_user_name_or_group_name.substring(0, this.to_user_name_or_group_name.length - 8) : this.to_user_name_or_group_name

			try {
				const response = await this.$axios({
					method: 'post',
					url: '/messages/' + this.selected_message_id,
					headers: {
						'Authorization' : 'Bearer ' + this.session_token
					},
					data: {
                        "to_user_name_or_group_name": name,
                        "is_group": is_group
                    }
				});

                this.$emit('messageForwarded', name);
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
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
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
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
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},
		async onReactionClick(el) {
			if (this.selected_message_id == 0) {
				return;
			}

			if (el.target == null|| el.target.value == '') {
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
                        "reaction": el.target.value
                    }
				});

				if (response.status != 201) {
                    this.errormsg = "Unexpected response " + response.status;
                    return;
                }

				this.$emit('messageModified', this.to_user_name);
				this.errormsg = '';
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		}
	},
	watch: {
		session_token(newValue, oldValue) {
			if (newValue == 0) {
				this.errormsg = '';
				this.to_user_name_or_group_name = '';
			}
		},
        selected_message_id(newValue, oldValue) {
        }
	},
	beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
            localStorage.setItem('to_user_name_or_group_name',  JSON.stringify(this.to_user_name_or_group_name));
        });

        try {
            if (localStorage.getItem('to_user_name_or_group_name') != null) this.to_user_name_or_group_name = JSON.parse(localStorage.getItem('to_user_name_or_group_name'));
        } catch {
        }
    }
}
</script>

<template>
	<div v-if="session_token != 0">
		<div class="message-container">
			<div v-if="selected_message_id != 0" style="position: relative; top: 0.7em;">
				<span style="position: relative; left:1em">.</span>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F600;">&#x1F600;</button>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F602;">&#x1F602;</button>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F621;">&#x1F621;</button>
				<button class="emoticon-button" @click="onReactionClick" value="&#x1F44D;">&#x1F44D;</button>
				<button @click="onDeleteReaction">Remove</button>
				<button @click="onDeleteMessage">Delete message</button>
				<span class="label-flat">Forward to</span>
				<select id="userNameInput" style="position:relative; height: 1.3em; width: 9.5em;" v-model="to_user_name_or_group_name" >
					<option v-for="(u, index) in all_users" :key="index">{{ u.is_group ? u.group_name + " (group)" : u.user_name }}</option>
				</select>
				<button @click="doForwardMessage">Send</button>
			</div>
		</div>
		<ErrorMsg :errormsg="errormsg" @error-dismissed="this.errormsg = '';"></ErrorMsg>
	</div>
</template>
