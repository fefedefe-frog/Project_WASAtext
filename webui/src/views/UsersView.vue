<script>
export default {
  data: function () {
    return {
      token: localStorage.getItem('authToken'),
      usrId: localStorage.getItem('usrId'),
      errormsg: null,
      loading: false,
      users: [],
      isOverlayVisible: false,
      showMessage: false,
      noUsers: false,
    }
  },
  mounted() {
    this.refresh();
  },
  methods: {
    async refresh() {
      this.loading= true;
      this.showMessage= true;

      this.noUsers= false;
      this.errormsg= null;

      this.users= []

      // Aggiungo dei delay grafici per rendere l'interfaccia più gradevole, in questo caso il delay
      // aspetterà minimo 500ms prima di mostrare i risultati della richiesta al backend che verrà eseguita
      // durante questo delay visivo
      const delay = new Promise(resolve => setTimeout(resolve, 500));
      const fetchUsers = this.$axios.get("/users", { headers: { Authorization: `${this.token}` } });

      try {
        let response= await Promise.all([fetchUsers, delay]).then(results => results[0]);

        if (response.data){
          this.showMessage= false;
          response.data["users"].forEach(user => {
            if(user.usrId !== this.usrId){
              this.users.push(user);
            }
          })
        }else{
          this.users= [];
        }
      } catch (e) {
        if (e.response) {
          if (e.response.status === 404) {
            // In questo caso, se ho l'errore 404 significa solo che non ci sono utenti nel database
            this.noUsers= true;

          }else{
            this.errormsg= e.toString();
          }
        }else{
          this.errormsg= e.toString();
        }
      } finally {
        this.loading= false;
      }
    },
    toggleOverlay() {
      this.isOverlayVisible = !this.isOverlayVisible;
    }
  }
}
</script>

<template>
  <div class="container-main">
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Users</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
            Refresh
          </button>
          <div class="btn-group me-2">
            <button type="button" class="btn btn-sm btn-outline-primary" @click="toggleOverlay">
              New Chat
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="container">
      <div :class="showMessage ? 'user-list-show-message' : 'users-list-container'">
        <div v-if="isOverlayVisible" class="overlay">
          <!-- overlay per startare una nuova chat o vedere le info di un utente -->
        </div>

        <div v-if="showMessage && noUsers" class="nousers">
          <p class="visually-visible centered">Non ci sono utenti nel sistema :(</p>
        </div>

        <LoadingSpinner v-if="loading" :loading="loading" loading-text="Caricamento users" /><LoadingSpinner />
        <UserBanner v-for="user in users" :key="user.usrId" :user="user" />
      </div>

      <ErrorMsg v-if="errormsg" :msg="errormsg" />
    </div>
  </div>
</template>

<style scoped>
.container-main {
  height: 100%;
  width: 100%;
}

.container {
  width: 100%;
  height: 85%;
  overflow-y: scroll;
  overflow-x: hidden;
}
.users-list-container{
  position: relative;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
}

/* usato per quando devo mostare un messaggio centrato nel div della lista */
.user-list-show-message{
  display: flex;
  justify-content: center;
  align-items: center;
}

.nousers {
  position: relative;
  width: fit-content;
  background-color: #f0f0f0;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  font-size: 18px;
  color: #333;
}

.overlay{
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(123, 174, 40, 0.5); /* Overlay semi-trasparente */
  z-index: 1;
}
</style>