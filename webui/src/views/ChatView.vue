<script>

import ForwardMessage from "../components/ForwardMessage.vue";

export default {
  components: {ForwardMessage},
  data: function () {
    return {
      usrId: '',
      token: '',
      errormsg: null,
      loading: false,

      chat: {
        chatId: -1,
        chatName: "",
        chatPhoto: "",
        isGroup: false,
        participants: []
      },
      participantNames: {},

      lastMsgId: -1,
      prevLastMsgId: -1,
      messages: [],

      respondTo: -1,
      respondMessageData: {},

      forwardMsgId: -1,

      getChatInfoIntervalId: null,
      getMessagesIntervalId: null,
      getMessagesIsRunning: false,

      addParticipantPanel: false,
      changeGroupInfo: false,
      forwardMessagePanel: false,

      newGroupName: "",
    }
  },
  computed: {
    groupNameIsValid() {
      const name=  this.newGroupName;
      return (name.length >= 3 && name.length <= 16 && ((/^\S.*\S$/).test(name)));
    },
  },
  watch: {
    '$route.params.chat_id': {
      immediate: true,
      async handler(newChatId, oldChatId) {
        if (newChatId !== oldChatId) {
          if (this.token === "") this.token = sessionStorage.getItem('authToken');

          this.loading = true;
          this.chat['chatId'] = newChatId;

          await this.getChatInfo();
          this.updateParticipantNamesDict();
        }
      }
    },
    respondTo(newId, oldId){
      if (newId !== oldId && newId !== -1){
        let respondMessage= this.messages.filter(message => message['msgId'] === newId)[0];
        if (respondMessage){
          this.respondMessageData= {
            senderName: this.participantNames[respondMessage['senderId']],
            textContent: respondMessage['textContent'],
            photoContent: respondMessage['photoContent'],
          };
        }else {
          this.respondTo= -1;
        }
      }
    },
    forwardMsgId(newId, oldId){
      if (newId !== oldId && newId !== -1){
        this.addParticipantPanel= false;
        this.changeGroupInfo= false;

        this.forwardMessagePanel= true;
      }
    }
  },
  async mounted(){
    this.chat['chatId']= this.$route.params.chat_id;
    this.usrId= sessionStorage.getItem('usrId')
    this.token= sessionStorage.getItem('authToken');

    this.loading= true;
    await this.getChatInfo();
    this.updateParticipantNamesDict();


    this.getChatInfoIntervalId= setInterval( async () => {
      await this.getChatInfo();
      this.updateParticipantNamesDict();
    }, 30000);


    // Funzione per avviare un nuovo intervallo di ricezione dei messaggi
    await this.getMessagesSetInterval();
  },
  beforeUnmount() {
    clearInterval(this.getChatInfoIntervalId);
    clearInterval(this.getMessagesIntervalId);
  },
  methods: {
    async getMessagesSetInterval(){
      if (this.getMessagesIntervalId){
        clearInterval(this.getMessagesIntervalId);
      }

      await this.getMessages();
      await this.updateReadStatus();
      this.getMessagesIntervalId= setInterval( async () => {
        this.lastMsgId= -1; //Per ora recupero tutti i messagi della chat in ogni caso
        await this.getMessages();
        await this.updateReadStatus();
      }, 5000);
    },
    async getChatInfo() {
      try {
        let response = await this.$axios.get(`/chats/${this.chat['chatId']}`, {
          headers: {Authorization: `${this.token}`}
        });

        if (response.data) {
          this.chat['chatName']= response.data['chatName'];
          this.chat['chatPhoto']= response.data['chatPhoto'];
          this.chat['isGroup']= response.data['isGroup'];
          this.chat['participants']= response.data['participants'];
        }
      } catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      } finally {
        this.loading= false;
      }
    },
    async getMessages() {
      // Evito la sovrapposizione della funzione quando chiamata dagli intervalli
      if (this.getMessagesIsRunning) return;
      this.getMessagesIsRunning= true;

      try {
        let response= await this.$axios.put(`/chats/${this.chat['chatId']}/messages`, {
          msgId: this.lastMsgId
        }, {
          headers: {Authorization: this.token}
        });

        if (response.data) {
          if (this.lastMsgId === -1){
            this.messages= [];
          }
          if (Array.isArray(response.data['messages']) && response.data['messages'].length > 0){
            response.data['messages'].forEach(message => {
              this.messages.push(message);
            });

            this.messages.reverse();
            this.lastMsgId= this.messages[0]['msgId'];
          }
        }
      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }
      this.getMessagesIsRunning= false;
    },
    async deleteMessage(msgId){

      // Recupero il messaggio dalla lista di messaggi
      let msg= this.messages.filter(message => message['msgId'] === msgId)[0];
      if (this.usrId !== msg['senderId']) {
        this.errormsg= "non puoi eliminare il messaggio di un altro utente"
        return;
      }

      try {
        await this.$axios.delete(`/chats/${this.chat['chatId']}/messages/${msgId}`, {
          headers: {Authorization: this.token}
        });
        this.messages= this.messages.filter(message => message['msgId'] !== msgId)

      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }
    },
    async updateReadStatus() {
      if (this.prevLastMsgId === this.lastMsgId) return;

      this.lastMsgId= this.messages[0]['msgId'];
      try {
        let response= await this.$axios.put(`/chats/${this.chat['chatId']}/messages/${this.lastMsgId}`, {}, {
          headers: {Authorization: this.token}
        });
        if (response.status > 400) {
          throw new Error("unable to update the messages status");
        }
      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }

      this.prevLastMsgId= this.lastMsgId;
    },
    async sendMessage(rawInput){
      const requestFormData= new FormData();
      requestFormData.append('textContent', rawInput['textContent']);
      requestFormData.append('photoContent', rawInput['photoContent']);
      requestFormData.append('respondTo', this.respondTo);

      try{
        let response= await this.$axios.post(`/chats/${this.chat['chatId']}/messages`, requestFormData, {
          headers: {
            Authorization: this.token,
          }
        });

        if(response.data){
          if (response.data['message']){
            this.messages.reverse();
            this.messages.push(response.data['message']);
            this.messages.reverse();
            this.lastMsgId= response.data['message']['msgId'];
          }
        }
      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }finally {
        this.respondTo= -1;
        this.respondMessageData= {};
      }

    },
    async leaveChat() {
      try {
        await this.$axios.delete(`/chats/${this.chat['chatId']}/users`, {
          headers: {Authorization: this.token}
        });
      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }finally {
        this.$router.replace('/chats');
      }
    },
    async addParticipant(participant){
      if (participant['usrId'] in this.participantNames){
        this.errormsg= new Error("Non puoi aggiungere un utente già partecipante").toString();
      }else {
        this.errormsg= null;

        try {
          let response= await this.$axios.post(`chats/${this.chat['chatId']}/users`,
              {usrId: participant['usrId']},
              {headers: {Authorization: this.token}}
          );

          if (response.data) {
            //Aggiorno la lista di partecipanti
            this.chat['participants']= response.data['participants']
            this.updateParticipantNamesDict();
          }
        }catch(e) {
          let error_string= ""
          if (e.response.status === 400 ||  //Bad request
              e.response.status === 401 ||  //Unauthorized
              e.response.status === 403 ||  //Forbidden
              e.response.status === 404 ||  //Not found
              e.response.status === 500){   //Internal server error
            error_string= `Error: ${e.response.status}. ${e.response.data}`;
          }else{  //Axios error
            error_string= `Internal axios error: ${e}`;
            console.log(e);
          }
          this.errormsg= error_string;
        }
      }
      this.addParticipantPanel= false;
    },
    imageUpload() {
      const input= document.createElement('input');
      input.type= "file";
      input.accept= "image/*";

      input.addEventListener("change", async (event) => await this.changeGroupPhoto(event));
      input.click();
    },
    async changeGroupPhoto(event){
      this.errormsg= null;

      if (!this.chat['isGroup']){
        this.errormsg= "non puoi cambiare l'immagine profilo di un altro utente";
        return
      }
      let oldGroupPhoto= this.chat['chatPhoto'];

      const file= event.target.files[0];
      if (file && file.type.startsWith("image/")) {

        // Faccio la richiesta per modificare l'immagine al backend
        // Preparo il formData per la richiesta
        const requestFormData= new FormData();
        requestFormData.append('newGroupPhoto', file);

        try{
          let response= await this.$axios.put(`/chats/${this.chat['chatId']}/propic`, requestFormData, {
            headers: {Authorization: this.token},
          });

          if (response.data){
            this.chat['chatPhoto']= response.data['chatPhoto']
          }
        }catch(e) {
          let error_string= ""
          if (e.response.status === 400 ||  //Bad request
              e.response.status === 401 ||  //Unauthorized
              e.response.status === 403 ||  //Forbidden
              e.response.status === 404 ||  //Not found
              e.response.status === 500){   //Internal server error
            error_string= `Error: ${e.response.status}. ${e.response.data}`
          }else{  //Axios error
            error_string= `Internal axios error: ${e}`
            console.log(e)
          }
          this.errormsg= error_string;
          this.chat['chatPhoto']= oldGroupPhoto;
        }
      }else {
        this.chat['chatPhoto']= oldGroupPhoto;
      }
    },
    async updateGroupName(){
      this.errormsg= null;

      if (!this.chat['isGroup']){
        this.errormsg= "non puoi cambiare l'immagine profilo di un altro utente";
        return
      }

      try{
        let response= await this.$axios.put(`/chats/${this.chat['chatId']}`, {
          newGroupName: this.newGroupName,
        }, {
          headers: {Authorization: this.token},
        });

        if (response.data){
          this.chat['chatName']= response.data['chatName'];
        }
      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }
    },
    respondMessageContentPrep(msgId){
      let respondMessage= this.messages.filter(message => message['msgId'] === msgId)[0];

      let resultData= null;

      if (respondMessage){
        resultData= {
          senderName: this.participantNames[respondMessage['senderId']],
          textContent: respondMessage['textContent'],
          photoContent: respondMessage['photoContent'],
        };
      }
      return resultData;
    },
    updateParticipantNamesDict(){
      this.participantNames= {};
      this.chat['participants'].forEach(participant => {
        this.participantNames[participant['usrId']]= participant['userName'];
      });
    },
    errorHandler(e){
      this.errormsg = e;
    },
  }
}
</script>

