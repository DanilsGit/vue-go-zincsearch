import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

import { OhVueIcon, addIcons } from "oh-vue-icons";
import { IoReloadCircle, FaArrowRight, FaArrowLeft , IoArrowBackCircleSharp, LaSearchSolid } from "oh-vue-icons/icons";

addIcons( IoReloadCircle, FaArrowRight, FaArrowLeft , IoArrowBackCircleSharp, LaSearchSolid);

createApp(App).component("v-icon", OhVueIcon).mount('#app')
