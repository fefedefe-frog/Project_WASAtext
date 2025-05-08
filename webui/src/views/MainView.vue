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

      showUserInfo: false,
      loadedUserInfo: {},

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
    // this.setIntervalId= setInterval(async () => {
    //   await this.getUsers();
    //   await this.getUserChats();
    // }, 20000);
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
    prepSendMessage(rawInput){
      console.log(rawInput)
      let userToSend= rawInput['sendTo'];
      let messageData= rawInput['messageData']

      // Controllo se ho già una chat diretta con l'utente a cui voglio mandare il messaggio
      let chatToSend= -1;
      if (this.userChats.length > 0){
        let filteredChat= this.userChats.filter(chat => chat['isGroup'] === false);
        filteredChat.forEach(chat => {
          if(chat['participants'].includes(userToSend)){
            console.log("la chat esiste già: "+ chat['chatId']);
            chatToSend= chat['chatId'];
          }
        });
      }

      // Preparo il formData per la richiesta
      const requestFormData= new FormData();
      if (chatToSend === -1){ // Non esiste una chat diretta con l'utente

        // Assegno le informazioni sulla chat
        requestFormData.append('chatInfo', JSON.stringify({
          chatName: "",
          chatPhoto: "",
          isGroup: false,
          participants: [userToSend]
        }));

        // Assegno le informazioni sul messaggio
        requestFormData.append('contentType', messageData['contentType']);
        requestFormData.append('content', messageData['content']);

        this.startNewChat(requestFormData);
      }else { // La chat esiste già e quindi procedo a preparare il form per l'invio del messaggio
        // Assegno le informazioni sul messaggio
        requestFormData.append('contentType', messageData['contentType']);
        requestFormData.append('content', messageData['content']);
        requestFormData.append('respondTo', -1);

        this.sendMessage(chatToSend, requestFormData);
      }


    },
    async startNewChat(formData){
      console.log("nuova chat");
      console.log(formData);
      this.errormsg= null;
      try{
        console.log("invio dati")
        let response= await this.$axios.post(`/chats`, formData, {
          headers: {
            Authorization: this.token,
          }
        });

        if(response.data){
          console.log("risposta")
          console.log(response.data);
        }
      }catch (e){
        console.log(e)
        this.errormsg= e;
      }
    },
    async sendMessage(chatId, formData){
      console.log("invio messaggio")
      this.errormsg= null;

      try{
        let response= await this.$axios.post(`/chats/${chatId}/messages`, formData, {
          headers: {
           Authorization: this.token,
          }
        });

        if(response.data){
          console.log(response.data);
        }
      }catch (e){
        this.errormsg= e;
      }
    },
    loadChat(bannerData){
      this.showChat= false
      this.loadedChatInfo= bannerData['chatData'];
      this.loadedChatMessages= bannerData['messages'];

      this.showChat= true;
    },
    doLogout(){
      sessionStorage.removeItem("authToken");
      sessionStorage.removeItem("usrId");
      this.$router.push('/login');
    },
    loadUserInfo(bannerData){
      this.showUserInfo= false;
      this.loadedUserInfo= bannerData;

      this.showUserInfo= true;
    },
    closeUserInfo(){
      this.showUserInfo= false;
      this.loadedUserInfo= null;
    },

    componentsErrorHandler(error){
      this.errormsg= error.toString;
    }
  }
}
</script>

<template>
  <div class="main-container">
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
            <userBanner v-for="user in userFilteredResult" :key="user.usrId" :user-data="user" @userClicked="loadUserInfo"/>
          </div>
        </div>
      </div>

      <div class="list-footer">
        <div class="btn-toolbar mb-2 mb-md-0">
          <div class="btn-group me-2">
            <button type="button" class="btn btn-sm btn-primary shadow-none" @click="console.log('TODO: PROFILO')">
              Profilo
            </button>
            <button type="button" class="btn btn-sm btn-danger shadow-none" @click="doLogout">
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="chat-container bobby">
      <UserInfo v-if="showUserInfo" :key="loadedUserInfo['usrId']" :user-data="loadedUserInfo" @close-user-info="closeUserInfo" @reqNewChat="prepSendMessage"/>
      <Chat v-if="showChat" :key="loadedChatInfo['chatId']" :initial-messages="loadedChatMessages" :chat-data="loadedChatInfo" @close-chat="closeChat" />
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
.main-container {
  padding: 0.7rem;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row;

  position: relative;
}

@media (min-width: 2000px) {
  .main-container {
    max-width: 1400px;
  }
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

.list-footer{
  border-top: 2px solid grey;
  padding: 5px;
  width: 100%;
  display: flex;
  justify-content: center;
}

.chat-container {
  height: 100%;
  width: 100%;

  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

</style>
