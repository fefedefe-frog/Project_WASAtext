<script>
export default {
  props: {
    messageData: {
      type: Object,
      required: true
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
      }

    }
  },
  computed: {
    dynamicMarginSide(){
      return{
        'margin-left': this.usrId === this.message.senderId ? 'auto' : '0',
      }
    }
  },
  mounted () {
    this.usrId= sessionStorage.getItem('usrId');
    this.message= this.messageData;

    this.prepMessage();
  },
  beforeUnmount() {
    this.message= {};
  },
  methods: {
    prepMessage() {
      switch (this.message['deliveryStatus']) {
        case 'sent':
          this.status = 'minus';
          break;
        case 'received':
          this.status = 'chevron-up';
          break;
        case 'read':
          this.status = 'chevrons-up';
          break;
        default:
          this.status = 'minus';
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
};
</script>

<template>
  <div class="message-div">
    <div class="message-container" :style="dynamicMarginSide">
      <div class="message-info">
        <div v-if="chatIsGroup && message['senderId'] !== usrId" class="sender-id">{{ senderName }}</div>
        <div class="message-status">
          <span class="timestamp">{{ date }}</span>
          <svg v-if="message['senderId'] === usrId" class="feather"><use :href="'/feather-sprite-v4.29.0.svg#'+ status" /></svg>
        </div>
      </div>
      <div class="message-content">
        <div v-if="message['photoContent'].length === 0" class="message-content-text-container">
          <p>{{ message["textContent"] }}</p>
        </div>
        <div v-if="message['photoContent'].length > 0" class="message-content-image-container">
          <img :src="'data:image/png;base64,'+message['photoContent']" alt="Image" draggable="false">
        </div>
      </div>
    </div>
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


/* container info stato messaggio */
.message-info {
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 100%;
  border-bottom: 1px black solid;
  user-select: none;
}

.sender-id {
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

.message-status span {
  margin-right: 3px;
}

.message-status svg{
  border: 1px black solid;
  border-radius: 50%;
  margin-left: auto;
  width: 15px;
  height: 15px;
  fill: currentColor;
  transition: fill 0.3s ease;
  user-select: none;
}

.message-status span {
  margin-left: auto;
  user-select: none;
}
/* fine container info stato messaggio */


/* contenuto del messaggio */
.message-content {
  width: 100%;
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