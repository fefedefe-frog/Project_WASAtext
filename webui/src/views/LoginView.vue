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
        const response = await this.$axios.post('/session', {
          userName: (this.username).toLowerCase(),
        });

        // Estraggo il token dall'header
        const token = response.headers["authorization"];
        const usrId = response.data.usrId;
        this.usrId = usrId;

        if (token && usrId) {

          // Salvo il token e usrId
          sessionStorage.setItem('authToken', token);
          sessionStorage.setItem('usrId', usrId);


          // Ritarda la redirezione per consentire il caricamento
          setTimeout(() => {
            this.$router.push('/home'); // Redirigi alla schermata principale delle chat
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
  <div class="container">
    <div class="login-container" v-if="!loading">
      <!-- Form di login -->
      <form class="login-form" @submit.prevent="doLogin">
        <h3 class="h3">Login</h3>
        <label for="username">Nome utente:  <span v-if="!isUsernameValid && username" class="error">Nome non valido</span></label>
        <input id="username" v-model="username" type="text" placeholder="Inserisci il nome utente" required>
        <button type="submit" :disabled="!isUsernameValid || loading" :class="{ disabled: !isUsernameValid || loading }">Login</button>
      </form>
    </div>
    <!-- Mostra messaggi di errore -->
    <ErrorMsg v-if="errormsg" :msg="errormsg" />

    <!-- Spinner di caricamento -->
    <LoadingSpinner :loading="loading" :loading-text="'Benvenuto/a '+ username +'! Caricamento in corso'" />
  </div>
</template>



<style scoped>
.container {
  display: flex;
  align-items: center;
  justify-content: center;
}

form {
  display: flex;
  flex-direction: column;

  width: 300px;
  margin: 0 auto;
  padding: 20px;

  border-radius: 8px;
  background: #ccc;
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