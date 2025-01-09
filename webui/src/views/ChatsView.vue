<script>
export default {
  data: function () {
    return {
      token: localStorage.getItem('authToken'),
      errormsg: null,
      loading: false,
      chats: [],
      isOverlayVisible: false,
    }
  },
  methods: {
    async refresh() {
      this.loading = true
      this.errormsg = null

      try {
        let response = await this.$axios.get("/chats", {headers: {Authorization: `${this.token}`}});
        this.chats = []
        if (response.data) {
          response.data["chats"].forEach(chat => {
            this.chats.push(chat)
          })
        }
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
        this.loading = false;
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
    <div id="chat_list_container">
      <chatBanner v-for="chat in chats" :chatId="chat['chatId']" :chatName="chat['chatName']" :chatPhoto="chat['chatPhoto']" lastMessage="franco"></chatBanner>
    </div>


    <LoadingSpinner v-if="loading" loading="{{ loading }}" loadingText="Caricamento chats"/><LoadingSpinner/>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style scoped>
#chat_list_container{
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
}

.container{
  position: relative;
  height: 100%;
  width: 100%;
}
</style>