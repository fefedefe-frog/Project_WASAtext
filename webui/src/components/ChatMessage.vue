<script>
export default {
  props: {
    messageData: {
      type: Object,
      required: true
    },
    respondMessageData: {
      type: Object,
      required: false,
    },
    senderName: {
      type: String,
      required: true
    },
    chatIsGroup: {
      type: Boolean,
      required: true
    }
  },
  emits: ['respondMsg', 'forwardMsg', 'deleteMsg'],
  data: function () {
    return{
      status: "minus",
      date: "",
      textContent: "",
      photoContent: {},
      usrId: "",

      message: {
        msgId: -1,
        senderId: "",
        respondTo: -1,
        textContent: "",
        photoContent: [],
        timestamp: "",
        comments: [],
      },

      respondToData: {}
    }
  },
  computed: {
    dynamicMessageSide(){
      return{
        'justify-content': this.usrId === this.message['senderId'] ? 'flex-end' : 'flex-start',
      }
    }
  },
  mounted () {
    this.usrId= sessionStorage.getItem('usrId');
    this.message= this.messageData;

    if (this.message['respondTo'] !== -1 && this.respondMessageData){
      this.respondToData= this.respondMessageData;
    }
    this.prepMessage();
  },
  beforeUnmount() {
    this.message= {};
  },
  methods: {
    prepMessage() {
      switch (this.message['deliveryStatus']) {
        case 'sent':
          this.status = 'clock';
          break;
        case 'received':
          this.status = 'chevron-down';
          break;
        case 'read':
          this.status = 'chevrons-down';
          break;
        default:
          this.status = 'clock';
      }


      if (this.message['timestamp']){
        const dateObject= new Date(this.message['timestamp']);

        // Formatto giorno mese anno dal timestamp
        let dateFormatter = new Intl.DateTimeFormat('it-IT', { dateStyle: 'short' });
        let formattedDate = dateFormatter.format(dateObject);

        // Formatto ore e minuti
        let timeFormatter = new Intl.DateTimeFormat('it-IT', { hour: '2-digit', minute: '2-digit' });
        let formattedTime = timeFormatter.format(dateObject);

        // Unisco nel formato che mi interessa
        this.date= formattedDate +" "+ formattedTime;
      }
    },
  }
}
</script>

<template>
  <div class="message-div" :style="dynamicMessageSide">
    <MessageDropdownMenu
      v-if="usrId === message['senderId']"
      :message-id="message['msgId']"
      :sender-id="message['senderId']"
      @respond-to="$emit('respondMsg', message['msgId'])"
      @forward-msg="$emit('forwardMsg', message['msgId'])"
      @comment-msg="console.log('TODO')"
      @delete-msg="$emit('deleteMsg', message['msgId'])"
    />
    <div class="message-container">
      <div v-if="message['respondTo'] !== -1" class="respond-message">
        <RespondMsgContent v-if="respondToData" :sender-name="respondToData['senderName']" :message-data="respondToData" />
      </div>
      <div class="message-content">
        <div v-if="message['photoContent'].length === 0" class="message-content-text-container">
          <p>{{ message["textContent"] }}</p>
        </div>
        <div v-if="message['photoContent'].length > 0" class="message-content-image-container">
          <img :src="'data:image/png;base64,'+message['photoContent']" alt="Image" draggable="false">
        </div>
      </div>
      <div class="message-info">
        <div v-if="chatIsGroup && message['senderId'] !== usrId" class="sender-name">{{ senderName }} </div>
        <div class="message-status">
          <span class="timestamp">{{ date }}</span>
          <svg v-if="message['senderId'] === usrId" class="feather delivery-status"><use :href="'/feather-sprite-v4.29.0.svg#'+ status" /></svg>
        </div>
      </div>
    </div>
    <MessageDropdownMenu
      v-if="usrId === message['senderId']"
      :message-id="message['msgId']"
      :sender-id="message['senderId']"
      @respond-to="$emit('respondMsg', message['msgId'])"
      @forward-msg="$emit('forwardMsg', message['msgId'])"
      @comment-msg="console.log('TODO')"
      @delete-msg="$emit('deleteMsg', message['msgId'])"
    />
  </div>
</template>

<style scoped>
.message-div {
  position: relative;

  display: flex;
  flex-direction: row;

  width: 100%;

  box-sizing: border-box;
  padding: 0 2px 0 2px;
}

.message-container {
  height: fit-content;
  width: fit-content;
  max-width: 50%;

  border-radius: 8px;
  border: 1px deepskyblue solid;

  padding: 2px;
  margin: 5px;

  display: flex;
  flex-direction: column;

  justify-content: flex-start;
  align-items: flex-start;

  background-color: lightskyblue;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.respond-message{
  width: 96%;
  height: 10%;

  margin: 2px 2% 2px 2%;
  padding-bottom: 2px;
  border-bottom: 1px black solid;
}

/* container info stato messaggio */
.message-info {
  width: 100%;
  height: 10%;

  margin-top: 2px;

  display: flex;
  flex-direction: row;
  align-items: center;

  user-select: none;
  border-top: 1px black solid;
}

.sender-name {
  margin-left: 2px;
  margin-right: 10px;
  font-size: 1.1em;
  font-style: italic;
  font-weight: bold;
  align-self: flex-end;
  color: #333;
}

.message-status {
  display: flex;
  flex-direction: row;
  margin-left: auto;
  align-items: center;
  width: fit-content;
}

.timestamp {
  margin-left: auto;
  margin-right: 2px;
  user-select: none;
}

.delivery-status{
  border: 1px black solid;
  border-radius: 50%;
  margin-left: auto;

  width: 15px;
  height: 15px;


  transition: fill 0.3s ease;
  user-select: none;
}

/* fine container info stato messaggio */


/* contenuto del messaggio */
.message-content {
  width: 100%;
  height: 80%;
  margin-top: 2px;
}

/* Stile in caso di testo */
.message-content-text-container {
  width: 100%;
  height: auto;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.2);
  padding: 2px;

}

.message-content-text-container p {
  padding: 0 2px 0 3px;

  white-space: normal;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

/* Stile in caso di immagine */
.message-content-image-container {
  width: 200px;
  height: fit-content;
  object-fit: cover;
  user-select: none;
}

.message-content-image-container img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  user-select: none;
  border-radius: 8px;
}

/* fine contenuto del messaggio */
</style>