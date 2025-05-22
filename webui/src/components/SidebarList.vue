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
      }, 15000);
    }else {
      this.getChats();
      this.setIntervalId= setInterval(async () => {
        await this.getChats();
      }, 15000);
    }
  },
  beforeUnmount() {
    clearInterval(this.setIntervalId);
  },
  methods: {
    async getUsers() {
      this.errormsg= null;

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
          this.errormsg = e;
        }
      }finally {
        this.loading = false;
      }
    },
    async getChats() {
      this.errormsg= null;

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
          console.log("user doesn't have any chat");
        }else {
          this.errormsg = e;
        }
      }finally {
        this.loading = false;
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
        v-for="item in filteredResult"
        :key="item[this.items === 'users' ? 'usrId' : 'chatId']"
        :is="bannerComponent"
        :inputData="item"
        @bannerClicked="bannerClicked"
    />
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
</style>