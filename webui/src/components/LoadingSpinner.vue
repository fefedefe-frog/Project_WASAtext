<script>
export default {
	props: ["loading", "loadingText"]
}
</script>

<template>
	<div v-if="loading">
		<div class="loading-spinner">
			<div class="spinner-border" role="status">
				<span class="visually-hidden"></span>
			</div>
    <p class="visually-visible centered">{{ loadingText }}...</p>
		</div>
	</div>
	<div v-if="!loading"><slot /></div>
</template>

<style scoped>
.centered {
  text-align: center; /* Centra il testo */
}

.loading-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: #f0f0f0;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  font-size: 18px;
  color: #333;
}
</style>
