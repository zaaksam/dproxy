<template>
    <div>
        <Row>
            <Col span="20">
            <MyPage :total="table.list.total" :pageIndex="table.list.pageIndex" :pageSize="table.list.pageSize" @onLoad="onLoad"></MyPage>
            </Col>
            <Col push="3" span="1">
            <div style="padding:0px 0px 10px 0px;">
                <Button @click="onModalShow" type="primary" icon="plus-round"></Button>
            </div>
            </Col>
        </Row>
        <Table stripe border :columns="table.columns" :data="table.list.items">
        </Table>
        <Modal title="端口映射" v-model="modal.isShow" :mask-closable="false" :closable="false">
            <Form :label-width="80">
                <Form-item label="标题">
                    <Input v-model="modal.data.title"></Input>
                </Form-item>
                <Form-item label="源IP">
                    <Input v-model="modal.data.sourceIP"></Input>
                </Form-item>
                <Form-item label="源端口">
                    <Input v-model="modal.data.sourcePort" @on-keydown="onKeyDown" number></Input>
                </Form-item>
                <Form-item label="目标IP">
                    <Input v-model="modal.data.targetIP"></Input>
                </Form-item>
                <Form-item label="目标端口">
                    <Input v-model="modal.data.targetPort" @on-keydown="onKeyDown" number></Input>
                </Form-item>
                <Form-item label="创建者">
                    <Input v-model="modal.data.userName"></Input>
                </Form-item>
            </Form>
            <Alert v-show="modal.isLoading" type="warning" show-icon>数据提交中，请勿关闭此窗口...</Alert>
            <Alert v-show="modal.isErr" type="error" show-icon>{{modal.errMsg}}</Alert>
            <div slot="footer">
                <Button :disabled="modal.isLoading" @click="onModalCancel">取消</Button>
                <Button type="primary" :loading="modal.isLoading" @click="onModalOK">确定</Button>
            </div>
        </Modal>
    </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { Component } from 'vue-property-decorator'
import _ from 'lodash'
import moment from 'moment'
import Axios, { AxiosResponse, AxiosError } from 'axios'
import MyPage from './page.vue'

interface modalDataModel {
    id: number
    title: string
    sourceIP: string
    sourcePort: number
    targetIP: string
    targetPort: number
    userID: string
    userName: string
}

interface modalModel extends BaseModalModel {
    data: modalDataModel
}

interface listModel extends BaseListModel {
    items: PortMapModel[]
}

interface tableModel {
    columns: any[]
    list: listModel
}

@Component({
    components: {
        MyPage
    }
})
export default class MyPortMap extends Vue {
    modal: modalModel = {
        isShow: false,
        isLoading: false,
        isErr: false,
        errMsg: '',
        data: {
            id: 0,
            title: '',
            targetIP: '',
            targetPort: 0,
            sourceIP: '',
            sourcePort: 0,
            userID: '',
            userName: ''
        }
    }

    table: tableModel = {
        list: {
            total: 0,
            pageIndex: 1,
            pageSize: 10,
            pageCount: 0,
            items: <PortMapModel[]>[],
        },
        columns: [
            {
                title: '标题',
                key: 'title'
            },
            {
                title: '源IP',
                key: 'sourceIP'
            },
            {
                title: '源端口',
                key: 'sourcePort'
            },
            {
                title: '目标IP',
                key: 'targetIP'
            },
            {
                title: '目标端口',
                key: 'targetPort'
            },
            {
                title: '创建者',
                key: 'userName'
            },
            {
                title: '创建时间',
                key: 'created',
                render: (h: Vue.CreateElement, params: any): Vue.VNode => {
                    return h('span', moment.unix(params.row.created).format('YYYY-MM-DD HH:mm:ss'))
                }
            },
            {
                title: '操作',
                width: 220,
                render: (h: Vue.CreateElement, params: any): Vue.VNode => {
                    return h('div', [
                        h('Button', {
                            props: {
                                type: this.getStartOrStopType(params.index),
                                icon: this.getStartOrStopIcon(params.index)
                            },
                            style: {
                                marginRight: '10px'
                            },
                            on: {
                                click: () => {
                                    this.onStartOrStop(params.index)
                                }
                            }
                        }),
                        h('Button', {
                            props: {
                                icon: 'edit',
                                type: 'info',
                                disabled: this.getIsStart(params.index)
                            },
                            style: {
                                marginRight: '10px'
                            },
                            on: {
                                click: () => {
                                    this.onEdit(params.index)
                                }
                            }
                        }),
                        h('Poptip', {
                            props: {
                                confirm: true,
                                placement: 'left',
                                title: '您确认要删除这条内容吗？'
                            },
                            on: {
                                'on-ok': () => {
                                    this.onDel(params.index)
                                }
                            }
                        }, [
                                h('Button', {
                                    props: {
                                        icon: 'close-round',
                                        type: 'error',
                                        disabled: this.getIsStart(params.index)
                                    }
                                })
                            ])
                    ])
                }
            }
        ]
    }

