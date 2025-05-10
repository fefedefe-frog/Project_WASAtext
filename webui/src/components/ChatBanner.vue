<script>
export default {
  props: {
    chatData: {
      type: Object,
      required: true
    },
  },
  emits: [
    'error', 'chatBannerData'
  ],
  data() {
    return {
      errormsg: null,
      token: "",
      chatId: -1,
      participantNames: {},
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
  mounted() {
    this.chatId= this.chatData['chatId'];
    this.token= sessionStorage.getItem('authToken');

    this.chatData['participants'].forEach(participant => {
      this.participantNames[participant['usrId']]= participant['userName'];
    });

    this.getChatMessages();
    this.setIntervalId= setInterval(async () => {
      await this.getChatMessages();
    }, 7000)
  },
  beforeUnmount() {
    clearInterval(this.setIntervalId);
  },
  deactivated() {
    clearInterval(this.setIntervalId);
  },
  computed: {
    messagePreview() {
      let message= this.lastMessage

      let messagePreview= `${this.participantNames[message['senderId']]}:"`;
      if (message['textContent'] !== ""){
        if (message['textContent'].length > 10){
          messagePreview= `${messagePreview} ${message['textContent'].substring(0, 10)}...`;
        }
        messagePreview= `${messagePreview} ${message['textContent']}`;
      }
      return messagePreview
    }
  },
  methods: {
    async getChatMessages(){
      console.log(this.chatData)
      this.errormsg= null;

      try {
        let response= await this.$axios.put(`/chats/${this.chatId}/messages`, {
          //msgId: this.lastMsgId
          msgId: -1
        }, {
          headers: {Authorization: this.token},
        });

        if (response.data) {
          if (Array.isArray(response.data['messages']) && response.data['messages'].length > 0){
            response.data['messages'].forEach(message => {
              this.messages.push(message);
            })
            this.lastMessage= this.messages[this.messages.length-1];
            //this.lastMsgId= this.lastMessage['msgId'];
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
    <!-- Immagine chat a destra -->
    <div class="chat-image-container">
      <img :src="'data:image/png;base64,'+chatData['chatPhoto']" alt="Chat Image">
    </div>

    <!-- Nome della chat e ultimo messaggio -->
    <div class="chat-text-container">
      <span class="chat-name"> {{ chatData['chatName'] }} </span>
      <span class="last-message-text">{{this.messagePreview}}</span><svg v-if="lastMessage['photoContent'].length > 0" class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
    </div>
  </div>
</template>

<style scoped>

.chat-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;
  margin-bottom: 5px;
  background-color: lightgray;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  height: 70px;
  width: 100%;
}

.chat-banner:hover {
  background-color: darkgray;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.4);
}

.chat-image-container {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.chat-image-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.chat-text-container {
  height: 100%;

  flex: 1;
  padding: 0 0 0 5px;

  display: flex;
  flex-direction: column;
  justify-content: center;

  margin-left: 10px;
  user-select: none;
}

.chat-name {
  height: fit-content;
  user-select: none;

  color: black;
  font-weight: bold;
  font-size: 1rem;
}

.last-message-text {
  user-select: none;
  height: fit-content;

  margin-left: 10px;
  color: dimgray;
  font-size: 0.9em;
  font-weight: normal;
}
</style>