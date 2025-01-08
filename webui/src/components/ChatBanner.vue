<script>
export default {
  props: ["chatName", "chatPhoto", "lastMessage", "chatId"],
  methods: {
    loadChat(){
      console.log("loading chat"+ this.chatId)
      this.$router.push('/chats/' + this.chatId);
      let div = document.getElementById("chat-banner");
      div.style.backgroundColor= d5d5d5;
    }
  },
};
</script>

<template>
  <div class="chat-banner" @onclick="loadChat">
    <!-- Foto Profilo a destra -->
    <div class="chat-immage-container">
      <img :src="'data:image/png;base64,'+chatPhoto" alt="User Profile" />
    </div>

    <!-- Contenuto del banner -->
    <div class="text-container">
      <div class="chat-name">{{ chatName }}</div>
      <div class="last-message">{{ lastMessage }}</div>
    </div>
  </div>
</template>

<style scoped>

.chat-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;
  background-color: #f5f5f5;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.chat-banner:hover {
  background-color: #e5e5e5;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.4);
}

.chat-immage-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.chat-immage-container {
  max-width: 100%;  /* L'immagine non andrà mai oltre la larghezza del suo contenitore */
  max-height: 100%; /* L'immagine non andrà mai oltre l'altezza del suo contenitore */
  object-fit: contain; /* L'immagine si adatta dentro il box senza distorsioni */
  border: 5px solid #ffcd39;
}

.text-container {
  flex: 1;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-left: 10px;
  border: #0a53be 1px solid;
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