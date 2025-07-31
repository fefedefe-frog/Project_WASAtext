<script>
export default {
  props: {
    items: {
      type: String,
      required: true
    },
    bannerComponent: {
      type: [String, Object],
      required: true
    }
  },
  emits: ['error', 'bannerData'],
  data: function () {
    return {
      token: '',
      searchQuery: '',
      chats: [],
      users: [],

      setIntervalId: null,
    }
  },
  computed: {
    filteredResult(){
      let searchPool= [];
      let searchKey= "";
      if (this.items === 'users'){
        searchPool= this.users;
        searchKey= 'userName';
      }else {
        searchPool= this.chats;
        searchKey= 'chatName';
      }

      let query= this.searchQuery.trim();
      if (query === ''){
        return searchPool;
      }
      return searchPool.filter(item =>
          item[searchKey].toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.usrId= sessionStorage.getItem('usrId');

    if (this.items === 'users'){
      this.getUsers();
      this.setIntervalId= setInterval(async () => {
        await this.getUsers();
      }, 10000);
    }else {
      this.getChats();
      this.setIntervalId= setInterval(async () => {
        await this.getChats();
      }, 10000);
    }
  },
  beforeUnmount() {
    clearInterval(this.setIntervalId);
  },
  methods: {
    async getUsers() {
      try {
        let response= await this.$axios.get(`/users`, {headers: {Authorization: this.token}});

        if (response.data) {
          if(response.data['users']){
            this.users= [];
            response.data['users'].forEach(user => {
              this.users.push(user);
            });
          }
        }
      }catch(e) {
        if (e.status === 404){
          this.users= [];
        }else {
          this.$emit('error', e);
        }
      }
    },
    async getChats() {
      try {
        let response= await this.$axios.get(`/chats`, {headers: {Authorization: this.token}});
        if (response.data) {
          if (response.data['chats']){
            this.chats= [];
            response.data['chats'].forEach(chat => {
              this.chats.push(chat);
            });
          }
        }

      }catch(e) {
        if (e.status === 404){
          this.chats= [];
        }else {
          this.$emit('error', e);
        }
      }
    },
    bannerClicked(bannerData){
      this.$emit('bannerData', bannerData);
    },
  }
}
</script>

<template>
  <div class="search-box">
    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
    <input v-model="searchQuery" type="text" :placeholder="items === 'users' ? 'cerca utente' : 'cerca chat'" required>
  </div>

  <div class="banner-lists">
    <component
      :is="bannerComponent"
      v-for="item in filteredResult"
      :key="item[items === 'users' ? 'usrId' : 'chatId'] + '-' + Math.floor(Math.random() * 10) "
      :input-data="item"
      @banner-clicked="bannerClicked"
    />
  </div>

  <div class="list-footer">
    <button type="button" class="btn btn-sm btn-primary shadow-none" @click="items === 'users' ? getUsers : getChats">
      <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#rotate-cw" /></svg> Ricarica {{ items === 'users' ? "utenti" : "chat" }}
    </button>
  </div>
</template>

<style scoped>
.search-box {
  position: relative;
  width: 100%;
  max-width: 300px;
  margin-bottom: 5px;
}

.search-box input {
  width: 100%;
  height: 30px;
  padding: 10px 40px 10px 20px; /* Spazio per l'icona */
  border: 1px solid #ccc;
  border-radius: 50px; /* Arrotonda i bordi */
  outline: none;
}

.search-box svg {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #888;
}

.banner-lists{
  height: fit-content;
  width: 100%;
  padding: 5px;

  overflow: hidden;
  overflow-y: scroll;
}

.list-footer{
  border-top: 2px solid gray;
  width: 85%;
  height: 10%;
  margin-bottom: 0;
  margin-top: auto;

  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}
</style>