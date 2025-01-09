<script setup>
import { watch } from 'vue';
const props = defineProps({
	errormsg: {
		type: String,
		required: true,
	},
});

watch(() => props.errormsg, (newValue) => {
	console.log(newValue);
});
</script>

<script>
export default {
	emits: ["errorWindowClosed"],
}
</script>

<template>
	<div v-if="errormsg.length" class="error-container">
		<div class="error-header">
		<h2>Error</h2>
		</div>
		<p>{{errormsg}}</p>
		<button class="error-button" @click="$emit('errorWindowClosed', '')">OK</button>
	</div>
</template>
