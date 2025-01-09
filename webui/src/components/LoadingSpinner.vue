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

<style>
.centered {
  text-align: center; /* Centra il testo */
}

.loading-spinner {
  position: fixed; /* Rimane centrato anche durante lo scroll */
  top: 50%; /* Posiziona al centro verticalmente */
  left: 50%; /* Posiziona al centro orizzontalmente */
  transform: translate(-50%, -50%); /* Compensa il 50% della larghezza e altezza per centrare perfettamente */
  background-color: #f0f0f0; /* Colore di sfondo (opzionale) */
  padding: 20px; /* Spazio interno */
  border-radius: 10px; /* Arrotonda i bordi */
  text-align: center; /* Centra il testo */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); /* Ombra per stile */
  font-size: 18px; /* Dimensione del testo */
  color: #333; /* Colore del testo */
}
</style>