<template>
  <LoadingSpinner :loading="loading" loading-text="Caricando la chat... " />
  <div v-if="!loading" class="chat-container">
    <ErrorMsg v-if="errormsg" :msg="errormsg" @close="errormsg= null" />
    <div class="chat-info-container">
      <div class="chat-image-container">
        <img :src="'data:image/png;base64,'+ chat['chatPhoto']" alt="Chat Image" draggable="false">
      </div>
      <div class="chat-info-text-container">
        <div class="chat-name">
          <h3>{{ chat['chatName'] }}</h3>
        </div>
        <div class="participants">{{ chat['isGroup'] ? Object.values(participantNames).join(", ") : "..." }}</div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary shadow-none" @click="getMessagesSetInterval">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#rotate-cw" /></svg> Ricarica Messaggi
          </button>
          <button v-if="chat['isGroup']" type="button" class="btn btn-sm btn-outline-dark shadow-none" @click="changeGroupInfo= false; forwardMessagePanel= false; addParticipantPanel= !addParticipantPanel">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus" /></svg> Aggiungi Partecipante
          </button>
          <button v-if="chat['isGroup']" type="button" class="btn btn-sm btn-outline-dark shadow-none" @click="addParticipantPanel= false; forwardMessagePanel= false; changeGroupInfo= !changeGroupInfo; newGroupName= chat['chatName']">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#menu" /></svg> Modifica info gruppo
          </button>
          <button type="button" class="btn btn-sm btn-outline-danger shadow-none" @click="$router.replace('/chats')">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg> Chiudi
          </button>
          <button type="button" class="btn btn-sm btn-danger shadow-none" @click="leaveChat">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out" /></svg> Abbandona
          </button>
        </div>
      </div>
    </div>

    <div class="message-sender">
      <div v-if="respondTo !== -1" class="respond-message-content">
        <button v-if="respondTo !== -1" class="cancel-respond-to" type="button" @click="respondTo= -1">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
        </button>
        <RespondMsgContent v-if="respondMessageData" :key="respondMessageData['msgId'] +'-'+ Math.floor(Math.random() * 10)" :message-data="respondMessageData" />
      </div>

      <div class="message-form">
        <messageForm @prep-message="sendMessage" />
      </div>
    </div>

    <div class="messages-main">
      <div class="messages-container">
        <div class="messages">
          <ChatMessage
            v-for="message in messages"
            :key="`${message['msgId']}-${message['deliveryStatus']}`"
            :message-data="message"
            :respond-message-data="respondMessageContentPrep(message['respondTo'])"
            :sender-name="participantNames[message['senderId']]"
            :chat-is-group="chat['isGroup']"
            @respond-msg="(msgId) => respondTo= msgId"
            @forward-msg="(msgId) => forwardMsgId= msgId"
            @delete-msg="deleteMessage"
          />
        </div>
      </div>

      <transition name="add-participant-panel">
        <div v-if="addParticipantPanel" class="add-participant-panel bobby">
          <div class="select-participant">
            <sidebarList :banner-component="'userBanner'" items="users" @error="errorHandler" @banner-data="addParticipant" />
          </div>
        </div>
      </transition>

      <transition name="change-group-info-panel">
        <div v-if="changeGroupInfo" class="change-group-info-panel bobby">
          <div class="update-group-photo">
            <button class="chat-image-button" type="button" @click="imageUpload">
              <img :src="'data:image/png;base64,'+ chat['chatPhoto'] || '/images/def_group.png'" alt="Anteprima" draggable="false">
              <span>
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
              </span>
            </button>
          </div>
          <form class="update-group-name" @submit.prevent="updateGroupName">
            <span>
              <label for="new-chatname">Nome della chat:  <span v-if="!groupNameIsValid && newGroupName" style="color: red; margin-top: 10px;">Nome non valido</span></label>
              <input id="new-chatname" v-model="newGroupName" type="text" placeholder="Inserisci il nome" required>
            </span>
            <button class="new-group-name-button" type="submit" :disabled="!groupNameIsValid || newGroupName === chat['chatName']">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigation" /></svg>
            </button>
          </form>
        </div>
      </transition>

      <transition name="forward-message-panel">
        <div v-if="forwardMessagePanel" class="forward-message-panel bobby">
          <forward-message :chat-id="this.chat['chatId']" :msg-id="forwardMsgId" @close="forwardMessagePanel= false; forwardMsgId= -1" @error="errorHandler" />
        </div>
      </transition>
    </div>

  </div>
