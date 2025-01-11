<script setup>
import ErrorMsg from '../components/ErrorMsg.vue';
</script>

<script>
export default {
    props: ['session_token', 'user', 'selected_message_id'],
    data: function() {
		return {
			errormsg: '',
		}
	},
	methods: {
	},
	watch: {
        selected_message_id(newValue, oldValue) {

        }
	},
}
</script>

<template>
    <div class="message-container" v-if="session_token != 0" style="color:lightgreen">
        <div style="position: relative; left: 1em">selected message: {{ selected_message_id }}
        </div>
        <div style="position: relative; left: 1em">TO DO: manage comment/uncomment/delete</div>
    </div>
    <ErrorMsg :errormsg="errormsg" @errorWindowClosed="this.errormsg = '';"></ErrorMsg>
</template>
