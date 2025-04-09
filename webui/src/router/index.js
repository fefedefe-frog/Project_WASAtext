import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainView from "../views/MainView.vue";


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		// schermata di login
		{path: '/session', component: LoginView, meta: {requiresAuth: false}},
		// schermata "principale" dell'utente, mostra le chat di cui fa parte
		{path: '/home', component: MainView, meta: {requiresAuth: true}},
		// route per catturare tutti gli url non validi, riporterà alla schermata di login
		{path: '/:pathMatch(.*)*', redirect: '/session' }
	]
})

// Controllo che si occupa di verificare se l'utente è autenticato prima di reindirizzarlo sulla pagina principale
router.beforeEach((to, from, next) => {
	// La doppia ! in js indica che la variabile viene forzata ad uno stato booleano,
	// quindi se ha contenuto -> true, se non ha conteunto("") o è nulla -> false
	const isAuthenticated = !!sessionStorage.getItem("authToken");

	// Controllo che reindirizza tutti gli utenti non autenticati alla pagina di login
	if (to.meta.requiresAuth && !isAuthenticated) {
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