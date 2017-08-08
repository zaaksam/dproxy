import Vue from 'vue'
import VueRouter from 'vue-router'
import iView from 'iview'
import Router from './router'

import App from '../vue/app.vue'

Vue.use(iView)

new Vue({
    el: '#app',
    router: Router,
    render: h => h(App)
})