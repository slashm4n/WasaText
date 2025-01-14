<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
import SendMessageView from './SendMessageView.vue';
import GroupManagementView from './GroupManagementView.vue';
import ConversationsView from './ConversationsView.vue';
import ConversationView from './ConversationView.vue';
import MessageView from './MessageView.vue'
</script>

<script>

export default {
    data: function() {
		return {
            errormsg: '',
            username: '',
            new_name: '',
            session_token: 0,
            user: null,
            selected_conversation_id : 0,
            selected_message_id : 0,
            need_update_conversations_list : false,
            need_update_conversation : false,
            need_update_all_users_list : false
		}
	},
	methods: {
        async doLogin() {
			try {
                if (this.username.length < 3) {
                    this.errormsg = "user name too short or too long";
                    return;
                }

                const res = await this.$axios({
                    method: 'post',
                    url: '/session',
                    data: {
                        name: this.username
                    }
                });
                
                const authHeader = res.headers['authorization'];
                
                if (authHeader && authHeader.startsWith('Bearer ')) {
                    this.session_token = authHeader.substring(7);
                    console.log("session token ", this.session_token);
                } else {
                    this.errormsg = "no valid authorization header found.";
                }
                
                if (res.status == 200) {
                    console.log("Login done successfully. User already existing. Session token ", this.session_token);

                    // also if the user exist, the all users list need to be updated soon after login
                    this.need_update_all_users_list = true;
                }
                else if (res.status == 201) {
                    console.log("Login done successfully. New user created. Session token ", this.session_token);

                    // need refresh the all users list!
                    this.need_update_all_users_list = true;
                }
                else {
                    this.errormsg = "Response code" + res.statusText + " diverso da quello atteso";
                    this.session_token = 0 // must reinitialize session token !
                    return;
                }

                // Retrieve the user data from the content
                this.user = res.data
                if (this.user == null) {
                    this.errormsg = "Not received user data at login";
                    this.session_token = 0 // must reinitialize session token !
                    return;
                }

                // Check the photo profile and if missing assign the default one
                if (this.user.photo == '') {
                    console.log("missing photo profile");
                }

			} catch (e) {
				this.errormsg = "Error: " + e;
			}
		},

		async doLogout() {
            // this.username = '';
            this.new_name = '';
            this.session_token = 0;
            this.user = null;
            this.selected_conversation_id = 0;
            this.selected_message_id = 0;
            this.need_update_conversation = false;
        },

		async doSetMyUserName(new_name) {
			try {
                const res = await this.$axios({
                    method: 'put',
                    url: '/user',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "new_name": new_name
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Problem updating the user name";
                    return;
                }
                this.username = this.new_name;
                this.user.name = this.new_name;
                this.new_name = '';

                // need refresh the all users list!
                this.need_update_all_users_list = true;
            } catch (e) {
                this.errormsg = "Error: " + e;
            }
        },

        // this code cannot be placed inside doSetMyPhoto because
        // await need to be on the top level of an async method
        async doUploadMyPhoto(img)
        {
            try {
                const res = await this.$axios({
                    method: 'put',
                    url: '/user/photo',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "photo": img
                    }
                });

                if (res.status != 200) {
                    this.errormsg = "Problem while updating the profile photo";
                }
            } catch (e) {
                this.errormsg = "Error: " + e;
            }
        },
		
        async doSetMyPhoto(e) {
            const img = e.target.files[0];
            const reader = new FileReader();
            reader.readAsDataURL(img);
            reader.onload = evn =>{
                this.doUploadMyPhoto(evn.target.result);
                this.user.photo = evn.target.result;
            };
        },
        
        async onSelectedConversationChanged(selected_conversation_id) {
            this.selected_conversation_id = selected_conversation_id;
            // must reset the selected message!
            this.selected_message_id = 0;
        },
        
        async onConversationsListUpdated() {
            this.need_update_conversations_list = false;
        },
        
        async onSelectedMessageChanged(selected_message_id) {
            this.selected_message_id = selected_message_id;
        },
        
        async onConversationUpdated() {
            this.need_update_conversation = false;
        },
        
        async onMessageSent(to_user_name) {
            if (this.selected_conversation_id == 0) {
                // TO DO: should force the refreash of the list of conversations and the selection of the conversation
                this.need_update_conversations_list = true
            }
            this.need_update_conversation = true;
        },

        async onUserAddedToGroup() {
            this.need_update_conversations_list = true
        },

        async onMessageModified(to_user_name) {
            this.need_update_conversation = true;
        },

        async onAllUsersListUpdated() {
            this.need_update_all_users_list = false
        }
    }
}
</script>

<template>
    <div class="login-container">
        <div style="position: relative; top: 0.7em; float: left;">
            <input v-if="session_token == 0" v-model="username" type="text" autocomplete="off" placeholder="User name">
            <button v-if="session_token == 0" @click="doLogin">Login</button>
            
            <img v-if="session_token != 0 && user.photo !=''" class="photo-box" style="top:-0.7em" v-bind:src="user.photo">
            <img v-if="session_token != 0 && user.photo ==''" class="photo-box" style="top:-0.7em" src="../assets/profile.png">
            
            <span v-if="session_token != 0" style="position: relative; top:-1.8em">
                <span class="label-flat" style="font-weight: bold;">{{ user.name }}</span>
                <button @click="doLogout">Logout</button>
                <input type="text" v-model="new_name" placeholder="New name">
                <button @click="doSetMyUserName(new_name)">Apply</button>
                <label for="photoUploader" class="label-button">Set photo</label>
                <input type="file" accept="image/*" hidden="true" id="photoUploader" @change="doSetMyPhoto">
            </span>
        </div>
    </div>
    
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
    <SendMessageView :session_token="session_token" :need_update_all_users_list="need_update_all_users_list" @allUsersListUpdated="onAllUsersListUpdated" @messageSent="onMessageSent"></SendMessageView>
    <GroupManagementView :session_token="session_token" @userAddedToGroup="onUserAddedToGroup"></GroupManagementView>
    <ConversationsView :session_token="session_token" :user="user" :need_update_conversations_list="need_update_conversations_list" @selectedConversationChanged="onSelectedConversationChanged"></ConversationsView @conversationsListUpdated="onConversationsListUpdated">
    <ConversationView :session_token="session_token" :user="user" :selected_conversation_id="selected_conversation_id" :need_update_conversation="need_update_conversation" @selectedMessageChanged="onSelectedMessageChanged" @conversationUpdated="onConversationUpdated"></ConversationView>
    <MessageView :session_token="session_token" :selected_message_id="selected_message_id" @messageModified="onMessageModified"></MessageView>
</template>
