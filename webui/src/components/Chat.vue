<script>
export default {
  props: {
    chatData: {
      type: Object,
      required: true
    },
    initialMessages: {
      type: Array,
      required: true
    },
  },
  emits: [
      'closeChat'
  ],
  data: function () {
    return {
      token: '',
      errormsg: null,

      chat: {
        chatId: -1,
        chatName: "",
        chatPhoto: "",
        isGroup: false,
        participants: []
      },
      participantsNames: {},

      lastMsgId: -1,
      messages: [],
      textContent: '',

      getChatInfoIntervalId: null,
      getMessagesIntervalId: null,
      getMessagesIsRunning: false,
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.chat['chatId']= this.chatData['chatId'];
    this.chat['chatName']= this.chatData['chatName'];
    this.chat['chatPhoto']= this.chatData['chatPhoto'];
    this.chat['isGroup']= this.chatData['isGroup'];
    this.chat['participants']= this.chatData['participants'];

    this.lastMsgId= this.initialMessages[this.initialMessages.length - 1]['msgId'];
    this.messages= this.initialMessages;

    this.getChatInfo();
    this.getParticipants();
    this.getMessages();

    this.getChatInfoIntervalId= setInterval( async () => {
      await this.getChatInfo();
      await this.getParticipants();
    }, 30000);

    this.getMessagesIntervalId= setInterval( async () => {
      this.lastMsgId= -1; //Per ora recupero tutti i messagi della chat in ogni caso
      await this.getMessages();
      await this.updateReadStatus();
    }, 5000);

  },
  beforeUnmount() {
    clearInterval(this.getChatInfoIntervalId);
    clearInterval(this.getMessagesIntervalId);
  },
  deactivated() {
    clearInterval(this.getChatInfoIntervalId);
    clearInterval(this.getMessagesIntervalId);
  },
  methods: {
    async getChatInfo() {
      this.errormsg= null;

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
        this.errormsg = e;
      }
    },
    async getMessages() {
      // Evito la sovrapposizione della funzione quando chiamata dagli intervalli
      if (this.getMessagesIsRunning) return;
      this.getMessagesIsRunning= true;

      this.errormsg = null;

      try {
        let response= await this.$axios.put(`/chats/${this.chat['chatId']}/messages`, {
          msgId: this.lastMsgId
        }, {
          headers: {Authorization: this.token}
        });

        if (response.data) {
          //if (this.lastMsgId === -1){
          //  this.messages= [];
          //}
          if (Array.isArray(response.data['messages']) && response.data['messages'].length > 0){
            this.messages= [];
            response.data['messages'].forEach(message => {
              this.messages.push(message);
            });

            //this.lastMsgId= this.messages[this.messages.length-1]['msgId'];
            this.lastMsgId= -1;
          }
        }
      }catch(e) {
        this.errormsg = e;
      }
      this.getMessagesIsRunning= false;
    },
    async updateReadStatus() {
      this.errormsg = null;

      try {
        let response= await this.$axios.put(`/chats/${this.chat['chatId']}/messages/${this.lastMsgId}`, {}, {
          headers: {Authorization: this.token}
        });
        if (response.status > 400) {
          throw new Error("unable to update the messages status");
        }
      }catch(e) {
        this.errormsg = e;
      }
      this.getMessagesIsRunning= false;

    },
    async getParticipants() {
      this.errormsg = null

      try {
        let response= await this.$axios.get(`/chats/${this.chat['chatId']}/users`, {
          headers: {Authorization: this.token}
        });

        if (response.data) {
          response.data['participants'].forEach(participant => {
            this.participantsNames[participant['usrId']]= participant['userName'];
          });
        }
      }catch(e) {
        this.errormsg = e.toString();
      }
    },
    async sendMessage(){
      // TODO
    },
    closeChat(leaveChat){
      this.$emit('closeChat', leaveChat);
    },
    errorHandler(e){
      this.errormsg = e;
    },
  }
}
</script>

<template>
  <div class="w-100 h-100 p-2">
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div class="chat-image-container">
        <img :src="'data:image/png;base64,'+ chatData['chatPhoto']" alt="Chat Image">
      </div>
      <div class="text-container">
        <div class="chat-name">
          <h3>{{ chatData['chatName'] }}</h3>
        </div>
        <div class="participants">{{ chatData['isGroup'] ? Object.values(participantsNames).join(", ") : "..." }}</div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary shadow-none" @click="getMessages(-1)">
            Ricarica Messaggi
          </button>
          <button type="button" class="btn btn-sm btn-outline-dark shadow-none" @click="getMessages">
            Info
          </button>
          <button type="button" class="btn btn-sm btn-outline-danger shadow-none" @click="closeChat(false)">
            Chiudi
          </button>
          <button type="button" class="btn btn-sm btn-danger shadow-none" @click="closeChat(true)">
            Abbandona
          </button>
        </div>
      </div>
    </div>

    <div class="message-sender">
      <form class="sendMessage-form" @submit.prevent="sendMessage">
        <label for="text">
          <input id="textContent" v-model="textContent" type="text" placeholder="Scrivi un messaggio" required>
        </label>
        <button type="submit" :disabled="!textContent" :class="{ disabled: !textContent}">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigator" /></svg>
        </button>
      </form>
    </div>

    <div class="messages-main">
      <div class="messages-container">
        <ChatMessage v-for="message in messages" :key="message.msgId" :message="message" :sender-name="participantsNames[message['senderId']]" />
      </div>
    </div>
  </div>
</template>

<style scoped>

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

.text-container {
  flex: 1;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-left: 10px;
  user-select: none;
}

.chat-name {
  font-size: 1.1em;
  font-weight: bold;
  color: #333;
}

.participants {
  font-size: 0.9em;
  color: lightgrey;
}
/* fine barra superiore */

/* invio messaggio */
.message-sender{
  position: sticky;
  z-index: 1001;
  background: rgba(0, 0, 0, 0.5);
  height: 50px;
  width: 100%;
}

.message-sender svg{
  background-color: white;
  border: 1px black solid;
  border-radius: 50%;
  fill: currentColor;

  user-select: none;
}
/* fine invio messaggio */

/* gestione dei messaggi */
.messages-main {
  position: relative;
  width: 100%;
  height: 70vh;
  overflow-y: scroll;
  overflow-x: hidden;
}

.messages-container {
  position: relative;
  display: flex;
  flex-direction: column;
  width: auto;
  min-height: 100%;
  margin-left: 2px;
  margin-right: 2px;
}
/* fine gestione messaggi */
</style>
