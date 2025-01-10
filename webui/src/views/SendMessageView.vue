<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    emits: ['messageSent'],
	props: ['session_token'],
	data: function() {
		return {
            errormsg: '',
            to_user_name: '',
            message: ''
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
	},
    watch: {
    	session_token(newToken, oldToken) {
      		if (this.session_token == 0) {
                this.to_user_name = '';
                this.message = '';
            }
		}
	}
}
</script>

<template>
	<div class="send-message-container" v-if="session_token != 0">
    	<span style="position:relative; top: 0.7em;">
            <span style="position:relative; left: 1em; color: white">Send message to: </span>
			<input style="position:relative; left: 2em; width: 8em;" v-model="to_user_name" placeholder="User name"></input>
			<input style="position:relative; left: 4em; width: 8em;" v-model="message" placeholder="Message"></input>
			<button style="position:relative; left:5em; background-color: white" @click="doSendMessage">Send</button>
		</span>
	</div>
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
