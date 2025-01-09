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
  methods: {
    async doLogin() {
      // Imposta lo stato di caricamento
      this.loading = true;
      this.error = null;

      try {
        // Esegui la richiesta POST e aspetta la risposta con `await`
        const response = await this.$axios.post('http://localhost:3000/session', {
          userName: this.username
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
          }, 2000); // Attendi 2 secondi
        } else {
          throw new Error('Token o usrId non ricevuti dal server.');
        }
      } catch (error) {
        // Gestisci errori di rete o di altro tipo
        this.errormsg = error.toString();
      }
    },
  },
  created(){
    localStorage.setItem('authToken', "");
    localStorage.setItem('usrId', "");
  },
  mounted() {
    localStorage.setItem('authToken', "");
    localStorage.setItem('usrId', "");
  }
};
</script>

<template>
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Login</h1>
  </div>
  <div class="login-container">
    <!-- Form di login -->
    <form @submit.prevent="doLogin" class="login-form">
      <label for="username">Nome utente:</label>
      <input id="username" v-model="username" type="text" placeholder="Inserisci il nome utente" required/>
      <button type="submit" :disabled="!username || loading" :class="{ disabled: !username || loading }">Login</button>
    </form>

    <!-- Mostra messaggi di errore -->
    <ErrorMsg v-if="errormsg" :msg="errormsg" />

    <!-- Spinner di caricamento -->
    <LoadingSpinner v-if="loading" loading="{{ loading }}" :loadingText="'Benvenuto/a '+ username +'! Caricamento chat in corso'"/><LoadingSpinner />
  </div>
</template>



<style scoped>
/* Stili per il form di login */
.login {
  width: 300px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

form {
  display: flex;
  flex-direction: column;
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

.alert {
  display: block !important;
}

.welcome-message {
  text-align: center;
  margin-top: 20px;
}
</style>