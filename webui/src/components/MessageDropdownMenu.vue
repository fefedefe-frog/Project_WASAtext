<script>
export default {
  props: {
    messageId: {
      type: Number,
      required: true
    },
    senderId: {
      type: String,
      required: true
    }
  },
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
  }
}
</script>

<template>
  <div class="message-hamburger">
    <!-- Toggle button -->
    <button class="showMenu" @click="show = !show">
      <svg class="feather" :style="dynamicTransformDirection"><use href="/feather-sprite-v4.29.0.svg#more-vertical" /></svg>
    </button>

    <!-- Options menu -->
    <transition name="slideDown">
      <div v-if="show" class="dropdown-list" :style="dynamicMenuSide">

        <button @click="console.log('1')" class="dropdown-item" >
          <span>Inoltra</span>
          <svg class="feather feather-mod"><use href="/feather-sprite-v4.29.0.svg#corner-up-right" /></svg>
        </button>


        <button @click="console.log('2')" class="dropdown-item">
          <span>Rispondi</span>
          <svg class="feather feather-mod"><use href="/feather-sprite-v4.29.0.svg#corner-up-left" /></svg>
        </button>

      </div>
    </transition>
  </div>
</template>

<style scoped>
/* puntini button */
.message-hamburger{
  position: relative;
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
.dropdown-list{
  position: absolute;
  top: 22px;
  z-index: 10;

  width: 100px;

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
  transform: translateY(-25%); /* quando non cliccato */
}

.slideDown-enter-to,
.slideDown-leave-from {
  transform: translateY(25%); /* normalmente visibile */
}

.slideDown-enter-active,
.slideDown-leave-active {
  transition: transform 0.3s ease;
}
/* Fine transition per menù a tendina */
</style>