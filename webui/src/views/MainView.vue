<script>
export default {
  data: function () {
    return {
      token: '',
      loggedUser: {
        usrId: "",
        userName: "",
        userPhoto: "",
      },
      searchQuery: '',
      errormsg: null,
      loading: false,
      userChats: [],
      users: [],

      showChat: false,
      loadedChatInfo: {
        chatInfo: {},
        participantNames: {},
        messages: [],
      },

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

    // Recupero le info dell'utente
    let storedUser= sessionStorage.getItem('loggedUser');
    if (storedUser){
      this.loggedUser= JSON.parse(storedUser);
    }
    sessionStorage.removeItem('loggedUser');

    this.getUsers();
    this.getUserChats();
    this.setIntervalId= setInterval(async () => {
      await this.getUsers();
      await this.getUserChats();
    }, 20000);
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
        if (response.data) {
          if (response.data['chats']){
            this.userChats= [];
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
    async leaveChat(chatId){
      this.errormsg = null

      try {
        let response= await this.$axios.delete(`/chats/${chatId}/users`, {
          headers: {Authorization: this.token}
        });
        if (response.status < 400){
          await this.getUserChats();
        }
      }catch(e) {
        this.errormsg = e;
      }
    },
    closeChat(leave){
      this.showChat= false
      if (leave){
        this.leaveChat(this.loadedChatInfo['chatInfo']['chatId']);
      }
      this.loadedChatInfo= {}
    },
    prepSendMessage(rawInput){  // TODO adattarla alla creazione dei gruppi
      let userToSend= rawInput['sendTo'];
      let messageData= rawInput['messageData']

      // Controllo se ho già una chat diretta con l'utente a cui voglio mandare il messaggio
      let chatToSend= -1;
      if (this.userChats.length > 0){
        let filteredChat= this.userChats.filter(chat => chat['isGroup'] === false);
        filteredChat.forEach(chat => {
          if(chat['participants'].some(p => p['usrId'] === userToSend)){
            console.log("la chat esiste già: "+ chat['chatId']);
            chatToSend= chat['chatId'];
          }
        });
      }

      // Preparo il formData per la richiesta
      const requestFormData= new FormData();

      // Assegno le informazioni sul messaggio
      requestFormData.append('textContent', messageData['textContent']);
      requestFormData.append('photoContent', messageData['photoContent']);
      requestFormData.append('respondTo', -1);

      if (chatToSend === -1){ // Non esiste una chat diretta con l'utente

        // Assegno le informazioni sulla chat
        requestFormData.append('chatName', "");

        let emptyPhoto= new Blob([], {type: 'image/png'});
        requestFormData.append('chatPhoto', emptyPhoto);
        requestFormData.append('isGroup', false);
        requestFormData.append('participants', [userToSend]);


        this.startNewChat(requestFormData);
        this.showUserInfo= false;
      }else { // La chat esiste già e quindi la richiesta http sarà all'endpoint per inviare un messaggio
        this.sendMessage(chatToSend, requestFormData);
        this.showUserInfo= false;

      }
    },
    async startNewChat(formData){
      this.errormsg= null;

      try{
        let response= await this.$axios.post(`/chats`, formData, {
          headers: {
            Authorization: this.token,
          }
        });

        if(response.data){
          let newChat= {
            chatId: response.data['chatId'],
            isGroup: response.data['isGroup'],
            chatName: response.data['chatName'],
            chatPhoto: response.data['chatPhoto'],
            participants: response.data['participants']
          }
          this.userChats.push(newChat)
        }
      }catch (e){
        this.errormsg= e;
      }
    },
    async sendMessage(chatId, formData){
      this.errormsg= null;

      try{
        let response= await this.$axios.post(`/chats/${chatId}/messages`, formData, {
          headers: {
           Authorization: this.token,
          }
        });
      }catch (e){
        this.errormsg= e;
      }
    },
    loadChat(chatBannerData){
      this.showChat= false
      this.loadedChatInfo= chatBannerData;
      this.showChat= true;
    },
    doLogout(){
      sessionStorage.removeItem("authToken");
      sessionStorage.removeItem("usrId");
      this.$router.push('/login');
    },
    loadUserInfo(userBannerData){
      this.showUserInfo= false;
      this.loadedUserInfo= userBannerData;

      this.showUserInfo= true;
    },
    loadMyInfo(){
      this.showUserInfo= false;
      this.loadedUserInfo= this.loggedUser;

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
          <div class="banner-lists">
            <chatBanner v-for="chat in chatFilteredResult" :key="chat.chatId" :chat-data="chat" @error="componentsErrorHandler" @chat-banner-data="loadChat" />

          </div>
        </div>
        <div id="users" class="tab-pane fade" role="tabpanel">
          <div class="banner-lists">
            <userBanner v-for="user in userFilteredResult" :key="user.usrId" :user-data="user" @userClicked="loadUserInfo"/>
          </div>
        </div>
      </div>

      <div class="list-footer">
        <div class="btn-toolbar mb-2 mb-md-0">
          <div class="btn-group me-2">
            <button type="button" class="btn btn-sm btn-primary shadow-none" @click="loadUserInfo(this.loggedUser)">
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
      <Chat v-if="showChat" :key="loadedChatInfo['chatInfo']['chatId']" :chat-data="loadedChatInfo" @close-chat="closeChat" />
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
.main-container {
  height: 100%;
  width: 100%;

  padding: 0.7rem;
  display: flex;
  flex-direction: row;

  position: relative;
}

@media (min-width: 2000px) {
  .main-container {
    max-width: 1400px;
    max-height: calc(1400px * 9 / 16);
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

.banner-lists{
  height: fit-content;
  width: 100%;
  padding: 5px;
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