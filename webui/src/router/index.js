import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainView from "../views/MainView.vue";

import UsersView from "../views/UsersView.vue";
import UserInfoView from "../views/UserInfoView.vue";

import ChatsView from "../views/ChatsView.vue";
import ChatView from "../views/ChatView.vue";
import NewChatView from "../views/NewChatView.vue";


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		// schermata di login
		{path: '/session', component: LoginView, name: "login",},
		// schermata "principale"
		{path: '/home', component: MainView, name: "home"},
		// schermata dove l'utente può vedere tutti gli altri utenti, e selezionandone uno vedere le sue info e mandargli un messaggio
		{path: '/users', component: UsersView, name: "users", children:[{path: ':usr_id', component: UserInfoView, props: true}]},
		// schermata delle chat, con schermata figlia che carica la singloa chat
		{path: '/chats', component: ChatsView, name: "chat", children:[{path: ':chat_id', component: ChatView, props: true}]},
		// schermata per creare nuove chat
		{path: '/newChat', component: NewChatView, name: "newChat"},
		// route per catturare tutti gli url non validi, riporterà alla schermata di login, o alla home se già loggato
		{path: '/:pathMatch(.*)*', redirect: '/session' }
	]
})

// Controllo che si occupa di verificare se l'utente è autenticato prima di reindirizzarlo sulla pagina principale
router.beforeEach((to, from, next) => {
	// La doppia ! in js indica che la variabile viene forzata ad uno stato booleano,
	// quindi se ha contenuto -> true, se non ha conteunto("") o è nulla -> false
	const isAuthenticated = !!sessionStorage.getItem("authToken");

	// Controllo che reindirizza tutti gli utenti non autenticati alla pagina di login
	if (to.name !== "login" && !isAuthenticated) {
		next("/session");
	} else if (to.path === "/session" && isAuthenticated) {
		// Se l'utente invece è già autenticato e vuole accedere alla pagina per il login,
		// verrà reindirizzato direttamente alla pagina principale, ovvero in questo caso, la lista di chat
		next("/home");
	} else {
		next(); // Altrimenti, continua normalmente verso la pagina selezionata
	}
});

export default router