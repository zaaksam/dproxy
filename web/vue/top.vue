<template>
    <Menu mode="horizontal" theme="dark" :active-name="activeName" @on-select="onSelect">
        <div class="layout-logo">{{title}}</div>
        <div class="layout-nav">
            <Menu-item name="whiteList">
                <Icon type="ios-navigate"></Icon>
                白名单管理
            </Menu-item>
            <Menu-item name="portMap">
                <Icon type="arrow-swap"></Icon>
                端口映射
            </Menu-item>
            <Menu-item name="log">
                <Icon type="ios-analytics"></Icon>
                日志管理
            </Menu-item>
            <Menu-item name="doc">
                <Icon type="ios-paper"></Icon>
                API文档
            </Menu-item>
        </div>
    </Menu>
</template>

<style>
.layout-logo {
    width: 100px;
    height: 30px;
    background: #5b6270;
    border-radius: 3px;
    float: left;
    position: relative;
    top: 15px;
    left: 20px;
    color: #fff;
    text-align: center;
    line-height: 30px;
}

.layout-nav {
    width: 600px;
    margin: 0 auto;
}
</style>


<script lang="ts">
import Vue from 'vue'
import VueRouter from 'vue-router'
import _ from 'lodash'
import { Component } from 'vue-property-decorator'

@Component
export default class MyTop extends Vue {
    title: string = globalConfig.appName + ' ' + globalConfig.appVersion
    activeName: string = 'whiteList'

    mounted() {
        let name = _.last(_.split(this.$route.path, '/'))

        if (name != 'web') {
            this.activeName = name!
        } else {
            this.$router.replace(this.$route.path + '/' + this.activeName)
        }
    }

    onSelect(name: string) {
        this.$router.push({ name: name })
    }
}
</script>
