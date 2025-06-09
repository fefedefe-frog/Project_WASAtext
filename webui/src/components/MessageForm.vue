<script>
export default {
  emits: ['prepMessage'],
  data: function() {
    return {
      textContent: "",
      image: null,
      imageName: "",
      submit: false,
      maxLength: 1024
    }
  },
  computed: {
    isMessageValid() {
      return this.textContent.trim() || this.image ? false : true;
    }
  },
  methods: {
    prepMessage() {
      let emptyPhoto= new Blob([], {type: 'image/png'});
      if (this.textContent.trim() || this.image){
        let rawMessageData= {
          textContent: this.textContent.trim() ? this.textContent : "",
          photoContent: this.image ? this.image : emptyPhoto
        };
        this.$emit('prepMessage', rawMessageData);
        this.textContent= "";
        this.image= null;
        this.imageName= "";
      }
    },
    imageUpload() {
      const input= document.createElement('input');
      input.type= "file";
      input.accept= "image/*";
      input.addEventListener("change", this.handleFileChange);
      input.click();
    },
    handleFileChange(event) {
      const file= event.target.files[0];
      if(file && file.type.startsWith("image/")) {
        this.image= file;
        this.imageName= file.name;
      }else {
        this.image= null;
        this.imageName= "";
      }
    }
  }
}
</script>

<template>
  <form class="send-message-form" @submit.prevent="prepMessage">
    <textarea v-if="!image" ref="textareaMessage" v-model="textContent" class="textarea-content" placeholder="Scrivi un messaggio..." rows="2" :maxlength="maxLength" />
    <div v-if="image" class="image-name">
      <button class="form-buttons delete-button" type="button" @click="image=null; imageName= '';">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
      </button>
      <span>{{ imageName }}</span>
    </div>

    <div class="buttom-column">
      <button class="form-buttons" type="button" @click="imageUpload">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
      </button>

      <button class="form-buttons" type="submit" :disabled="isMessageValid">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigation" /></svg>
      </button>
    </div>
  </form>
</template>

<style scoped>
.send-message-form {
  width: 100%;
  height: 100%;

  display: flex;
  flex-direction: row;
  justify-content: right;

}

.textarea-content {
  width: 100%;
  max-height: 100px;
  margin: 2px;

  border-radius: 10px;
  padding: 2px 4px 0 4px;

  resize: none;
  overflow-x: hidden;
  font-size: 0.8rem;
  line-height: 1rem;
  box-sizing: border-box;
  border: 1px solid white;
}

.buttom-column {
  width: fit-content;
  height: fit-content;
  margin-top: auto;
  margin-left: 0.3vh;
  display: flex;
  flex-direction: column;
}

.form-buttons {
  width: 25px;
  height: 25px;

  margin: 1px 5px 2px 0;

  border-radius: 25%;
  border: 2px dashed lightseagreen;


  color: white;
  background-color: lightseagreen;
  cursor: pointer;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
  transition: .4s;

  display: flex;
  justify-content: center;
  align-items: center;
}

.form-buttons:not(:disabled):hover {
  transition: .4s;
  border: 2px dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.form-buttons:active {
  background-color: lightseagreen;
}

.form-buttons:disabled{
  background-color: gray;
  border: 2px solid gray;
  cursor: default;
}

.delete-button {
  border: 2px dashed red;
  background: red;
}
.delete-button:hover {
  border: 2px dashed darkred;
  background-color: white;
  color: red;
}
.delete-button:active {
  border: 2px dashed darkred;
}

.image-name {
  height: fit-content;
  object-fit: contain;

  margin-right: 2px;

  display: flex;
  flex-direction: row;

  align-items: center;
  justify-content: center;

}
</style>