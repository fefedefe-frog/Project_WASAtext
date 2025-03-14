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
  <div class="alert alert-danger" role="alert">
    {{ msg }}
  </div>
</template>

<style>
</style>
