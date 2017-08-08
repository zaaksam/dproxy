import Vue from 'vue'
import VueRouter from 'vue-router'

import MyWhiteList from '../vue/whiteList.vue'
import MyPortMap from '../vue/portMap.vue'
import MyLog from '../vue/log.vue'
import MyDoc from '../vue/doc.vue'

Vue.use(VueRouter)

const routes: VueRouter.RouteConfig[] = [
    { path: '/web/whitelist', name: 'whiteList', component: MyWhiteList },
    { path: '/web/portmap', name: 'portMap', component: MyPortMap },
    { path: '/web/log', name: 'log', component: MyLog },
    { path: '/web/doc', name: 'doc', component: MyDoc }
]

const Router = new VueRouter({
    mode: 'history',
    routes: routes
})

export default Router