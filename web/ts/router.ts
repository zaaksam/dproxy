import Vue from 'vue'
import VueRouter from 'vue-router'

import MyWhiteList from '../vue/whitelist.vue'
import MyPortMap from '../vue/portmap.vue'
import MyLog from '../vue/log.vue'
import MyDoc from '../vue/doc.vue'

Vue.use(VueRouter)

const routes = <VueRouter.RouteConfig[]>[
    { path: globalConfig.prefixPath + '/web/whitelist', name: 'whitelist', component: MyWhiteList },
    { path: globalConfig.prefixPath + '/web/portmap', name: 'portmap', component: MyPortMap },
    { path: globalConfig.prefixPath + '/web/log', name: 'log', component: MyLog },
    { path: globalConfig.prefixPath + '/web/doc', name: 'doc', component: MyDoc }
]

const Router = new VueRouter({
    mode: 'history',
    routes: routes
})

export default Router