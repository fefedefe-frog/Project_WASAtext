<script>
export default {
	props: {
    msg: {
      type: String,
      required: true
    }
  },
  emits: ['close'],
}
</script>

<template>
  <div class="alert alert-danger error-div" role="alert">
    <span>{{ msg }}</span>
    <button class="close-button" @click="this.$emit('close')">Close</button>
  </div>
</template>

<style scoped>
.error-div {
  position: absolute;

  top: 2%;
  left: 50%;

  transform: translate(-50%);
  user-select: none;

  display: flex;
  flex-direction: row;

  align-items: center;
}

.close-button{
  width: 50px;
  height: 30px;

  margin-left: 10px;

  border: 1px solid red;
  border-radius: 5px;
}
</style>
