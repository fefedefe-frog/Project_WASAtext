<script>
export default {
  data: function () {
    return {
      token: localStorage.getItem('authToken'),
      errormsg: null,
      loading: false,
      chats: [],
      isOverlayVisible: false,
      showMessage: false,
      noChats: false,
    }
  },
  methods: {
    async refresh() {
      this.loading= true;
      this.showMessage= true;

      this.noChats= false;
      this.errormsg= null;

      this.chats= []

      const delay = new Promise(resolve => setTimeout(resolve, 500));
      const fetchChats = this.$axios.get("/chats", {headers: {Authorization: `${this.token}`}});

      try {
        let response= await Promise.all([fetchChats, delay]).then(results => results[0]);

        if (response.data){
          this.showMessage= false;
          response.data["chats"].forEach(chat => {
            this.chats.push(chat)
          })
        }
      } catch (e) {
        if (e.response) {
          if (e.response.status === 404) {
            // In questo caso, se ho l'errore 404 significa solo che l'utente non ha chat, e non è un problema
            this.noChats= true;

          }else{
            this.errormsg= e.toString();
          }
        }else{
          this.errormsg= e.toString();
        }
      } finally {
        this.loading= false
      }
    }
  },
  mounted() {
    this.refresh()
  }
}
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Your chats</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Div che contiene la lista di chat -->
    <div :class="showMessage ? 'chat-list-show-message' : 'chat-list-container'">
      <div v-if="showMessage && noChats" class="nochat">
        <p class="visually-visible centered">Non hai nessuna chat :(</p>
      </div>

      <LoadingSpinner v-if="loading" :loading="loading" loadingText="Caricamento chats"/><LoadingSpinner/>
      <chatBanner v-for="chat in chats" :chat="chat" lastMessage="franco"></chatBanner>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style scoped>
.chat-list-container{
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
}

/* usato per quando devo mostare un messaggio centrato nel div della lista */
.chat-list-show-message{
  display: flex;
  justify-content: center;
  align-items: center;
}

.nochat {
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
</style>