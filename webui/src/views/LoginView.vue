<script>
export default {
  data() {
    return {
      username: '',
      usrId: '',
      errormsg: null, // Messaggio di errore
      loading: false, // Stato di caricamento
      welcomeMsg: false,
    };
  },
  computed: {
    isUsernameValid() {
      const username = this.username;
      return (username.length >= 3 && username.length <= 16 && ((/^\S.*\S$/).test(username)));
    },
  },
  methods: {
    async doLogin() {
      // Imposta lo stato di caricamento
      this.loading = true;
      this.error = null;

      // Controllo se il nome inserito rispetta le regex prestabilite
      const regex= /^\S.*\S$/;
      if (!regex.test(this.username)) {
      }else

      try {
        // Esegui la richiesta POST e aspetta la risposta con `await`
        const response = await this.$axios.post('http://localhost:3000/session', {
          userName: (this.username).toLowerCase(),
        });

        // Estraggo il token dall'header
        const token = response.headers["authorization"];
        const usrId = response.data.usrId;
        this.usrId = usrId;

        if (token && usrId) {

          // Salvo il token e usrId
          localStorage.setItem('authToken', token);
          localStorage.setItem('usrId', usrId);


          // Ritarda la redirezione per consentire il caricamento
          setTimeout(() => {
            this.$router.push('/chats'); // Redirigi alla schermata principale delle chat
          }, 1000); // Attendi 2 secondi
        }
      } catch (error) {
        // Gestisci errori di rete o di altro tipo
        this.errormsg = error.toString();
      }
    },
  }
};
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Login</h1>
    </div>
    <div class="login-container">
      <!-- Form di login -->
      <form class="login-form" @submit.prevent="doLogin">
        <label for="username">Nome utente:  <span v-if="!isUsernameValid" class="error">Nome non valido</span></label>
        <input id="username" v-model="username" type="text" placeholder="Inserisci il nome utente" required>
        <button type="submit" :disabled="!isUsernameValid || loading" :class="{ disabled: !isUsernameValid || loading }">Login</button>
      </form>

      <!-- Mostra messaggi di errore -->
      <ErrorMsg v-if="errormsg" :msg="errormsg" />

      <!-- Spinner di caricamento -->
      <LoadingSpinner v-if="loading" loading="{{ loading }}" :loading-text="'Benvenuto/a '+ username +'! Caricamento chat in corso'" /><LoadingSpinner />
    </div>
  </div>
</template>



<style scoped>
form {
  display: flex;
  flex-direction: column;
  width: 300px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

input {
  margin: 8px 0;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #ccc;
}

button {
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

button:disabled,
button.disabled {
  background-color: #a5a5a5; /* Colore più scuro */
  cursor: not-allowed; /* Cambia il cursore */
  color: #f0f0f0; /* Testo leggermente più chiaro */
}

.error {
  color: red;
  margin-top: 10px;
}
</style>