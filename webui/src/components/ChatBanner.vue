<script>
export default {
  props: {
    chatData: {
      type: Object,
      required: true
    },
  },
  data() {
    return {
      errormsg: null,
      token: "",
      chatId: -1,
      messages: [],

      lastMessage: {
        msgId: 3,
        senderId: "",
        contentType: "",
        content: "",
        deliveryStatus: "",
        timestamp: "",
        comments: [],
      },
      lastMsgId: -1,

      setIntervalId: null
    }
  },
  computed: {
    messageContent(){
      if (this.lastMessage){
        let content= this.lastMessage['content'];
        let sender= this.lastMessage['senderId'];
        let messsagePreview= sender + ": " + content;

        if (this.lastMessage['contentType'] == "photo"){
          return `<span class="last-message-text">${sender}: </span><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>`;
        }else {
          return `<span class="last-message-text">${messsagePreview.length > (10 + sender.length) ? messsagePreview.substring(0, (10 + sender.length))+ "..." : messsagePreview}</span>`;
        }
      }else {
        return '';
      }
    }
  },
  mounted() {
    this.chatId= this.chatData['chatId'];
    this.token= sessionStorage.getItem('authToken');

    this.getChatMessages();
    this.setIntervalId= setInterval(async () => {
      this.getChatMessages();
    }, 7000)
  },
  beforeUnmount() {
    clearInterval(this.setIntervalId);
  },
  deactivated() {
    clearInterval(this.setIntervalId);
  },
  methods: {
    async getChatMessages(){
      this.errormsg= null;

      try {
        let response= await this.$axios.put(`/chats/${this.chatId}/messages`, {
          msgId: this.lastMsgId
        }, {
          headers: {Authorization: this.token},
        });

        if (response.data) {
          if (Array.isArray(response.data['messages']) && response.data['messages'].length > 0){
            response.data['messages'].forEach(message => {
              this.messages.push(message);
            })
            this.lastMessage= this.messages[this.messages.length-1];
            this.lastMsgId= this.lastMessage['msgId'];
          }
        }
      }catch(e) {
        this.$emit('error', e)
      }
    },
    bannerClicked(){
      this.$emit('chatBannerData', {chatData: this.chatData, messages: this.messages})
    }
  },
};
</script>

<template>
  <div class="chat-banner" @click="bannerClicked">
    <!-- Foto Profilo a destra -->
    <div class="chat-image-container">
      <img :src="'data:image/png;base64,'+chatData['chatPhoto']" alt="Chat Image">
    </div>

    <!-- Contenuto del banner -->
    <div class="text-container">
      <div class="chat-name">{{ chatData['chatName'] }}</div>
      <div class="last-message">
        <div v-html="messageContent" />
      </div>
    </div>
  </div>
</template>

<style scoped>

.chat-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;
  margin: 5px;
  background-color: lightgray;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  height: 70px;
  width: 200px;
}

.chat-banner:hover {
  background-color: darkgray;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.4);
}

.chat-image-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.chat-image-container {
  max-width: 100%;  /* L'immagine non andrà mai oltre la larghezza del suo contenitore */
  max-height: 100%; /* L'immagine non andrà mai oltre l'altezza del suo contenitore */
  object-fit: contain; /* L'immagine si adatta dentro il box senza distorsioni */
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

.last-message {
  font-size: 0.9em;
  color: #777;
}
</style>