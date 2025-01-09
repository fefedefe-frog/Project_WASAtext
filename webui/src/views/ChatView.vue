<script>
import ChatMessage from "../components/ChatMessage.vue";

export default {
  components: {ChatMessage},
  data: function () {
    return {
      token: '',
      errormsg: null,
      loading: false,
      loadingMessages: false,
      isChatInfoVisible: false,
      chatId: -1,
      chatInfo: {
        chatName: '',
        chatPhoto: '',
        isGroup: false,
        participants: [],
      },
      messages: [],
    }
  },
  methods: {
    async refreshChatInfo() {
      this.loading= true
      this.errormsg= null

      try {
        let response = await this.$axios.get('/chats/'+ this.chatId, {headers: {Authorization: `${this.token}`}});
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
        let response= await this.$axios.put('/chats/'+ this.chatId+ '/messages', {}, {headers: {Authorization: this.token}});
        if (response.data) {
          this.messages= []
          response.data['messages'].forEach((message) => {
            this.messages.push(message)
          })
        }
      }catch(e) {
        this.errormsg = e.toString();
      }finally {
        this.loadingMessages = false
      }
    },
    async leaveChat() {
      console.log('leave chat')
    },
    showChatInfo() {
      console.log('showChatInfo')
      this.isChatInfoVisible = true
    },
    hideChatInfo() {
      console.log('hideChatInfo')
      this.isChatInfoVisible = false
    }
  },
  mounted() {
    this.chatId= this.$route.params.chat_id;
    this.token= localStorage.getItem('authToken');
    this.refreshChatInfo().then(
      this.getMessages()
    )
  },
  created() {
    this.chatId= this.$route.params.chat_id;
    this.token= localStorage.getItem('authToken');
  }
}
</script>

<template>
  <div class="main-container">
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div class="chat-image-container">
        <img :src="'data:image/png;base64,'+ chatInfo['chatPhoto']" alt="Chat Image" />
      </div>
      <div class="text-container">
        <div class="chat-name">{{ chatInfo['chatName'] }}</div>
        <div v-if="chatInfo['isGroup']" class="participants">{{ chatInfo['participants'].join(", ") }}</div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="getMessages()">
            Refresh messages
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="showChatInfo">
            Chat Info
          </button>
          <button type="button" class="btn btn-sm btn-outline-primary" @click="leaveChat">
            Leave chat
          </button>
        </div>
      </div>
    </div>
    <div class="messages-container">
      <div v-for="message in messages" class="message-div">
        <ChatMessage :message="message"></ChatMessage>
      </div>


      <LoadingSpinner v-if="loadingMessages" :loading="loadingMessages" loadingText="Caricamento Messaggi"/><LoadingSpinner/>

      <div v-if="isChatInfoVisible" class="chat-info-sidebar">
        <div class="chat-info-sidebar-header">
          <h3>Chat Info</h3>
          <button @click="hideChatInfo" class="close-chat-info-sidebar-btn">Chiudi</button>
        </div>
        <div class="chat-info-sidebar-content">
          <p>Contenuto del componente.</p>
          <p>Aggiungi tutto quello che vuoi qui.</p>
        </div>
      </div>

      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
  </div>
</template>

<style scoped>
.main-container {
  border: blue 2px solid;
  height: 100%;
}

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


/* gestione dei messaggi */
.messages-container {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: 100%;
  border-bottom: saddlebrown 1px solid;
}

.message-div {
  display: flex;
  flex-direction: row;
  width: 100%;
  border: black 1px solid;
}

.message-div ChatMessage {
  margin-left: auto;
  border: black 1px solid;
}

/* sidebar chat info*/
.chat-info-sidebar {
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  height: 400px;
  background-color: #f8f9fa;
  box-shadow: -2px 0 5px rgba(0, 0, 0, 0.5);
  z-index: 1000;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.chat-info-sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ddd;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.close-chat-info-sidebar-btn {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 5px 10px;
  cursor: pointer;
  border-radius: 4px;
}

.close-chat-info-sidebar-btn:hover {
  background-color: #c82333;
}

.chat-info-sidebar-content {
  flex-grow: 1;
  overflow-y: auto;
  color: #333;
  font-size: 14px;
}
/* fine sidebare chat info */

</style>
