import { createApp } from 'vue'
// import Vue from 'vue'; 这是vue2的用法
import App from './App.vue';
import router from './router'

createApp(App).use(router).mount('#app')

