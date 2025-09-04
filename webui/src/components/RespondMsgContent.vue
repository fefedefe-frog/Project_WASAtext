<script>
export default {
  props: {
    messageData: {
      type: Object,
      required: true
    },
  },
  data: function () {
    return{
      status: "minus",
      date: "",
      textContent: "",
      photoContent: {},
      usrId: "",

      message: {
        senderId: "",
        textContent: "",
        photoContent: [],
      }

    }
  },
  mounted () {
    this.usrId= sessionStorage.getItem('usrId');
    this.message= this.messageData;
  },
  beforeUnmount() {
    this.message= {};
  },
  methods: {
  }
}
</script>

<template>
  <div class="respond-container">
    <span class="respond-to-name">{{ message['senderName'] }}</span>
    <div v-if="!message['photoContent']" class="respond-content-text-container">
      <p>{{ message["textContent"] }}</p>
    </div>
    <div v-if="message['photoContent']" class="respond-content-image-container">
      <img :src="'data:image/png;base64,'+message['photoContent']" alt="Image" draggable="false">
    </div>
  </div>
</template>

<style scoped>
.respond-container {
  height: fit-content;
  width: 100%;

  border-radius: 8px;
  border: 1px steelblue solid;

  display: flex;
  flex-direction: column;

  justify-content: center;
  align-items: center;

  background-color: cadetblue;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}


/* contenuto del messaggio */
.respond-to-name {
  width: fit-content;
  max-width: 90%;
  height: 10%;

  margin-left: 3px;
  margin-right: auto;
  margin-bottom: 2px;

  font-size: 0.9rem;
  font-style: italic;
  font-weight: bold;

  color: #333;


  overflow: hidden;
  text-overflow: ellipsis;

  border-bottom: 1px solid black;
}

/* Stile in caso di testo */
.respond-content-text-container {
  width: 96%;
  height: 90%;

  margin-right: 2%;
  margin-left: 2%;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.2);
  padding: 2px;

  overflow-y: auto;
}

.respond-content-text-container p {
  padding: 0 2px 0 3px;

  white-space: normal;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

/* Stile in caso di immagine */
.respond-content-image-container {
  height: 100%;
  aspect-ratio: 1/1;
  user-select: none;

  display: flex;
  justify-content: center;
  align-items: center;
}

.respond-content-image-container img {
  max-width: 95%;
  max-height: 95%;
  user-select: none;
  border-radius: 8px;
}

/* fine contenuto del messaggio */
</style>