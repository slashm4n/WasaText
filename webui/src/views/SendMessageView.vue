<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['messageSent', 'allUsersListUpdated'],
	props: ['session_token', 'user', 'need_update_all_users_list'],
	data: function() {
		return {
            errormsg: '',
            to_user_name_or_group_name: '',
            is_group: false,
            message: '',
            allusers: ''
		}
	},
	methods: {
		async doSendMessage() {
            try {
                if (this.to_user_name_or_group_name == '') {
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
                        "to_user_name_or_group_name": this.to_user_name_or_group_name,
                        "is_group": this.is_group,
                        "message": this.message
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('messageSent', this.to_user_name_or_group_name);
                this.message = '';

                this.errormsg = '';
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
        async doUpdateAllUsersList() {
            try {
                // Get the list of all users
                var res = await this.$axios({
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

                // Remove myself
                const index = res.data.map(u => u.user_id).indexOf(this.user.id);
                res.data.splice(index, 1);
                this.allusers = res.data;
                
                

                // Get the list of groups that belongs to the user
                /*res = await this.$axios({
                    method: 'get',
                    url: '/groups',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.allusers = this.allusers.concat(res.data);


*/
                this.$emit('allUsersListUpdated');
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
      		if (this.session_token == 0) {
                this.to_user_name_or_group_name = '';
                this.message = '';
            }
            else {
                this.errormsg = '';
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
    <div v-if="session_token != 0">
        <div class="send-message-container">
            <div style="position:relative; top: 0.7em; float: left;">
                <span class="label-flat">Send message to</span>
                <select style="position:relative; font-size: 1em; min-width: 8em;" v-model="to_user_name_or_group_name" >
                    <!--option @value=u v-for="u in allusers">{{ u.user_name + "  " + u.group_name }}</option-->
                    <option v-for="u in allusers">{{ u.user_name }}</option>
                </select>
                <input v-model="message" placeholder="Message"></input>
                <button @click="doSendMessage">Send</button>
                <span style="position:relative;left:1em;color:darkorange">Send to groups not yet implemented</span>
            </div>
        </div>
        <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
    </div>
</template>