</template>

<style scoped>
.chat-container {
  width: 100%;
  height: 100%;

  margin: 0 2px 0 2px;

  display: flex;
  flex-direction: column;
}

.chat-info-container {
  height: 25%;
  width: 100%;

  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  padding: 1rem 0 0.5rem 0.2rem;

  margin-bottom: 1rem;
  border-bottom: 1px solid #dee2e6;
}

/* Barra superiore */
.chat-image-container {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  user-select: none;
}

.chat-image-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.chat-info-text-container {
  flex: 1;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-left: 10px;
  user-select: none;
}

.chat-name {
  width: fit-content;
  max-width: 200px;

  font-size: 1.1em;
  font-weight: bold;
  color: #333;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.participants {
  max-width: 180px;

  font-size: 0.9em;
  color: lightgrey;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
/* fine barra superiore */

/* invio messaggio */
.message-sender{
  height: 15%;
  width: 95%;

  display: flex;
  flex-direction: row;
  justify-content: right;


  margin: 0 5px 0 auto;
  padding: 2px 0 2px 0;
}

.message-form{
  width: 50%;
  height: 100%;
}

.respond-message-content{
  width: 50%;
  height: 100%;
  margin: 2px;

  display: flex;
  flex-direction: row;

  justify-content: center;
  align-items: center;
  
  overflow: hidden;
}

.cancel-respond-to{
  height: 30%;
  aspect-ratio: 1/1;

  margin: 0px 2px 0px 2px;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  border: 1px solid darkred;
  border-radius: 25%;
  background: red;
  color: white;

}
.cancel-respond-to:hover{
  opacity: 0.5;
  transition: opacity 0.3 ease;
}
/* fine invio messaggio */

/* gestione dei messaggi */
.messages-main {
  position: relative;
  width: 100%;
  height: 65%;
  overflow: hidden;
}

.messages-container {
  position: relative;
  height: 100%;
  width: 100%;

  overflow: hidden;
  overflow-y: auto;
}

.messages {
  position: relative;
  height: 100%;
  width: 100%;

  display: flex;
  flex-direction: column;

  margin: 0 2px 0 2px;
}
/* fine gestione messaggi */

/* aggiungere partecipante */
.add-participant-panel{
  position: absolute;
  width: 30%;
  height: 70%;

  top: 15%;
  right: 0.5%;

  background: white;
  box-shadow: -2px 0 6px rgba(0,0,0,0.2);
  z-index: 999;

  display: flex;
  flex-direction: column;
  align-items: center;
}

.select-participant {
  display: flex;
  flex-direction: column;
  align-items: center;

  width: 100%;
  height: 100%;
  padding: 4px;

  overflow: hidden;
}
/* fine aggiunta partecipante */

/* modifica info gruppo */
.change-group-info-panel{
  position: absolute;
  width: 25%;
  height: 50%;

  top: 1%;
  left: 0.5%;

  background: white;
  box-shadow: -2px 0 6px rgba(0,0,0,0.2);

  display: flex;
  flex-direction: column;

  justify-content: center;
  align-items: center;
}

.update-group-photo {
  height: 50%;
  aspect-ratio: 1/1;

  display: flex;
  justify-content: center;
  align-items: center;

  border-bottom: 1px solid black;
}

.chat-image-button {
  height: 90%;
  aspect-ratio: 1/1;
  border-radius: 50%;

  padding: 0;
  border: none;
  background: lightgray;
  position: relative;
  overflow: hidden;
}

.chat-image-button img {
  display: block;

  width: 100%;
  height: 100%;


  object-fit: cover;
  object-position: center;
  pointer-events: none;
}

.chat-image-button:not(:disabled):hover {
  filter: brightness(0.6);
  transition: filter 0.2s ease-in-out;
}

.chat-image-button span {
  position: absolute;
  width: 40%;
  height: 40%;

  border-radius: 50%;
  padding: 10px;

  top: 50%;
  left:50%;
  transform: translate(-50%, -50%);

  color: white;
  background: rgba(0, 0, 0, 0.5);
}

.chat-image-button svg {
  width: 100%;
  height: 100%;
}

.update-group-name {
  width: 90%;
  height: fit-content;

  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
}

#new-chatname{
  width: 90%;
}

