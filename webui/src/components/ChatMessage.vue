<script>
export default {
  props: {
    message: {
      type: Object,
      required: true
    },
    senderName: {
      type: String
    }
  },
  data: function () {
    return{
      status: 'minus',
      date: '',
      contentType: '',
      usrId: '',
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

    this.prepMessage()
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
      const dateObject= new Date(this.message['timestamp']);

      // Formatto giorno mese anno dal timestamp
      let dateFormatter = new Intl.DateTimeFormat('it-IT', { dateStyle: 'short' });
      let formattedDate = dateFormatter.format(dateObject);

      // Formatto ore e minuti
      let timeFormatter = new Intl.DateTimeFormat('it-IT', { hour: '2-digit', minute: '2-digit' });
      let formattedTime = timeFormatter.format(dateObject);

      // Unisco nel formato che mi interessa
      this.date= formattedDate +" "+ formattedTime;

      this.contentType = this.message['contentType'];

    }
  }
};
</script>

<template>
  <div class="message-div">
    <div class="message-container" :style="dynamicMarginSide">
      <div class="message-info">
        <div class="sender-id">{{ senderName === usrId ? "tu" : senderName }}</div>
        <div class="message-status">
          <span class="timestamp">{{ date }}</span>
          <svg class="feather"><use :href="'/feather-sprite-v4.29.0.svg#'+ status" /></svg>
        </div>
      </div>
      <div class="message-content">
        <div v-if="contentType === 'text' " class="message-content-text-container">
          <p>{{ message["content"] }}</p>
        </div>
        <div v-if="contentType === 'photo'" class="message-content-image-container">
          <img :src="'data:image/png;base64,'+message['content']" alt="Image">
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
}

.message-container {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
  padding: 2px;
  margin: 5px;
  background-color: lightskyblue;
  border-radius: 8px;
  border: 1px deepskyblue solid;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  height: fit-content;
  width: fit-content;
  max-width: 500px;
}


/* container info stato messaggio */
.message-info {
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 100%;
  border-bottom: 1px black solid;
}

.sender-id {
  margin-right: 10px;
  font-size: 1.1em;
  font-style: italic;
  font-weight: 550;
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
  display: block;
  align-items: center;
  width: 100%;
  margin: 2px;
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
  padding-left: 2px;
  padding-right: 2px;
}

/* Stile in caso di immagine */
.message-content-image-container {
  width: 200px;
  height: 200px;
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