    mounted() {
        this.onLoad()
    }

    onLoad() {
        let pi = _.parseInt(this.$route.query.pi)
        let ps = _.parseInt(this.$route.query.ps)
        if (_.isNaN(pi)) {
            pi = this.table.list.pageIndex
        }
        if (_.isNaN(ps)) {
            ps = this.table.list.pageSize
        }

        Axios.get('/api/portmap/list', { params: { pageIndex: pi, pageSize: ps } })
            .then((res: AxiosResponse) => {
                this.table.list = <listModel>res.data.data.list

                if (res.data.code === 90000) {
                    this.$Message.error({ content: res.data.msg + '(' + res.data.code.toString() + ')' })
                }
            })
            .catch((err: AxiosError) => {
                this.$Message.error({ content: err.message + '(' + err.code + ')' })
            })
    }

    getStartOrStopType(index: number): string {
        return this.table.list.items[index].isStart ? 'warning' : 'success'
    }

    getStartOrStopIcon(index: number): string {
        return this.table.list.items[index].isStart ? 'stop' : 'play'
    }

    getIsStart(index: number): boolean {
        return this.table.list.items[index].isStart
    }

    onKeyDown(event: any) {
        if (_.isNaN(_.parseInt(event.key))) {
            event.returnValue = false
        }
    }

    onStartOrStop(index: number) {
        let isStart = !this.table.list.items[index].isStart
        let str = isStart ? 'start' : 'stop'

        Axios.get('/api/proxy/' + str + '/' + this.table.list.items[index].id.toString())
            .then((res: AxiosResponse) => {
                if (res.data.code === 10000) {
                    this.table.list.items[index].isStart = isStart
                } else {
                    this.$Message.error({ content: '代理启动失败：' + res.data.msg + '(' + res.data.code.toString() + ')' })
                }
            }).catch((err: AxiosError) => {
                this.$Message.error({ content: '代理启动失败：' + err.message + '(' + err.code + ')' })
            })
    }

    onEdit(index: number) {
        let data = this.table.list.items[index]
        if (data) {
            this.onModalShow('', data)
        }
    }

    onDel(index: number) {
        let id: number = this.table.list.items[index].id

        Axios.delete('/api/portmap/' + id.toString())
            .then((res: AxiosResponse) => {
                if (res.data.code === 10000) {
                    this.onLoad()
                } else {
                    this.$Message.error({ content: res.data.msg + '(' + res.data.code.toString() + ')' })
                }
            })
            .catch((err: AxiosError) => {
                this.$Message.error({ content: err.message + '(' + err.code + ')' })
            })
    }

    onModalShow(event: any, data?: PortMapModel) {
        this.modal.isShow = true
        this.modal.isLoading = false
        this.modal.isErr = false
        this.modal.errMsg = ''

        if (data) {
            this.modal.data = _.pick(data, _.keys(this.modal.data))
        } else {
            this.modal.data = {
                id: 0,
                title: '',
                targetIP: '127.0.0.1',
                targetPort: 80,
                sourceIP: '127.0.0.1',
                sourcePort: 80,
                userID: '0',
                userName: ''
            }
        }
    }

    onModalOK() {
        this.modal.isErr = false
        this.modal.errMsg = ''

        let errMsg: string = ''
        if (this.modal.data.sourceIP == '') {
            errMsg = '源IP不能为空'
        } else if (this.modal.data.sourcePort <= 0 || this.modal.data.sourcePort > 65536) {
            errMsg = '源端口数值错误，不能为0，最大不能超过65536'
        } else if (this.modal.data.targetIP == '') {
            errMsg = '目标IP不能为空'
        } else if (this.modal.data.targetPort <= 0 || this.modal.data.targetPort > 65536) {
            errMsg = '目标端口数值错误，不能为0，最大不能超过65536'
        } else if (this.modal.data.userName == '') {
            errMsg = '创建者不能为空'
        }
        if (errMsg != '') {
            this.modal.errMsg = errMsg
            this.modal.isErr = true
            return
        }

        this.modal.isLoading = true
        let id: string = ''
        if (this.modal.data.id > 0) {
            id = '/' + this.modal.data.id.toString()
        }

        Axios.put('/api/portmap' + id, this.modal.data)
            .then((res: AxiosResponse) => {
                if (res.data.code === 10000) {
                    if (this.modal.data.id > 0) {
                        this.onLoad()
                    } else {
                        this.table.list.items.unshift(<PortMapModel>res.data.data.portMap)
                        this.table.list.total += 1
                    }
                    this.modal.isShow = false
                } else {
                    this.modal.errMsg = res.data.msg + '(' + res.data.code.toString() + ')'
                    this.modal.isErr = true
                }
                this.modal.isLoading = false
            }).catch((err: AxiosError) => {
                this.modal.errMsg = err.message + '(' + err.code + ')'
                this.modal.isErr = true
                this.modal.isLoading = false
            })
    }

    onModalCancel() {
        this.modal.isShow = false
    }
}
</script>
