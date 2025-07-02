<script>
export default {
  props: {
    inputData: {
      type: Object,
      required: true
    },
  },
  emits: [
    'error', 'bannerClicked'
  ],
  data() {
    return {
      errormsg: null,
      token: "",
      usrId: "",


      chat: {
        chatId: -1,
        chatName: "",
        chatPhoto: "",
        isGroup: false,
        participants: []
      },
      participantNames: {},

      lastMessage: {
        msgId: -1,
        senderId: "",
        respondTo: -1,
        textContent: "",
        photoContent: "",
        timestamp: "",
        comments: [],
      },
      lastMessPreview: "",
      lastMsgId: -1
    }
  },
  created() {
    this.chat= this.inputData;
    this.token= sessionStorage.getItem('authToken');
    this.usrId= sessionStorage.getItem('usrId')

    this.chat['participants'].forEach(p =>{
      this.participantNames[p['usrId']]= p['userName'];
    });

  },
  async mounted(){
    await this.getChatMessages();
    this.makeMessPreview();
  },
  methods: {
    async getChatMessages(){
      try {
        let response= await this.$axios.put(`/chats/${this.chat['chatId']}/messages`, {
          msgId: this.lastMsgId
        }, {
          headers: {Authorization: this.token},
        });

        if (response.data) {
          if (Array.isArray(response.data['messages']) && response.data['messages'].length > 0){
            this.messages= [];
            response.data['messages'].forEach(message => {
              this.messages.push(message);
            })
            this.lastMessage= this.messages[this.messages.length-1];
            this.lastMsgId= this.lastMessage['msgId'];
          }
        }
      }catch(e) {
        this.$emit('error', e);
      }
    },
    bannerClicked(){
      this.$emit('bannerClicked', this.chat['chatId'])
    },
    makeMessPreview() {
      let preview= "";

      if(this.lastMessage['senderId'] === this.usrId){
        preview= `tu:`;
      }else if (this.chat['isGroup']){
        preview= `${this.participantNames[this.lastMessage['senderId']]}:`;
      }
      if (this.lastMessage['textContent'] !== ""){
        preview= `${preview} ${this.lastMessage['textContent']}`;
      }
      this.lastMessPreview= preview;
    }
  },
};
</script>

<template>
  <div class="chat-banner" @click="bannerClicked">
    <!-- Immagine chat a destra -->
    <div class="chat-image-container">
      <img :src="'data:image/png;base64,'+chat['chatPhoto']" alt="Chat Image" draggable="false">
    </div>

    <!-- Nome della chat e ultimo messaggio -->
    <div class="chat-text-container">
      <span class="chat-name"> {{ chat['chatName'] }} </span>
      <div class="message-preview">
        <span class="last-message-text">{{ lastMessPreview }}</span>
        <svg v-if="lastMessage['photoContent'].length > 0" class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
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
  width: 70%;

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

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.message-preview {
  width: 100%;
  display: flex;
  flex-direction: row;

  align-items: center;
}

.last-message-text {
  user-select: none;
  height: fit-content;
  width: fit-content;
  max-width: 90%;


  margin: 0 2px 0 10px;

  color: dimgray;
  font-size: 0.9em;
  font-weight: normal;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

</style>