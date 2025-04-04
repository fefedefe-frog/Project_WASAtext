<script>
export default {
  props: {
    userData: {
      type: Object,
      required: true,
    }
  },
  data() {
    return {
      usrId: '',
      userName: '',
      userPhoto: '',
    }
  },
  mounted() {
    this.usrId= this.userData['usrId'];
    this.userName= this.userData['userName'];
    this.userPhoto= this.userData['userPhoto'];
  }
};
</script>

<template>
  <div class="user-banner">
    <!-- Foto Profilo a destra -->
    <div v-if="userPhoto" class="user-image-container">
      <img :src="'data:image/png;base64,'+userPhoto" alt="Profile Image">
    </div>

    <!-- Contenuto del banner -->
    <div class="text-container">
      <p>{{ userName }}</p>
    </div>
  </div>
</template>

<style scoped>

.user-banner {
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

.user-banner:hover {
  background-color: darkgray;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.4);
}

.user-image-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.user-image-container {
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

.text-container p{
  font-size: 1.1em;
  font-weight: bold;
  color: #333;
}
</style>