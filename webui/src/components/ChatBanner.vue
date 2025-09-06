<script>
export default {
  props: {
    inputData: {
      type: Object,
      required: true
    },
  },
  emits: [
    'error', 'bannerClicked',
  ],
  data() {
    return {
      token: "",
      usrId: "",

      chat: {},
      participantNames: {},

      lastMessage: {},
      lastMessPreview: "",
      lastMessTimestamp: ""
    }
  },
  created() {
    this.chat= this.inputData['chat'];
    this.lastMessage= this.inputData['lastMsg'];

    this.chat['participants'].forEach(p =>{
      this.participantNames[p['usrId']]= p['userName'];
    });

  },
  mounted(){
    this.chat= this.inputData['chat'];
    this.lastMessage= this.inputData['lastMsg'];

    this.chat['participants'].forEach(p =>{
      this.participantNames[p['usrId']]= p['userName'];
    });

    this.makeMessPreview();
  },
  methods: {
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

      //Timestamp
      if (this.lastMessage['timestamp']){
        const dateObject= new Date(this.lastMessage['timestamp']);
        const currDate= new Date();

        // Formatto giorno mese anno dal timestamp
        let dateFormatter = new Intl.DateTimeFormat('it-IT', { dateStyle: 'short' });
        let formattedDate = dateFormatter.format(dateObject);

        // Formatto ore e minuti
        let timeFormatter = new Intl.DateTimeFormat('it-IT', { hour: '2-digit', minute: '2-digit' });
        let formattedTime = timeFormatter.format(dateObject);

        // calcolo la differenza tra le date, se maggiore di un giorno mostro il giorno del messaggio, se minore mostro la data
        let giornoInMs= 24 * 60 * 60 * 1000;

        if ((currDate.getTime() - dateObject.getTime()) > giornoInMs ){
          this.lastMessTimestamp= formattedDate;
        }else {
          this.lastMessTimestamp= formattedTime;
        }
      }
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
        <span class="last-message-timestamp">{{ lastMessTimestamp }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>

.chat-banner {
  height: 50px;
  width: 100%;

  padding: 5px;
  margin-bottom: 5px;

  display: flex;
  justify-content: space-between;
  align-items: center;

  background-color: lightgray;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
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
  height: fit-content;
  width: fit-content;
  max-width: 90%;
  margin: 0 2px 0 10px;

  color: dimgray;
  font-size: 0.9em;
  font-weight: normal;
  user-select: none;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.last-message-timestamp{
  height: fit-content;
  width: fit-content;
  margin-right: 2%;

  font-style: oblique;
  margin-left: auto;
  font-size: 0.9em;

  user-select: none;
}
</style>