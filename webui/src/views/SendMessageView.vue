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
    	<span style="position:relative; top: 0.7em;">
            
            <label style="position:relative; left:1em; width: 10em; height: 4em; color: white; border: none;">All users:</label>
            <select name="languages" id="lang" style="position:relative; left:2em; font-size: 1em; min-width: 8em;">
                <option v-for="u in allusers">{{ u.user_name }}</option>
            </select>

            <span style="position:relative; left: 4em; color: white">Send message to: </span>
			<input style="position:relative; left: 5em; width: 8em;" v-model="to_user_name" placeholder="User name"></input>
			<input style="position:relative; left: 6em; width: 8em;" v-model="message" placeholder="Message"></input>
			<button style="position:relative; left: 7em; background-color: white" @click="doSendMessage">Send</button>
            <span style="position:relative; left: 8em; color:lightgreen">TO DO: manage group</span>
		</span>
	</div>
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
