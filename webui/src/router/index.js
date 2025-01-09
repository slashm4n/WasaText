import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ConversationsView from '../views/ConversationsView.vue'
import ConversationView from '../views/ConversationView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/conversations', component: ConversationsView},
		{path: '/conversations/:conversation_id', name: 'conversation', component: ConversationView, props: true},
	]
})

export default router
