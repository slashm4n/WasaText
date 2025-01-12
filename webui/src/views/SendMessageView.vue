<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['messageSent', 'allUsersListUpdated'],
	props: ['session_token', 'need_update_all_users_list'],
	data: function() {
		return {
            errormsg: '',
            to_user_name: '',
            message: '',
            allusers: ''
		}
	},
	methods: {
		async doSendMessage() {
            try {
                if (this.to_user_name == '') {
                    this.errormsg = "Receiver user not set";
                    return;
                }

                if (this.message == '') {
                    this.errormsg = "Message not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'post',
                    url: '/messages',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "to_user": this.to_user_name,
                        "message": this.message
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('messageSent', this.to_user_name);
                this.message = '';
            } catch (e) {
                this.errormsg = "Error: " + e;
            }
        },
        async doUpdateAllUsersList() {
            try {
                // Reload the list of all users
                const res = await this.$axios({
                    method: 'get',
                    url: '/users',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.allusers = res.data;
                this.$emit('allUsersListUpdated');
            } catch (e) {
                this.errormsg = "Error: " + e;
            }
        }
	},
    watch: {
    	session_token(newValue, oldValue) {
      		if (this.session_token == 0) {
                this.to_user_name = '';
                this.message = '';
            }
		},
        need_update_all_users_list(newValue, oldValue) {
            if (this.need_update_all_users_list)
                this.doUpdateAllUsersList();
        }
	}
}
</script>

<template>
	<div class="send-message-container" v-if="session_token != 0">
    	<div style="position:relative; top: 0.7em; float: left;">
            
            <span class="label-flat">All users:</span>
            <select name="languages" id="lang" style="position:relative; font-size: 1em; min-width: 8em;">
                <option v-for="u in allusers">{{ u.user_name }}</option>
            </select>

            <span class="label-flat">Send message to:</span>
			<input v-model="to_user_name" placeholder="User name"></input>
			<input v-model="message" placeholder="Message"></input>
			<button @click="doSendMessage">Send</button>
		</div>
	</div>
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
