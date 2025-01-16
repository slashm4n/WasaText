<script setup>
</script>

<script>
export default {
	emits: ["errorWindowClosed"],
	props: ['errormsg'],
}
</script>

<template>
	<div v-if="errormsg.length" class="alert-container alert-danger" role="alert">
		<span style="position: relative; top:0.7em">{{ errormsg }}</span>

		<!--div class="error-header">
		<h2>Error</h2>
		</div>
		<p>{{ errormsg }}</p>
		<button class="error-button" @click="$emit('errorWindowClosed', '')">OK</button-->
	</div>
</template>
