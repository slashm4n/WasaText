<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['messageSent', 'allUsersListUpdated', 'usersUpdated', 'myConversationsUpdated'],
	props: ['session_token', 'user', 'all_users', 'need_update_all_users_list'],
	data: function() {
		return {
            errormsg: '',
            to_user_name_or_group_name: '',
            message: ''
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
                
                const is_group = this.to_user_name_or_group_name.substring(this.to_user_name_or_group_name.length - 8) == " (group)";
                const name = is_group ? this.to_user_name_or_group_name.substring(0, this.to_user_name_or_group_name.length - 8) : this.to_user_name_or_group_name

                const res = await this.$axios({
                    method: 'post',
                    url: '/messages',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "to_user_name_or_group_name": name,
                        "is_group": is_group,
                        "message": this.message
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                //this.$emit('messageSent', this.to_user_name_or_group_name);
                this.$emit('myConversationsUpdated');
                this.message = '';

                this.errormsg = '';
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        // this code cannot be placed inside doSetMyPhoto because
        // await need to be on the top level of an async method
        async doSendPhotoAsync(img)
        {
            try {
                const is_group = this.to_user_name_or_group_name.substring(this.to_user_name_or_group_name.length - 8) == " (group)";
                const name = is_group ? this.to_user_name_or_group_name.substring(0, this.to_user_name_or_group_name.length - 8) : this.to_user_name_or_group_name

                const res = await this.$axios({
                    method: 'post',
                    url: '/messages',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "to_user_name_or_group_name": name,
                        "is_group": is_group,
                        "message": img
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                // this.$emit('messageSent', this.to_user_name_or_group_name);
                this.$emit('myConversationsUpdated');
                this.message = '';

                this.errormsg = '';
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        async doSendPhoto(e) {
            if (this.to_user_name_or_group_name == '') {
                this.errormsg = "Receiver user not set";
                return;
            }
            const img = e.target.files[0];

            if (img == null)
                return;
            const reader = new FileReader();
            reader.readAsDataURL(img);
            reader.onload = evn =>{
                this.doSendPhotoAsync(evn.target.result);
            }
        },

        async onImageUploaderClick() {
            this.$refs.photoSender.value = ''
        },

        async doAddToGroup() {
            try {
                if (this.to_user_name_or_group_name == '') {
                    this.errormsg = "User to add not set";
                    return;
                }
                
                const is_group = this.to_user_name_or_group_name.substring(this.to_user_name_or_group_name.length - 8) == " (group)";
                if (is_group) {
                    this.errormsg = "Can't add a group to another group";
                    return;
                }

                if (this.group_name == '') {
                    this.errormsg = "Group name not set";
                    return;
                }
                
                const res = await this.$axios({
                    method: 'put',
                    url: '/groups/0',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "user_name_to_add": this.to_user_name_or_group_name,
                        "group_name": this.group_name
                    }
                });

                if (res.status != 201) {
                    this.errormsg = "Unexpected response " + res.status;
                    return;
                }

                this.$emit('myConversationsUpdated');

                this.user_name_to_add = ''
                this.group_name = '';
                this.errormsg = "";
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
        all_users(newValue, oldValue) {
            // no more necessary. to remove
        }
	},
    beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
            localStorage.setItem('to_user_name_or_group_name',  JSON.stringify(this.to_user_name_or_group_name));
            localStorage.setItem('message',  JSON.stringify(this.message));
        });

        try {
            if (localStorage.getItem('to_user_name_or_group_name') != null) this.to_user_name_or_group_name = JSON.parse(localStorage.getItem('to_user_name_or_group_name'))
            if (localStorage.getItem('message') != null) this.message = JSON.parse(localStorage.getItem('message'))
        } catch {
        }
    }
}
</script>

<template>
    <div v-if="session_token != 0">
        <div class="send-message-container">
            <div style="position:relative; top: 0.7em; float: left;">
                <span class="label-flat">For user</span>
                <select style="z-index: 99; position:relative; height: 1.3em; width: 9.5em;" v-model="to_user_name_or_group_name" >
                    <option v-for="(u, index) in this.all_users" :key="index">{{ u.is_group ? u.group_name + " (group)" : u.user_name }}</option>
                </select>
                <input style="z-index: 100; position:absolute; top:0.2em; left: 5.3em; width: 8em;" v-model="to_user_name_or_group_name" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
                <input v-model="message" placeholder="Message">
                <button @click="doSendMessage">Send</button>
                <label for="photoSender" class="label-button">Send photo</label>
                <input type="file" accept="image/*" hidden="true" id="photoSender" ref="photoSender" @click="onImageUploaderClick" @change="doSendPhoto">
                <span class="label-flat">add to group</span>
                <input v-model="group_name" placeholder="Group name" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
                <button @click="doAddToGroup">Apply</button>
            </div>
        </div>
        <ErrorMsg :errormsg="errormsg" @error-dismissed="this.errormsg = '';"></ErrorMsg>
    </div>
</template>
