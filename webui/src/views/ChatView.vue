<script>
export default {
  data: function () {
    return {
      token: '',
      errormsg: null,
      loading: false,
      loadingMessages: false,
      isChatInfoVisible: false,
      chatInfo: {
        chatId: -1,
        chatName: '',
        chatPhoto: '',
        isGroup: false,
        participants: [], // can be a list of username, or usrId if the request for their name doesn't go as expected
      },
      participantsInfo: [],
      participantsName: [],
      messages: [],
      textContent: '',
    }
  },
  mounted() {
    this.chatInfo['chatId']= this.$route.params.chat_id;
    this.token= localStorage.getItem('authToken');
    this.getChatInfo().then(
      this.getMessages()
    )
    this.getParticipantsInfo()
  },
  created() {
    this.chatInfo['chatId']= this.$route.params.chat_id;
    this.token= localStorage.getItem('authToken');
  },
  methods: {
    async getChatInfo() {
      this.loading= true
      this.errormsg= null
      this.messages=[]

      try {
        let response = await this.$axios.get(`/chats/${this.chatInfo['chatId']}`, {headers: {Authorization: `${this.token}`}});
        if (response.data) {
          this.chatInfo['chatName']= response.data.chatName
          this.chatInfo['chatPhoto']= response.data.chatPhoto
          this.chatInfo['isGroup']= response.data.isGroup
          this.chatInfo['participants']= response.data.participants
        }
      } catch (e) {
        this.errormsg = e.toString();
      }finally {
        this.loading = false;
      }
    },
    async getMessages() {
      this.loadingMessages = true
      this.errormsg = null

      try {
        let response= await this.$axios.put(`/chats/${this.chatInfo['chatId']}/messages`, {}, {headers: {Authorization: this.token}});
        if (response.data) {
          this.messages= []
          response.data['messages'].forEach(message => {
            this.messages.push(message)
          })
        }
      }catch(e) {
        this.errormsg = e.toString();
      }finally {
        this.loadingMessages = false
      }
    },
    async getParticipantsInfo() {
      this.participantsInfo = []
      this.errormsg = null

      try{
        let response= await this.$axios.get(`/chats/${this.chatInfo['chatId']}/users`, {headers: {Authorization: `${this.token}`}});
        if (response.data) {
          response.data['participants'].forEach(participant => {
            this.participantsInfo.push(participant)
            this.participantsName.push(participant['userName'])
          })
        }
      }catch(e) {
        this.errormsg = e.toString();
      }
    },
    async leaveChat() {
      console.log('leave chat')
    },
    toggleChatInfo() {
      this.isChatInfoVisible = !this.isChatInfoVisible
    },
    errorHandler(e){
      this.errormsg = e.toString();
    },
  }
}
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div class="chat-image-container">
        <img :src="'data:image/png;base64,'+ chatInfo['chatPhoto']" alt="Chat Image">
      </div>
      <div class="text-container">
        <div class="chat-name">
          <h3>{{ chatInfo['chatName'] }}</h3>
        </div>
        <div v-if="chatInfo['isGroup']" class="participants">{{ participantsName.length>0 ? participantsName.join(", ") : chatInfo['participants'].join(", ") }}</div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="getMessages()">
            Refresh messages
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="toggleChatInfo">
            {{ isChatInfoVisible ? "Close info" : "Chat Info" }}
          </button>
          <button type="button" class="btn btn-sm btn-outline-primary" @click="leaveChat">
            Leave chat
          </button>
        </div>
      </div>
    </div>

    <div class="message-sender">
      <form class="sendMessage-form" @submit.prevent="sendMessage">
        <label for="text">
          <input id="textContent" v-model="textContent" type="text" placeholder="Scrivi un messaggio" required>
        </label>
        <button type="submit" :disabled="!textContent || loading" :class="{ disabled: !textContent}">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigator" /></svg>
        </button>
      </form>
      <ChatInfo :is-visible="isChatInfoVisible" :chat-id="chatInfo['chatId']" :participants-info="participantsInfo" @error="errorHandler" @visibility="toggleChatInfo" />
    </div>

    <div class="messages-main">
      <div class="messages-container">
        <ChatMessage v-for="message in messages" :key="message.msgId" :message="message" />
      </div>
    </div>


    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
/* Barra superiore */
.chat-image-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.chat-image-container {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
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
  color: #777;
}
/* fine barra superiore */

.messages-main {
  position: relative;
  width: 100%;
  height: 70vh;
  overflow-y: scroll;
  overflow-x: hidden;
}

/* gestione dei messaggi */
.messages-container {
  position: relative;
  display: flex;
  flex-direction: column;
  width: auto;
  min-height: 100%;
  margin-left: 2px;
  margin-right: 2px;
}

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
</style>
