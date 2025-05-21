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

      showCreateGroup: false,
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
    this.getChats();
    this.setIntervalId= setInterval(async () => {
      await this.getUsers();
      await this.getChats();
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
    async getChats() {
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
    prepSendMessage(rawInput){
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
    componentsErrorHandler(error){
      this.errormsg= error;
    },
  }
}
</script>

<template>
  <div class="main-container">
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

.create-group {
  position: sticky;
  top: 0;
  background: white;
  z-index: 1;
  width: 100%;
  height: 50px;
  border: 1px solid red;
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