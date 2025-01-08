import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ChatsView from '../views/ChatsView.vue'
import UserProfileView from '../views/UserProfileView.vue'
import ChatView from '../views/ChatView.vue'
import ChatInfoView from '../views/ChatInfoView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		// schermata di login
		{path: '/session', component: LoginView},
		// schermata contentente il profilo di un utente
		{path: '/profile', component: UserProfileView, meta: {requiresAuth: true}},
		// schermata contentente tutti gli utenti e una barra di ricerca
		{path: '/users', component: UserProfileView, meta: {requiresAuth: true}},
		// schermata "principale" dell'utente, mostra le chat di cui fa parte
		{path: '/chats', component: ChatsView, meta: {requiresAuth: true}},
		// schermata di una chat, contentente anche i messagi
		{path: '/chats/:chat_id', component: ChatView, meta: {requiresAuth: true}},
		// schermata che consente di modificare nome e foto di una chat, mostra i partecipanti e la possibiltà di aggiungerne
		{path: '/chats/:chat_id/info', component: ChatInfoView, meta: {requiresAuth: true}},
		// route per catturare tutti gli url non validi, riporterà alla schermata di login
		{ path: '/:pathMatch(.*)*', redirect: '/session' }
	]
})

// Guard che si occupa di controllare se l'utente è autenticato
router.beforeEach((to, from, next) => {
	const isAuthenticated = localStorage.getItem("authToken");

	if (to.meta.requiresAuth && !isAuthenticated) {
		// Se la rotta richiede autenticazione ma l'utente non è autenticato
		next("/session"); // Reindirizza alla pagina di login
	} else {
		next(); // Procedi normalmente
	}
});

export default router