<script>
export default {
  components: {UserBanner},
  data: function () {
    return {
      token: localStorage.getItem('authToken'),
      errormsg: null,
      loading: false,
      users: [],
    }
  },
  methods: {
    async refresh() {
      this.loading = true
      this.errormsg = null

      try {
        let response = await this.$axios.get("/users", {headers: {Authorization: `${this.token}`}});
        this.users = []
        if (response.data){
          response.data["users"].forEach(user => {
            this.users.push(user)
          })
        }
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
        this.loading = false;
      }
    },
    showOverlay() {
      this.isOverlayVisible = true;
    },
    hideOverlay() {
      this.isOverlayVisible = false;
    },

  },
  mounted() {
    this.refresh()
  }
}
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Users</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
            Refresh
          </button>
          <div class="btn-group me-2">
            <button type="button" class="btn btn-sm btn-outline-primary" @click="showOverlay">
              New Chat
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="container">
      <div id="users_list_container">
        <UserBanner v-for="user in users" :user="user"></UserBanner>

        <LoadingSpinner v-if="loading" loading="{{ loading }}" loadingText="Caricamento users"/><LoadingSpinner/>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>

</template>

<style scoped>
#users_list_container{
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  border: #0c4128 4px solid;
}

/* Contenitore principale */
.container {
  position: relative;
  width: 100%;
  height: auto;
  padding: 20px;
  border: 4px solid darkorchid;
}

.overlay {
  position: fixed; /* Cambia a "fixed" per far sì che stia sopra l'intera finestra */
  top: 50%; /* Posizione centrata verticalmente */
  left: 50%; /* Posizione centrata orizzontalmente */
  transform: translate(-50%, -50%); /* Centra esattamente l'overlay */
  background-color: rgba(0, 0, 0, 0.7); /* Sfondo semi-trasparente */
  color: white;
  padding: 20px;
  border-radius: 10px;
  border: #0a58ca 6px solid;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  z-index: 100; /* Assicurati che l'overlay stia sopra l'altro contenuto */
}

.modal {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  color: #333;
}

</style>