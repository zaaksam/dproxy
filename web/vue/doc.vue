<template>
    <Row>
        <Col span="5">
            <Menu theme="light" width="auto" @on-select="onSelect">
                <Submenu name="proxy">
                    <template slot="title">
                        <Icon type="ios-world-outline"></Icon>
                        请求代理
                    </template>
                    <Menu-item name="proxy.start">启动</Menu-item>
                    <Menu-item name="proxy.stop">停止</Menu-item>
                </Submenu>
                <Submenu name="whiteList">
                    <template slot="title">
                        <Icon type="ios-navigate"></Icon>
                        白名单
                    </template>
                    <Menu-item name="whiteList.get">获取</Menu-item>
                    <Menu-item name="whiteList.list">列表</Menu-item>
                    <Menu-item name="whiteList.clear">清理过期</Menu-item>
                    <Menu-item name="whiteList.add">添加</Menu-item>
                    <Menu-item name="whiteList.update">修改</Menu-item>
                    <Menu-item name="whiteList.delete">删除</Menu-item>
                </Submenu>
                <Submenu name="portMap">
                    <template slot="title">
                        <Icon type="arrow-swap"></Icon>
                        端口映射
                    </template>
                    <Menu-item name="portMap.get">获取</Menu-item>
                    <Menu-item name="portMap.list">列表</Menu-item>
                    <Menu-item name="portMap.add">添加</Menu-item>
                    <Menu-item name="portMap.update">修改</Menu-item>
                    <Menu-item name="portMap.delete">删除</Menu-item>
                </Submenu>
                <Submenu name="log">
                    <template slot="title">
                        <Icon type="ios-analytics"></Icon>
                        日志
                    </template>
                    <Menu-item name="log.list">列表</Menu-item>
                </Submenu>
            </Menu>
        </Col>
        <Col push="2" span="17">
            <div>
                <Form>
                    <Form-item label="地址：">{{ doc.path }}</Form-item>
                    <Form-item label="方法：">{{ doc.method }}</Form-item>
                    <Form-item label="说明：">{{ doc.title }}</Form-item>
                </Form>
            </div>
        </Col>
    </Row>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import _ from 'lodash'

interface docModel {
    alias: string
    title: string
    method: string
    path: string
}

@Component
export default class Doc extends Vue {
    preContent: string = ''

    doc: docModel = {
        alias: '',
        title: '',
        method: '',
        path: ''
    }
    docs: docModel[] = [
        {
            alias: 'proxy.start',
            title: '启动代理',
            method: 'GET',
            path: '/api/proxy/start/:id'
        }, {
            alias: 'proxy.stop',
            title: '停止代理',
            method: 'GET',
            path: '/api/proxy/stop/:id'
        }, {
            alias: 'whiteList.get',
            title: '获取白名单信息',
            method: 'GET',
            path: '/api/whitelist/:id'
        }, {
            alias: 'whiteList.list',
            title: '获取白名单列表',
            method: 'GET',
            path: '/api/whitelist/list'
        }, {
            alias: 'whiteList.clear',
            title: '清理过过期白名单',
            method: 'DELETE',
            path: '/api/whitelist/clear'
        }, {
            alias: 'whiteList.add',
            title: '添加白名单',
            method: 'PUT',
            path: '/api/whitelist'
        }, {
            alias: 'whiteList.update',
            title: '修改白名单',
            method: 'PUT',
            path: '/api/whitelist/:id'
        }, {
            alias: 'whiteList.delete',
            title: '删除白名单',
            method: 'DELETE',
            path: '/api/whitelist/:id'
        }, {
            alias: 'portMap.get',
            title: '获取端口映射信息',
            method: 'GET',
            path: '/api/portmap/:id'
        }, {
            alias: 'portMap.list',
            title: '获取端口映射列表',
            method: 'GET',
            path: '/api/portmap/list'
        }, {
            alias: 'portMap.add',
            title: '添加端口映射',
            method: 'PUT',
            path: '/api/portmap'
        }, {
            alias: 'portMap.update',
            title: '更新端口映射',
            method: 'PUT',
            path: '/api/portmap/:id'
        }, {
            alias: 'portMap.delete',
            title: '删除端口映射',
            method: 'DELETE',
            path: '/api/portmap/:id'
        }, {
            alias: 'log.list',
            title: '日志列表',
            method: 'GET',
            path: '/api/log/list'
        }
    ]

    onSelect(name: string) {
        let doc = _.find(this.docs, { alias: name })
        this.doc = doc!
    }
}
</script>
