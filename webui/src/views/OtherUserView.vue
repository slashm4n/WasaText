<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['messageSent', 'allUsersListUpdated', 'usersUpdated', 'reloginNeeded', 'groupsUpdated'],
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
		async onErrorDismissed() {
			this.errormsg = '';
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
            try {
                localStorage.setItem('for_user_name_or_group_name',  JSON.stringify(this.for_user_name_or_group_name));
                localStorage.setItem('message',  JSON.stringify(this.message));
                localStorage.setItem('group_name',  JSON.stringify(this.group_name));
            } catch {
                this.$emit('reloginNeeded');
            }
        });

        try {
            if (localStorage.getItem('for_user_name_or_group_name') != null) this.for_user_name_or_group_name = JSON.parse(localStorage.getItem('for_user_name_or_group_name'))
            if (localStorage.getItem('message') != null) this.message = JSON.parse(localStorage.getItem('message'))
            if (localStorage.getItem('group_name') != null) this.group_name = JSON.parse(localStorage.getItem('group_name'))
        } catch {
            this.$emit('reloginNeeded');
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
            </div>
        </div>
		<ErrorMsg :errormsg="errormsg" @errorDismissed="onErrorDismissed"></ErrorMsg>
    </div>
</template>
