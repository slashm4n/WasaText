<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
import SendMessageView from './SendMessageView.vue';
import ConversationsView from './ConversationsView.vue';
import ConversationView from './ConversationView.vue';
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
            sent_message_to_user_name : ''
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
                }
                else if (res.status == 201) {
                    console.log("Login done successfully. New user created. Session token ", this.session_token);
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
            this.sent_message_to_user_name = '';
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
                }
                this.username = this.new_name
                this.user.name = this.new_name
                this.new_name = ''
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
        },

        async onMessageSent(to_user_name) {
            this.sent_message_to_user_name = to_user_name;
        }
    }
}
</script>

<template>
    <div class="login-container">
        <input id="username" autocomplete="off" v-if="session_token == 0" style="position:relative; top:0.7em; left:1em; width: 8em;" type="text" v-model="username" placeholder="User name">
        <button v-if="session_token == 0" style="position:relative; top:0.7em; left:2em; background-color: white" @click="doLogin">Login</button>
    
        <img v-if="session_token != 0 && user.photo !=''" v-bind:src="user.photo" style="top:0em; left:1em; width: 3em; height: 3em;">
        <img v-if="session_token != 0 && user.photo ==''" src="../assets/profile.png" style="top:0em; left:1em; width: 3em; height: 3em;">
        
        <span v-if="session_token != 0" style="position:relative; top:-1.2em;">
            <span style="position:relative; left:1em; color:white; font-weight: bold;">{{ user.name }}</span>
            <button style="position:relative; left:2em; background-color: white" @click="doLogout">Logout</button>
            <input style="position:relative; left:6em; width: 8em;" type="text" v-model="new_name" placeholder="New name">
            <button style="position:relative; left:7em; background-color: white" @click="doSetMyUserName(new_name)">Apply</button>
            <span>
                <label for="photouploader" class="btn" style="position:relative; left:11em; background-color: white">Change profile photo</label>
                <input type="file" accept="image/*" hidden="true" id="photouploader" @change="doSetMyPhoto">
            </span>
        </span>
    </div>
    
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
    <SendMessageView :session_token="session_token" @messageSent="onMessageSent" ></SendMessageView>
    <ConversationsView :session_token="session_token" :user="user" :sent_message_to_user_name="sent_message_to_user_name" @selectedConversationChanged="onSelectedConversationChanged"></ConversationsView>
    <ConversationView :session_token="session_token" :user="user" :selected_conversation_id="selected_conversation_id"></ConversationView>
</template>
