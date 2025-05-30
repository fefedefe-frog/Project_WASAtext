<script>
export default {
  data: function () {
    return {
      token: '',
      errormsg: null,
      loading: false,
      users: [],

      chatInfo: {
        chatName: "",
        chatPhoto: "",
        isGroup: false,
        participantsId: [],
      },

      selectedParticipantInfo: [],

      initialMessage: {
        textContent: "",
        photoContent: null,
      },

      setIntervalId: null,
    }
  },
  computed: {
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
    this.setIntervalId= setInterval(async () => {
      await this.getUsers();
    }, 5000);
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
    async startNewChat(){
      this.errormsg= null;

      // Preparo il formData per la richiesta
      const requestFormData= new FormData();


      // Assegno le informazioni sulla chat
      let chatName= ""
      let chatPhoto= new Blob([], {type: 'image/png'});
      if (this.chatInfo['isGroup']){
        chatName= this.chatInfo['chatName']
        chatPhoto= this.chatInfo['chatPhoto']
      }
      requestFormData.append('chatName', chatName);
      requestFormData.append('chatPhoto', chatPhoto);
      requestFormData.append('isGroup', this.chatInfo['isGroup']);
      requestFormData.append('participants', this.chatInfo['participantsId']);

      // Assegno le informazioni sul messaggio
      requestFormData.append('messageTextContent', this.initialMessage['textContent']);
      requestFormData.append('messagePhotoContent', this.initialMessage['photoContent']);

      try{
        let response= await this.$axios.post(`/chats`, requestFormData, {
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
          this.$router.push(`/chats/${newChat['chatId']}`);
        }
      }catch (e){
        this.errormsg= e;
      }
    },
    addParticipant(bannerData){

      // Aggiungo l'utente alla lista di utenti, solo se non è già presente
      if (!this.chatInfo['participantsId'].some(usrId => usrId === bannerData['usrId'])){
        this.selectedParticipantInfo.push(bannerData);
        this.chatInfo['participantsId'].push(bannerData['usrId'])
      }
    },
    removeParticipant(id){
      this.selectedParticipantInfo= this.selectedParticipantInfo.filter(item =>
        item['usrId'] !== id
      );
      this.chatInfo['participantsId']= this.chatInfo['participantsId'].filter(pId =>
          pId !== id
      );
    },
    componentsErrorHandler(error){
      this.errormsg= error;
    },
  }
}
</script>

<template>
  <div class="main-container bobby">
    <div class="select-participant">
      <sidebarList :bannerComponent="'userBanner'" items="users" @error="componentsErrorHandler" @bannerData="addParticipant"/>
    </div>
    <div class="new-chat-form">
      <div class="new-chat-info">
        <div class="new-chat-image">
          immagine selezionabile, se cliccato
          apre un menù per cambiare immagine
        </div>
        <div class="new-chat-text">
          <div class="new-chat-name">
            nome
          </div>
          <div class="participants">
            <span v-for="p in selectedParticipantInfo" :key="p['usrId']" class="participant">
              <img :src="'data:image/png;base64,'+ p['userPhoto']" alt="Profile Image" draggable="false">
              <span>{{p['userName']}}</span>
              <button @click="removeParticipant(p['usrId'])">
                <svg class="feather"> <use href="/feather-sprite-v4.29.0.svg#x"/></svg>
              </button>
            </span>
          </div>
        </div>
      </div>
      <div class="initial-message">
        <div class="roberto">
          <messageForm/>
        </div>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
.main-container {
  height: 100%;
  width: 100%;

  padding: 10px;
  display: flex;
  flex-direction: row;
}

@media (min-width: 2000px) {
  .main-container {
    max-width: 1400px;
    max-height: calc(1400px * 9 / 16);
  }
}

.select-participant {
  border: 1px solid orange;
  width: 25%;
  padding: 4px;
}

.new-chat-form{
  width: 75%;
  height: 100%;
  border: 1px solid black;
}

.initial-message {
  width: 75%;
  height: 50%;

  border: 1px solid darkgreen;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.roberto {

  width: 50%;
  height: 30%;
}

.new-chat-info {
  border: 1px solid red;

  width: 75%;
  height: 50%;

  padding: 10px;

  display: flex;
  flex-direction: row;
  align-items: center;
}

.new-chat-image {
  height: 70%;
  aspect-ratio: 1/1;
  border: 1px solid blue;
}

.new-chat-text {
  height: 70%;
  width: 70%;

  margin-left: 5px;

  border: 1px solid violet;

  display: flex;
  flex-direction: column;
}

.new-chat-name {
  width: 45%;
  height: 60%;
  border: 1px solid cyan;

}

.participants {
  width: 100%;
  height: 40%;
  border: 1px solid peru;

  padding-left: 5px;
  padding-right: 5px;
  gap: 5px;

  display: flex;
  flex-direction: row;
  align-items: center;

  overflow: hidden;
  overflow-x: scroll;
}

.participant {
  width: fit-content;
  height: fit-content;

  max-width: 100px;

  padding: 3px;

  background: skyblue;
  border-radius: 10px;

  display: flex;
  flex-direction: row;
  align-items: center;
}

.participant:hover {
  background: dodgerblue;
}

.participant img {
  height: 20px;
  width: 20px;

  border-radius: 10px;
  user-select: none;
}

.participant span{
  font-size: 14px;
  text-align: center;

  margin-left: 2px;

  user-select: none;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.participant button {
  background: none;
  border: none;
  width: 20px;
  height: 20px;
  border-radius: 20px;

  cursor: pointer;

  display: inline-flex;
  align-items: center;
  justify-content: center;

  color: red;
}

.participant svg {
  stroke-width: 3;
}

.participant button:hover {
  background: rgba(0, 0, 0, 0.5);
}
</style>