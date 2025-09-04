<script>
export default {
  data() {
    return {
      username: "",
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
      if (regex.test(this.username)) {
        try {
          // Esegui la richiesta POST e aspetta la risposta con `await`
          const response = await this.$axios.post('/session', {
            userName: (this.username).toLowerCase(),
          });

          // Estraggo il token dall'header
          const token = response.headers['authorization'];
          if (response.data && response.data['user']){

            this.username= response.data['user']['userName'];

            // Salvo il token e usrId
            sessionStorage.setItem('authToken', token);
            sessionStorage.setItem('usrId', response.data['user']['usrId']);

            // Ritarda la redirezione per consentire il caricamento
            setTimeout(() => {
              this.$router.push('/home'); // Redirigi alla schermata principale delle chat
            }, 1000); // Attendi 2 secondi
          }
        } catch(e) {
          let error_string= ""
          if (e.response.status === 400 ||  //Bad request
              e.response.status === 500){   //Internal server error
            error_string= `Error: ${e.response.status}. ${e.response.data}`
          }else{  //Axios error
            error_string= `Internal axios error: ${e}`
            console.log(e)
          }
          this.errormsg= error_string;
        }
      }
    },
  }
};
</script>

<template>
  <div class="main-container">
    <div v-if="!loading" class="login-container">
      <!-- Form di login -->
      <form class="login-form" @submit.prevent="doLogin">
        <h3 class="h3">Login</h3>
        <label for="username">Nome utente:  <span v-if="!isUsernameValid && username" class="error">Nome non valido</span></label>
        <input id="username" v-model="username" type="text" placeholder="Inserisci il nome utente" required>
        <button type="submit" :disabled="!isUsernameValid || loading" :class="{ disabled: !isUsernameValid || loading }">Login</button>
      </form>
    </div>
    <!-- Mostra messaggi di errore -->
    <ErrorMsg v-if="errormsg" :msg="errormsg" @close="errormsg= null" />

    <!-- Spinner di caricamento -->
    <LoadingSpinner :loading="loading" :loading-text="'Benvenuto/a '+ username +'! Caricamento in corso'" />
  </div>
</template>



<style scoped>
.main-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (min-width: 2000px) {
  .main-container {
    max-width: 1400px;
  }
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