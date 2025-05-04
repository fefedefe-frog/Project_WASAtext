<script>
export default {
	props: {
    userData: {
      type: Object,
      required: true
    }
  },
  data: function () {
    return {
    }
  },
}
</script>

<template>
  <div class="user-container">
    <div class="user-info">
      <div v-if="this.userData['userPhoto']" class="user-image-container">
        <img :src="'data:image/png;base64,'+userPhoto" alt="Profile Image">
      </div>

      <div class="text-container">
        <p>{{ this.userData['username'] }}</p>
      </div>
    </div>
    <div class="send-message"></div>
  </div>
</template>

<style scoped>
.user-container {
  width: fit-content;
  height: fit-content;

  padding: 20px;
  border-radius: 10px;

  background-color: #a0a0a0;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.2);
}
.user-info {
  border: 1px violet solid;

  width: 100%;
}

.send-message {
  border: 1px orangered solid;

  width: 100%;
}
</style>
