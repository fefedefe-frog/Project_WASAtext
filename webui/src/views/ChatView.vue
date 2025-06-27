<script>
export default {
  data: function () {
    return {
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

      getChatInfoIntervalId: null,
      getMessagesIntervalId: null,
      getMessagesIsRunning: false,

      addParticipantPanel: false,
    }
  },
  async mounted(){
    this.chat['chatId']= this.$route.params.chat_id;
    this.token= sessionStorage.getItem('authToken');

    this.loading= true;
    await this.getChatInfo();
    this.updateParticipantNamesDict();
    

    this.getChatInfoIntervalId= setInterval( async () => {
      await this.getChatInfo();
      this.updateParticipantNamesDict();
    }, 30000);

    // Funzione per avviare un nuovo intervallo di ricezione dei messaggi
    this.getMessagesSetInterval();
  },
  beforeUnmount() {
    clearInterval(this.getChatInfoIntervalId);
    clearInterval(this.getMessagesIntervalId);
  },
  watch: {
    '$route.params.chat_id': {
      immediate: true,
      async handler(newChatId, oldChatId){
        if (newChatId !== oldChatId){
          if (this.token === "") this.token= sessionStorage.getItem('authToken');

          this.loading= true;
          this.chat['chatId']= newChatId;

          await this.getChatInfo();
          this.updateParticipantNamesDict();
        }
      }
    }
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
      } catch (e) {
        this.errormsg= e.toString();
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
        this.errormsg= e.toString();
      }
      this.getMessagesIsRunning= false;
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
        this.errormsg= e.toString();
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
      }catch (e){
        this.errormsg= e.toString();
      }finally {
        this.respondTo= -1;
      }

    },
    async leaveChat() {
      try {
        let response= await this.$axios.delete(`/chats/${this.chat['chatId']}/users`, {
          headers: {Authorization: this.token}
        });
      }catch(e) {
        this.errormsg= e.toString();
      }finally {
        this.$router.replace('/chats');
      }
    },
    async addParticipant(participant){
      if (participant['usrId'] in this.participantNames){
        this.errormsg= new Error("Non puoi aggiungere un utente già partecipante").toString();
      }else {
        this.errormsg= null

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
          this.errormsg= e.toString();
        }
      }


    },
    updateParticipantNamesDict(){
      this.participantNames= {};
      this.chat['participants'].forEach(participant => {
        this.participantNames[participant['usrId']]= participant['userName'];
      });
    },
    errorHandler(e){
      this.errormsg = e.toString();
    },
  }
}
</script>

<template>
  <LoadingSpinner :loading="loading" loading-text="Caricando la chat... " />
  <div v-if="!loading" class="chat-container">
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
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
          <button v-if="chat['isGroup']" type="button" class="btn btn-sm btn-outline-dark shadow-none" @click="addParticipantPanel= !addParticipantPanel">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus" /></svg> Aggiungi Partecipante
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
      <messageForm @prep-message="sendMessage" />
    </div>

    <div class="messages-main">
      <div class="messages-container">
        <ChatMessage v-for="message in messages" :key="`${message['msgId']}-${message['deliveryStatus']}`" :message-data="message" :sender-name="participantNames[message['senderId']]" :chat-is-group="chat['isGroup']" />
      </div>
      <transition name="add-participant-panel">
        <div class="add-participant-panel" v-if="addParticipantPanel">
          <div class="select-participant">
            <sidebarList :banner-component="'userBanner'" items="users" @error="errorHandler" @banner-data="addParticipant" />
          </div>
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
  height: 60px;
  max-height: 100px;
  width: 50%;

  margin: 0 5px 0 auto;
}
/* fine invio messaggio */

/* gestione dei messaggi */
.messages-main {
  position: relative;
  width: 100%;
  height: 100%;

  overflow: hidden;
  overflow-y: auto;
}

.messages-container {
  height: fit-content;
  width: 100%;

  display: flex;
  flex-direction: column;

  margin: 0 2px 0 2px;
}
/* fine gestione messaggi */

/* aggiungere partecipante */
.add-participant-panel{
  position: absolute;
  top: 0;
  right: 0;
  width: 30%;
  height: 50%;
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
  padding: 4px;

  overflow: hidden;
}
/* fine aggiunta partecipante */

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
</style>
