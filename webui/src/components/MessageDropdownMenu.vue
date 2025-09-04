<script>
export default {
  props: {
    senderId: {
      type: String,
      required: true,
    },
    messageId: {
      type: Number,
      required: true,
    }
  },
  emits: ['respondTo', 'forwardMsg', 'commentMsg', 'deleteMsg'],
  data: function () {
    return{
      usrId: '',
      show: false
    }
  },
  computed: {
    dynamicTransformDirection(){
      return {
        'transform': this.usrId === this.senderId ? 'translate(65%, 15%)' : 'translate(-65%, 15%)'
      }
    },
    dynamicMenuSide(){
      if (this.usrId === this.senderId){
        return{
          'right' : '0'
        }
      }else {
        return{
          'left' : '0'
        }
      }
    }
  },
  mounted() {
    this.usrId= sessionStorage.getItem('usrId');

    document.addEventListener('click', this.closeIfClickOutside);
  },
  beforeUnmount() {
    document.removeEventListener('click', this.closeIfClickOutside);
  },
  methods: {
    closeIfClickOutside(event) {
      event.stopPropagation();


      // Controllo se il click è avennuto fuori dal menù dropdown e fuori dal pulsante per mostrarlo
      let clickedOnButt= this.$refs.hamburgerButton.contains(event.target)

      if(!this.$refs.menu.contains(event.target) && !clickedOnButt) {
        this.show= false;
      }
    }
  }
}
</script>

<template>
  <div class="message-hamburger">
    <!-- Toggle button -->
    <button ref="hamburgerButton" class="showMenu" @click="show = !show">
      <svg class="feather" :style="dynamicTransformDirection"><use href="/feather-sprite-v4.29.0.svg#more-vertical" /></svg>
    </button>

    <!-- Options menu -->
    <div class="dropdown-list-container" :style="dynamicMenuSide">
      <transition name="slideDown">
        <div v-show="show" ref="menu" class="dropdown-list">
          <button class="dropdown-item" @click="$emit('forwardMsg'); show= false">
            <span>Inoltra</span>
            <svg class="feather feather-mod"><use href="/feather-sprite-v4.29.0.svg#corner-up-right" /></svg>
          </button>

          <span class="dropdown-spacer" />

          <button class="dropdown-item" @click="$emit('respondTo'); show= false">
            <span>Rispondi</span>
            <svg class="feather feather-mod"><use href="/feather-sprite-v4.29.0.svg#corner-up-left" /></svg>
          </button>

          <span class="dropdown-spacer" />

          <button class="dropdown-item" @click="$emit('commentMsg'); show= false">
            <span>Commenta</span>
            <svg class="feather feather-mod"><use href="/feather-sprite-v4.29.0.svg#smile" /></svg>
          </button>

          <span v-if="usrId === senderId" class="dropdown-spacer" />

          <button v-if="usrId === senderId" class="dropdown-item" @click="$emit('deleteMsg'); show= false">
            <span>Elimina</span>
            <svg class="feather feather-mod"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
          </button>
        </div>
      </transition>
    </div>
  </div>
</template>

<style scoped>
/* puntini button */
.message-hamburger{
  position: relative;
  width: fit-content;
  height: fit-content;
}

.showMenu{
  position: relative;

  width: fit-content;
  height: fit-content;

  margin-top: 5px;
  margin-bottom: 5px;

  padding: none;
  border: none;

  background: none;

  display: flex;
  align-items: center;
  justify-content: center;

}
/* Fine puntini button */


/* Menù a tendina */
.dropdown-list-container{
  position: absolute;
  top: 22px;
  z-index: 10;

  width: 100px;
  height: fit-content;

  overflow: hidden;
}

.dropdown-list{
  position: relative;
  width: 100%;

  display: flex;
  flex-direction: column;
  align-items: center;

  justify-content: flex-start;

  padding: 4px;

  margin-top: 8px;

  border: 2px solid cornflowerblue;
  border-radius: 10px;

  background: steelblue;
  overflow: hidden;
}

.dropdown-item {
  position: relative;
  width: 100%;

  display: flex;
  flex-direction: row;
  justify-content: space-between;

  padding-right: 2px;
  padding-left: 2px;

  border-radius: 5px;

  font-size: small;
  color: rgba(200, 200, 200);
}

.dropdown-item:hover {
  transition: .4s;

  background-color: white;
  color: steelblue;
}

.dropdown-spacer{
  width: 90%;
  margin-top: 1px;
  border-bottom: 2px solid cornflowerblue;
}

.feather-mod{
  width: 20px;
  height: 20px;

  padding: none;
  margin: none;
}
/* Fine Menù a tendina */

/* Transition per menù a tendina */
.slideDown-enter-from,
.slideDown-leave-to {
  transform: translateY(-100%); /* quando non cliccato */
}

.slideDown-enter-to,
.slideDown-leave-from {
  transform: translateY(0); /* normalmente visibile */
}

.slideDown-enter-active,
.slideDown-leave-active {
  transition: transform 0.3s ease;
}
/* Fine transition per menù a tendina */
</style>