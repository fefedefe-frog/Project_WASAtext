import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ChatsView from '../views/ChatsView.vue'
import ChatView from '../views/ChatView.vue'
import ProfileView from '../views/ProfileView.vue'
import UsersView from "../views/UsersView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		// schermata di login
		{path: '/session', component: LoginView, meta: {requiresAuth: false}},
		// schermata contentente il profilo di un utente
		{path: '/profile', component: ProfileView, meta: {requiresAuth: true}},
		// schermata "principale" dell'utente, mostra le chat di cui fa parte
		{path: '/chats', component: ChatsView, meta: {requiresAuth: true}},
		// schermata che mostra la lista di tutti gli utenti dell'app e ne consente la ricerca
		{path: '/users', component: UsersView, meta: {requiresAuth: true}},
		// schermata di una chat, contentente anche i messagi
		{path: '/chats/:chat_id', component: ChatView, meta: {requiresAuth: true}},
		// route per catturare tutti gli url non validi, riporterà alla schermata di login
		{path: '/:pathMatch(.*)*', redirect: '/session' }
	]
})

// Guard che si occupa di controllare se l'utente è autenticato
router.beforeEach((to, from, next) => {
	// La doppia ! in js indica che la variabile viene forzata ad uno stato booleano,
	// quindi se ha contenuto -> true, se non ha conteunto("") o è nulla -> false
	const isAuthenticated = !!localStorage.getItem("authToken");

	// Controllo che reindirizza tutti gli utenti non autenticati alla pagina di login
	if (to.meta.requiresAuth && !isAuthenticated) {
		next("/session");
	} else if (to.path === "/session" && isAuthenticated) {
		// Se l'utente invece è già autenticato e vuole accedere alla pagina per il login,
		// verrà reindirizzato direttamente alla pagina principale, ovvero in questo caso, la lista di chat
		next("/chats");
	} else {
		next(); // Altrimenti, continua normalmente verso la pagina selezionata
	}
});

export default router