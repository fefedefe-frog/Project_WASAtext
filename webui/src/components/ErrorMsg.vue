<script>
export default {
	props: {
    msg: {
      type: Error,
      required: true
    }
  }
}
</script>

<template>
  <div class="alert alert-danger error-div" role="alert">
    {{ msg }}
  </div>
</template>

<style scoped>
.error-div {
  position: absolute;
  left: 50%;
  transform: translate(-50%);
  user-select: none;
}
</style>
