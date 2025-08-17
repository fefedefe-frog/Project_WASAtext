import {createApp} from 'vue';
import App from './App.vue';
import router from './router';
import axios from './services/axios.js';

import ErrorMsg from './components/ErrorMsg.vue';
import LoadingSpinner from './components/LoadingSpinner.vue';

import SidebarList from "./components/SidebarList.vue";
import ChatBanner from './components/ChatBanner.vue';
import UserBanner from './components/UserBanner.vue';

import ChatMessage from './components/ChatMessage.vue';
import RespondMsgContent from "./components/RespondMsgContent.vue";
import MessageForm from "./components/MessageForm.vue";
import MessageDropdownMenu from "./components/MessageDropdownMenu.vue";

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;

app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);

app.component("SidebarList", SidebarList);
app.component("ChatBanner", ChatBanner);
app.component("UserBanner", UserBanner);

app.component("ChatMessage", ChatMessage);
app.component("RespondMsgContent", RespondMsgContent);
app.component("MessageForm", MessageForm);
app.component("MessageDropdownMenu", MessageDropdownMenu);

app.use(router)
app.mount('#app')