.new-group-name-button {
  width: 25px;
  height: 25px;

  margin: 1px 5px 2px 0;

  border-radius: 25%;
  border: 2px dashed lightseagreen;


  color: white;
  background-color: lightseagreen;
  cursor: pointer;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
  transition: .4s;

  display: flex;
  justify-content: center;
  align-items: center;
}

.new-group-name-button:not(:disabled):hover {
  transition: .4s;
  border: 2px dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.new-group-name-button:active {
  background-color: lightseagreen;
}

.new-group-name-button:disabled{
  background-color: gray;
  border: 2px solid gray;
  cursor: default;
}
/* fine modifica info gruppo */

/* Panel per inoltrare un messaggio */
.forward-message-panel{
  position: absolute;
  width: 30%;
  height: 80%;

  top: 10%;
  right: 40%;

  background: white;
  box-shadow: -2px 0 6px rgba(0,0,0,0.2);
  z-index: 999;

  display: flex;
  flex-direction: column;
  align-items: center;
}
/* Fine panel per inoltrare un messaggio */

/* Transition per il pannello dei partecipanti */
.add-participant-panel-enter-from,
.add-participant-panel-leave-to {
  transform: translateX(100%); /* fuori dallo schermo a destra */
}

.add-participant-panel-enter-to,
.add-participant-panel-leave-from {
  transform: translateX(0); /* posizione normale, visibile */
}

.add-participant-panel-enter-active,
.add-participant-panel-leave-active {
  transition: transform 0.3s ease;
}
/* Fine transition per il pannello dei partecipanti */

/* Transition per il pannello delle info del gruppo */
.change-group-info-panel-enter-from,
.change-group-info-panel-leave-to {
  transform: translateY(-100%); /* fuori dallo schermo a destra */
}

.change-group-info-panel-enter-to,
.change-group-info-panel-leave-from {
  transform: translateY(0); /* posizione normale, visibile */
}

.change-group-info-panel-enter-active,
.change-group-info-panel-leave-active {
  transition: transform 0.3s ease;
}
/* Fine transition per il pannello delle info del gruppo */

/* Transition per il pannello che inoltra un messaggio */
.forward-message-panel-enter-from,
.forward-message-panel-leave-to {
  transform: translateY(100%); /* fuori dallo schermo a destra */
}

.forward-message-panel-enter-to,
.forward-message-panel-leave-from {
  transform: translateY(0); /* posizione normale, visibile */
}

.forward-message-panel-enter-active,
.forward-message-panel-leave-active {
  transition: transform 0.3s ease;
}
/* Fine transition per il pannello che inoltra un messaggio */
</style>
