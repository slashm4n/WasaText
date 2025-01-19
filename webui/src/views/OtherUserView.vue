<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['messageSent', 'allUsersListUpdated', 'usersUpdated'],
	props: ['session_token', 'user', 'all_users', 'my_groups'],
	data: function() {
		return {
            errormsg: '',
            for_user_name_or_group_name: '',
            message: '',
            group_name: ''
		}
	},
	methods: {
		async doSendMessage() {
            try {
                if (this.for_user_name_or_group_name == '') {
                    this.errormsg = "Receiver user not set";
                    return;
                }

                if (this.message == '') {
                    this.errormsg = "Message not set";
                    return;
                }
                
                const is_group = this.for_user_name_or_group_name.substring(this.for_user_name_or_group_name.length - 4) == " (G)";
                const name = is_group ? this.for_user_name_or_group_name.substring(0, this.for_user_name_or_group_name.length - 4) : this.for_user_name_or_group_name

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

                this.$emit('messageSent', this.for_user_name_or_group_name);
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
                const is_group = this.for_user_name_or_group_name.substring(this.for_user_name_or_group_name.length - 4) == " (G)";
                const name = is_group ? this.for_user_name_or_group_name.substring(0, this.for_user_name_or_group_name.length - 4) : this.for_user_name_or_group_name

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

                this.$emit('messageSent', this.for_user_name_or_group_name);

                this.errormsg = '';
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        async doSendPhoto(e) {
            if (this.for_user_name_or_group_name == '') {
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
                if (this.for_user_name_or_group_name == '') {
                    this.errormsg = "User to add not set";
                    return;
                }
                
                const is_group = this.for_user_name_or_group_name.substring(this.for_user_name_or_group_name.length - 4) == " (G)";
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
                        "user_name_to_add": this.for_user_name_or_group_name,
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
      		if (newValue == 0) {
                this.errormsg = '';
                this.for_user_name_or_group_name = '';
                this.message = '';
                this.group_name = '';
            }
		}
	},
    beforeMount: function () {
        window.addEventListener('beforeunload', (e) => {
            localStorage.setItem('for_user_name_or_group_name',  JSON.stringify(this.for_user_name_or_group_name));
            localStorage.setItem('message',  JSON.stringify(this.message));
            localStorage.setItem('group_name',  JSON.stringify(this.group_name));
        });

        try {
            if (localStorage.getItem('for_user_name_or_group_name') != null) this.for_user_name_or_group_name = JSON.parse(localStorage.getItem('for_user_name_or_group_name'))
            if (localStorage.getItem('message') != null) this.message = JSON.parse(localStorage.getItem('message'))
            if (localStorage.getItem('group_name') != null) this.group_name = JSON.parse(localStorage.getItem('group_name'))
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
                <select id="forUserNameInput" style="position:relative; height: 1.3em; width: 10em;" v-model="for_user_name_or_group_name">
                    <option v-for="(u, index) in all_users" :key="index">{{ u.is_group ? u.group_name + " (G)" : u.user_name }}</option>
                </select>
                <input id="messageInput" v-model="message" placeholder="Message">
                <button @click="doSendMessage">Send</button>
                <label for="photoSender" class="label-button">Send photo</label>
                <p></p>
                <input id="photoSender" ref="photoSender" type="file" accept="image/*" hidden="true" @click="onImageUploaderClick" @change="doSendPhoto">
                <span class="label-flat">add to group</span>
                <select id="groupNameInput" style="position:relative; height: 1.3em; width: 9em;" v-model="group_name">
                    <option v-for="(u, index) in my_groups" :key="index">{{ u.group_name }}</option>
                </select>
                <span class="label-flat">or</span>
                <input id="newGroupInput" v-model="group_name" placeholder="New group" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
                <button @click="doAddToGroup">Apply</button>
            </div>
        </div>
        <ErrorMsg :errormsg="errormsg" @error-dismissed="this.errormsg = '';"></ErrorMsg>
    </div>
</template>
