<script>
export default {
  data: function() {
    return {
      textContent: "",
      image: null,
      imagePreview: null,
      submit: false,
      maxLength: 1024
    }
  },
  emits: ['prepMessage'],
  methods: {
    prepMessage() {
      let emptyPhoto= new Blob([], {type: 'image/png'});
      if (this.textContent.trim() || this.image){
        let rawMessageData= {
          textContent: this.textContent.trim() ? this.textContent : "",
          photoContent: this.image ? this.image : emptyPhoto
        };
        this.$emit('prepMessage', rawMessageData);
      }
    },
    autoResize() {
      const el = this.$refs.textareaMessage;
      el.style.height = "auto";
      el.style.height = el.scrollHeight + "px";
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
        this.imagePreview= URL.createObjectURL(file);
      }else {
        this.image= null;
        this.imagePreview= null;
      }
    }
  }
}
</script>

<template>
  <form class="send-message-form" @submit.prevent="prepMessage">
    <div v-if="imagePreview" class="image-preview">
      <img :src="imagePreview" alt="Anteprima immagine" />
    </div>
    <textarea v-if="!imagePreview" class="textarea-content" v-model="textContent" ref="textareaMessage" placeholder="Scrivi un messaggio..." rows="2" :maxlength="maxLength" @input="autoResize" ></textarea>

    <div class="buttom-column">
      <button v-if="imagePreview" class="form-buttons delete-button" type="button" @click="imagePreview=null; image= null;">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-minus" /></svg>
      </button>

      <button class="form-buttons" type="button" @click="imageUpload">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
      </button>

      <button class="form-buttons" type="submit">
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
  border: 1px solid black;
}

.textarea-content {
  width: 100%;
  max-height: 20vh;
  margin: 0.2vh;

  border-radius: 2vh;
  padding: 1vh 1.5vh 0 1.5vh;

  resize: none;
  overflow-x: hidden;

  border: 1px solid white;
  font-size: 2.5vh;
  line-height: 3vh;
  box-sizing: border-box;
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
  width: 4.5vh;
  height: 4.5vh;

  margin-bottom: 0.4vh;

  border-radius: 1.8vh;
  border: 0.4vh dashed lightseagreen;


  color: white;
  background-color: lightseagreen;
  cursor: pointer;
  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
  transition: .4s;

  display: flex;
  justify-content: center;
  align-items: center;
}

.form-buttons:hover {
  transition: .4s;
  border: 0.4vh dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.form-buttons:active {
  background-color: lightseagreen;
}

.delete-button {
  border: 0.4vh dashed red;
  background: red;
}
.delete-button:hover {
  border: 0.4vh dashed darkred;
  background-color: white;
  color: red;
}
.delete-button:active {
  border: 0.4vh dashed darkred;
}

.image-preview {
  border: 1px solid red;
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.image-preview img {
  max-height: 20vh;
  object-fit: cover;
  object-fit: cover;
  user-select: none;
}
</style>