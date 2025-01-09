<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    props: ['session_token'],
	data: function() {
		return {
            errormsg: '',
            to_user: '',
            message: ''
		}
	},
	methods: {
		async doSendMessage() {
            try {

                const res = await this.$axios({
                    method: 'post',
                    url: '/messages',
                    headers: {
                        'Authorization' : 'Bearer ' + this.session_token
                    },
                    data: {
                        "to_user": this.to_user,
                        "message": this.message
                    }
                });

                if (res.status != 200) {
                    console.log("Problem sending a new conversation");
                    alert("Problem sending a new conversation");
                }
            } catch (e) {
                    alert("Error: " + e);
            }
        },	
	}
}
</script>

<template>
	<div class="avolio-send-message-container" v-if="session_token != 0">
		<span style="position:relative; top: 1em;">
			<span style="position:relative; left: 1em;">Send message to: </span>
			<input style="position:relative; left: 3em; width: 15em;" v-model="to_user"></input>
			<input style="position:relative; left: 3em; width: 15em;" v-model="message"></input>
			<button style="position:relative; left:5em; background-color: white" @click="doSendMessage">Send</button>
		</span>
	</div>
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
