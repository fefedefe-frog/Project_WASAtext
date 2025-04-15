<script>
export default {
  data: function () {
    return {
      token: '',
      usrId: '',
      searchQuery: '',
      errormsg: null,
      loading: false,
      userChats: [],
      users: [],

      showChat: false,
      loadedChatInfo: {},
      loadedChatMessages: [],

      setIntervalId: null,
    }
  },
  computed: {
    chatFilteredResult(){
      let searchPool= this.userChats
      let query= this.searchQuery.trim();

      if (query === ''){
        return searchPool
      }
      return searchPool.filter(chat =>
          chat['chatName'].toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
    userFilteredResult(){
      let searchPool= this.users;
      let query= this.searchQuery.trim();

      if (query === ''){
        return searchPool;
      }
      return searchPool.filter(user =>
        user['userName'].toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.usrId= sessionStorage.getItem('usrId');

    this.getUsers();
    this.getUserChats();
    this.setIntervalId= setInterval(async () => {
      this.getUsers();
      this.getUserChats();
    }, 19000);
  },
  beforeUnmount() {
    clearInterval(this.setIntervalId);
  },
  methods: {
    async getUsers() {
      this.errormsg= null

      try {
        let response= await this.$axios.get(`/users`, {headers: {Authorization: this.token}});
        if (response.data) {
          this.users= [];

          response.data['users'].forEach(user => {
            this.users.push(user);
          });
        }
      }catch(e) {
        if (e.status === 404){
          this.users= [];
        }else {
          this.errormsg = e;
        }
      }finally {
        this.loading = false
      }
    },
    async getUserChats() {
      this.errormsg= null

      try {
        let response= await this.$axios.get(`/chats`, {headers: {Authorization: this.token}});
        this.userChats= [];

        if (response.data) {
          if (response.data['chats']){
            response.data['chats'].forEach(chat => {
              this.userChats.push(chat);
            });
          }
        }

      }catch(e) {
        if (e.status === 404){
          this.userChats= [];
          console.log("user doesn't have any chat");
        }else {
          this.errormsg = e;
        }
      }finally {
        this.loading = false
      }
    },
    async closeChat(leave){
      this.showChat= false
      if (leave){
        this.errormsg = null

        try {
          let response= await this.$axios.delete(`/chats/${this.chatId}/users`, {
            headers: {Authorization: this.token}
          });
          if (response.status < 400){
            await this.getUserChats();
          }
        }catch(e) {
          this.errormsg = e;
        }
      }
    },
    loadChat(bannerData){
      this.showChat= false
      this.loadedChatInfo= bannerData['chatData'];
      this.loadedChatMessages= bannerData['messages'];

      this.showChat= true;
    },
    componentsErrorHandler(error){
      this.errormsg= error.toString;
    }
  }
}
</script>

<template>
  <div class="container">
    <div class="lists bobby">
      <div class="search-box">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
        <input v-model="searchQuery" type="text" placeholder="Inserisci nome utente o chat" required>
      </div>

      <ul id="tab" class="nav nav-tabs" role="tablist">
        <li class="nav-item" role="presentation">
          <button class="nav-link active" data-bs-toggle="tab" data-bs-target="#chats" type="button">Chats</button>
        </li>
        <li class="nav-item" role="presentation">
          <button class="nav-link" data-bs-toggle="tab" data-bs-target="#users" type="button">Users</button>
        </li>
      </ul>

      <div id="tabContent" class="tab-content">
        <div id="chats" class="tab-pane fade show active" role="tabpanel">
          <div class="chats-list">
            <chatBanner v-for="chat in chatFilteredResult" :key="chat.chatId" :chat-data="chat" @error="componentsErrorHandler" @chat-banner-data="loadChat" />
          </div>
        </div>
        <div id="users" class="tab-pane fade" role="tabpanel">
          <div class="users-list">
            <userBanner v-for="user in userFilteredResult" :key="user.usrId" :user-data="user" />
          </div>
        </div>
      </div>
    </div>
    <div class="chat-container bobby">
      <Chat v-if="showChat" :key="loadedChatInfo['chatId']" :initial-messages="loadedChatMessages" :chat-data="loadedChatInfo" @close-chat="closeChat" />
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
.container {
  display: flex;
  flex-direction: row;

  position: relative;
  height: 100%;
  width: 100%;
}

.lists {
  display: flex;
  flex-direction: column;

  height: 100%;
  width: 25%;
  padding: 5px;

  margin-right: 5px;
}

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


.tab-content {
  height: 100%;

  justify-content: center;
  align-items: center;

  padding-top: 5px;

  overflow-y: auto;

}

.chats-list .users-list{
  height: 100%;
  width: 100%;
  padding: 5px;

  overflow: hidden;
  overflow-y: auto;
}


.chat-container {
  height: 100%;
  width: 100%;

  overflow: hidden;
}

</style>
