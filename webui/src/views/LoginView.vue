<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
import NewConversationView from './NewConversationView.vue';
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
            session_token: 0,
            username: '',
            new_name: '',
            user: null,
            all_users: null,
            all_users_and_my_groups: null,
            my_groups: null,
            selected_conversation : null,
            selected_conversation_id : 0,
            selected_message_id : 0,
            need_update_conversations_list : false,
            need_update_conversation : false,
            // need_update_all_users_list : false,
            need_update_groups_list : false
		}
	},
	methods: {
        async doLogin() {
            try {
                 if (this.username.length < 3 || this.username.length > 16) {
                    this.errormsg = "user name too short or too long (min 3, max 16)";
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
                    // console.log("session token ", this.session_token);
                } else {
                    this.errormsg = "no valid authorization header found.";
                }

                if (res.status == 200) {
                    console.log("Login done successfully. User already existing. Session token ", this.session_token);

                    // also if the user exist, the all users list need to be updated soon after login
                    // this.need_update_all_users_list = true;
                }
                else if (res.status == 201) {
                    console.log("Login done successfully. New user created. Session token ", this.session_token);

                    // need refresh the all users list!
                    // this.need_update_all_users_list = true;
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
                    // console.log("missing photo profile");
                }

                // Soon after the login Update the users list and my groups list
                this.doUpdateAllUsersAndMyGroupsList();
                this.doUpdateMyGroupsList();

                // Exit
                this.errormsg = "";
			} catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
			}
		},

		async doLogout() {
            this.errormsg = '';
            this.session_token = 0;
            // this.username = ''; // we leave the username for ease
            this.new_name = '';
            this.user = null;
            this.all_users_and_my_groups = null;
            this.my_groups = null;
            this.selected_conversation = null;
            this.selected_message_id = 0;
            this.need_update_conversations_list = false;
            this.need_update_conversation = false;
            // this.need_update_all_users_list = false;
            this.need_update_groups_list = false;
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
                this.errormsg = "";

                // need refresh the all users list!
                // this.need_update_all_users_list = true;
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        // this code cannot be placed inside doSetMyPhoto because
        // await need to be on the top level of an async method
        async doSetMyPhotoAsync(img)
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
                
                this.errormsg = "";
            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
		
        async doSetMyPhoto(e) {
            const img = e.target.files[0];
            if (img == null)
                return;
            const reader = new FileReader();
            reader.readAsDataURL(img);
            reader.onload = evn =>{
                this.doSetMyPhotoAsync(evn.target.result);
                this.user.photo = evn.target.result;
            };
        },

        async doUpdateAllUsersAndMyGroupsList() {
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

                // Remove myself!
                const index = res.data.map(u => u.user_id).indexOf(this.user.id);
                res.data.splice(index, 1);
                this.all_users_and_my_groups = res.data;
                
                // We have the first result with all users
                this.all_users = res.data
                
                // Add the column 'is_group' with value false
                res.data.forEach(function(e) {
                    e["is_group"] = false;
                });

                // Get the list of groups that belongs to the user
                res = await this.$axios({
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

                // Add the column 'is_group' with value true
                if (res.data != null) {
                    res.data.forEach(function(e) {
                        e["is_group"] = true;
                    });
                    this.all_users_and_my_groups = this.all_users_and_my_groups.concat(res.data);
                }

                // Merge id and name
                if (this.all_users_and_my_groups != null) {
                    this.all_users_and_my_groups.forEach(function(e) {
                        if (e.is_group) {
                            e["user_id_or_group_id"] = e["group_id"];
                            e["user_name_or_group_name"] = e["group_name"];
                        }
                        else {
                            e["user_id_or_group_id"] = e["user_id"];
                            e["user_name_or_group_name"] = e["user_name"];
                        }
                    });
                }

            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },
        
        async doUpdateMyGroupsList() {
            try {
                // Get the list of groups that belongs to the user
                const res = await this.$axios({
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

                this.my_groups = res.data;

            } catch (e) {
                if (e.response != null && e.response.data != "")
                    this.errormsg = "Error: " + e.response.data;
                else
                    this.errormsg = "Error: " + e;
            }
        },

        


        async onNeedUpdateAllUsersList() {
            this.doUpdateAllUsersAndMyGroupsList();
            this.need_update_conversations_list = true;
        },

        async onNeedUpdateMyGroupsList() {
            this.doUpdateMyGroupsList();
            this.need_update_conversations_list = true;
            this.doUpdateAllUsersAndMyGroupsList();
            // this.need_update_all_users_list = true;
        },

        async onSelectedConversationChanged(selected_conversation) {
            this.selected_conversation = selected_conversation;
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
            this.need_update_conversations_list = true;
            this.need_update_conversation = true;
        },

        async onMessageModified(to_user_name) {
            this.need_update_conversations_list = true;
            this.need_update_conversation = true;
        },

        async onAllUsersListUpdated() {
            this.doUpdateAllUsersAndMyGroupsList();
        },

        async onPhotoUploaderClick() {
            this.$refs.photoUploader.value = '';
        },

        async onMyConversationsUpdated() {
            this.need_update_conversations_list = true;
        },

        async onReloginNeeded() {
            // non funziona cambiare evento?
            // window.addEventListener('load', (e) => {
                // this.errormsg = "some values get lost during refresh, need to relogin";
            // })
        },

        async onErrorDismissed() {
            this.errormsg = '';
        }
    },
    beforeMount: function () {
        // uncomment ONLY for testing a new installation!
        /*
        localStorage.removeItem('session_token');
        localStorage.removeItem('username');
        localStorage.removeItem('new_name');
        localStorage.removeItem('user');
        localStorage.removeItem('all_users_and_my_groups');
        localStorage.removeItem('my_groups');
        localStorage.removeItem('selected_conversation');
        localStorage.removeItem('selected_message_id');
        localStorage.removeItem('need_update_conversations_list');
        localStorage.removeItem('need_update_conversation');
        */

        // save the state
        window.addEventListener('beforeunload', (e) => {
            localStorage.setItem('session_token',  JSON.stringify(this.session_token));
            localStorage.setItem('username', JSON.stringify(this.username));
            localStorage.setItem('new_name', JSON.stringify(this.new_name));
            localStorage.setItem('user', JSON.stringify(this.user));
            localStorage.setItem('all_users_and_my_groups', JSON.stringify(this.all_users_and_my_groups));
            localStorage.setItem('my_groups', JSON.stringify(this.my_groups));
            localStorage.setItem('selected_conversation_id', JSON.stringify(this.selected_conversation_id));
            localStorage.setItem('selected_message_id', JSON.stringify(this.selected_message_id));
            localStorage.setItem('need_update_conversations_list', JSON.stringify(this.need_update_conversations_list));
            localStorage.setItem('need_update_conversation', JSON.stringify(this.need_update_conversation));
            // localStorage.setItem('need_update_all_users_list', JSON.stringify(this.need_update_all_users_list));
        });

        // retrieve the state
        try {
            if (localStorage.getItem('session_token') != null) this.session_token = JSON.parse(localStorage.getItem('session_token'));
            if (localStorage.getItem('username') != null) this.username = JSON.parse(localStorage.getItem('username'));
            if (localStorage.getItem('new_name') != null) this.new_name = JSON.parse(localStorage.getItem('new_name'));
            if (localStorage.getItem('user') != null) this.user = JSON.parse(localStorage.getItem('user'));
            if (localStorage.getItem('all_users_and_my_groups') != null) this.all_users_and_my_groups = JSON.parse(localStorage.getItem('all_users_and_my_groups'));
            if (localStorage.getItem('my_groups') != null) this.my_groups = JSON.parse(localStorage.getItem('my_groups'));
            if (localStorage.getItem('selected_conversation_id') != null) this.selected_conversation_id = JSON.parse(localStorage.getItem('selected_conversation_id'));
            if (localStorage.getItem('need_update_conversations_list') != null) this.need_update_conversations_list = JSON.parse(localStorage.getItem('need_update_conversations_list'));
            if (localStorage.getItem('need_update_conversation') != null) this.need_update_conversation = JSON.parse(localStorage.getItem('need_update_conversation'));
            // if (localStorage.getItem('need_update_all_users_list') != null) this.need_update_all_users_list = JSON.parse(localStorage.getItem('need_update_all_users_list'));
        } catch {
            this.session_token = 0
        }
    }
}
</script>

<template>
    <div class="login-container">
        <div style="position: relative; top: 0.7em; float: left;">
            <input id="userNameInput" v-if="session_token == 0" v-model="username" type="text" placeholder="User name" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
            <button v-if="session_token == 0" @click="doLogin">Login</button>
            <img v-if="session_token != 0 && user != null && user.photo !=''" class="photo-box" style="top:-0.7em" v-bind:src="user.photo">
            <img v-if="session_token != 0 && user != null && user.photo ==''" class="photo-box" style="top:-0.7em" src="../assets/profile.png">
            <span v-if="session_token != 0 && user != null" style="position: relative; top:-1.8em">
                <span class="label-flat" style="font-weight: bold;">{{ user.name }}</span>
                <button @click="doLogout">Logout</button>
                <input id="newNameInput" type="text" v-model="new_name" placeholder="New name" autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="false">
                <button @click="doSetMyUserName(new_name)">Apply</button>
                <label for="photoUploader" class="label-button">Set photo</label>
                <input id="photoUploader" ref="photoUploader" type="file" accept="image/*" hidden="true" @click="onPhotoUploaderClick" @change="doSetMyPhoto">
            </span>
        </div>
    </div>
    
    <ErrorMsg :errormsg="errormsg" @errorDismissed="onErrorDismissed"></ErrorMsg>
    <GroupManagementView :session_token="session_token" :all_users="all_users"  :my_groups="my_groups" @groupsUpdated="onNeedUpdateMyGroupsList" @reloginNeeded="onReloginNeeded"></GroupManagementView>
    <NewConversationView :session_token="session_token" :user="user" :all_users="all_users" @messageSent="onMessageSent"  ></NewConversationView>
    <ConversationsView :session_token="session_token" :user="user" :need_update_conversations_list="need_update_conversations_list" @selectedConversationChanged="onSelectedConversationChanged" @conversationsListUpdated="onConversationsListUpdated" @reloginNeeded="onReloginNeeded"></ConversationsView>
    <ConversationView :session_token="session_token" :user="user" :selected_conversation="selected_conversation" :need_update_conversation="need_update_conversation" @selectedMessageChanged="onSelectedMessageChanged" @conversationUpdated="onConversationUpdated" @reloginNeeded="onReloginNeeded"></ConversationView>
    <MessageView :session_token="session_token" :user="user" :all_users_and_my_groups="all_users_and_my_groups" :selected_message_id="selected_message_id" :selected_conversation="selected_conversation" @messageSent="onMessageSent" @messageModified="onMessageModified" @messageForwarded="onMessageSent" @reloginNeeded="onReloginNeeded"></MessageView>
</template>
