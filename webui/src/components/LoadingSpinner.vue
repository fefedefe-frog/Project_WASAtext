<script>
export default {
	props: {
    loading:{
      type: Boolean,
      required: true
    },
    loadingText:{
      type: String,
      required: false
    }
  },
}
</script>

<template>
  <div v-if="loading">
    <div class="loading-spinner">
      <div class="spinner-border" role="status">
        <span class="visually-hidden" />
      </div>
      <p class="visually-visible centered">{{ loadingText }}...</p>
    </div>
  </div>
</template>

<style scoped>
.centered {
  text-align: center; /* Centra il testo */
}

.loading-spinner {
  position: relative;
  width: fit-content;
  height: fit-content;
  background-color: #f0f0f0;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  font-size: 18px;
  color: #333;
}
</style>
