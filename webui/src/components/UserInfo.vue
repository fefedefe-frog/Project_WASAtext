<script>
export default {
	props: {
    userData: {
      type: Object,
      required: true
    }
  },
  emits: ['closeUserInfo', 'reqToSendMessage'],
  data: function () {
    return {
      errormsg: null,

      user: {
        usrId: "",
        userName: "",
        userPhoto: "",
      },

      messageContent: "",
    }
  },
  methods: {
    sendMessage(){
      let usrId= this.user['usrId']
      let rawMess= {
        contentType: "text",
        content: this.messageContent
      }

      this.$emit('reqToSendMessage', {sendTo: usrId, rawMess: rawMess})
    }
  },
  mounted (){
    this.user['usrId']= this.userData['usrId'];
    this.user['userName']= this.userData['userName'];
    this.user['userPhoto']= this.userData['userPhoto'];
  }
}
</script>

<template>
  <div class="user-container">
    <div class="btn-toolbar mb-2 mb-md-0 w-100">
      <div class="btn-group me-2">
        <button type="button" class="btn btn-sm btn-outline-danger shadow-none" @click="this.$emit('closeUserInfo')">
          Chiudi
        </button>
      </div>
    </div>
    <div class="user-info">
      <div v-if="user['userPhoto']" class="user-image-container">
        <img :src="'data:image/png;base64,'+ user['userPhoto']" alt="Profile Image">
      </div>

      <div class="text-container">
        <p>{{ user['userName'] }}</p>
      </div>
    </div>

    <div class="send-message">
      <h4>Invia un messaggio</h4>
      <form class="sendMessage-form" @submit.prevent="sendMessage">
        <label for="text">
          <input id="textContent" v-model="messageContent" type="text" placeholder="Scrivi un messaggio" required>
        </label>
        <button type="submit" :disabled="!messageContent" :class="{ disabled: !messageContent}">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigator" /></svg>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.user-container {
  height: 70%;
  width: 50%;

  display: flex;
  flex-direction: column;

  align-items: center;

  padding: 5px;
  margin-bottom: 5px;

  background-color: lightgray;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.user-info {
  border-bottom: 1px grey solid;

  display: flex;
  flex-direction: column;
  align-items: center;

  width: 80%;
}

.user-image-container {
  max-width: 100%;  /* L'immagine non andrà mai oltre la larghezza del suo contenitore */
  max-height: 100%; /* L'immagine non andrà mai oltre l'altezza del suo contenitore */
  object-fit: contain; /* L'immagine si adatta dentro il box senza distorsioni */
}

.user-image-container img {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
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

.text-container p{
  font-size: 2rem;
  color: #333;
}


.send-message {
  display: flex;
  flex-direction: column;

  align-items: center;
  border: 2px violet solid;

  width: 100%;
  height: 20%;
}
</style>
