<template>
    <div class="layout">
        <div class="layout-header">
            <MyTop></MyTop>
        </div>
        <div class="layout-content">
            <div class="layout-content-main">
                <router-view></router-view>
            </div>
        </div>
        <div class="layout-footer">
            <MyBottom></MyBottom>
        </div>
    </div>
</template>

<style>
html,
body {
    height: 100%
}

body {
    margin: 0
}

.layout {
    border: 1px solid #d7dde4;
    background: #f5f7f9;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
    /* height: 100%; */
    width: 100%;
}

.layout-header {
    height: 80px;
    background: #fff;
    box-shadow: 0 1px 1px rgba(0, 0, 0, .1);
}

.layout-content {
    /* height: 90%; */
    margin: 15px;
    overflow: hidden;
    background: #fff;
    border-radius: 4px;
}

.layout-content-main {
    /* height: 90%; */
    padding: 10px;
}

.layout-footer {
    height: 30px;
}
</style>

<script lang="ts">
import Vue from 'vue'
import { Component } from 'vue-property-decorator'
import MyTop from './top.vue'
import MyBottom from './bottom.vue'

@Component({
    components: {
        MyTop,
        MyBottom
    }
})
export default class App extends Vue {
}
</script>
