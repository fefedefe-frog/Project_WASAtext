import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';

import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import ChatBanner from './components/ChatBanner.vue'
import UserBanner from './components/UserBanner.vue'
import ChatMessage from './components/ChatMessage.vue'
import Chat from './components/Chat.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;

app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("ChatBanner", ChatBanner);
app.component("UserBanner", UserBanner);
app.component("ChatMessage", ChatMessage);
app.component("Chat", Chat);

app.use(router)
app.mount('#app')
