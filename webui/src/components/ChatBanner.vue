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
      chatId: -1,
      message: {}
    }
  },
  mounted() {
    this.chatId= this.chatData['chatId'];
    console.log(this.chatData)
  },
  methods: {
    // TODO
    loadChat(){
      this.$router.push('/chats/'+ this.chatId);
    }
  },
};
</script>

<template>
  <div class="chat-banner" @click="loadChat">
    <!-- Foto Profilo a destra -->
    <div class="chat-immage-container">
      <img :src="'data:image/png;base64,'+this.chatData['chatPhoto']" alt="Chat Image">
    </div>

    <!-- Contenuto del banner -->
    <div class="text-container">
      <div class="chat-name">{{ chatData['chatName'] }}</div>
      <div class="last-message"> messaggo tem...</div>
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

.chat-immage-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.chat-immage-container {
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