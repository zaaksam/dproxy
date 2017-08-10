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
        <Modal title="白名单申请" v-model="modal.isShow" :mask-closable="false" :closable="false">
            <Form :label-width="100">
                <Form-item label="申请者">
                    <Input v-model="modal.data.userName"></Input>
                </Form-item>
                <Form-item label="申请者IP">
                    <Input v-model="modal.data.ip" readonly></Input>
                </Form-item>
                <Form-item label="预计有效时间">
                    <Input v-model="getExpired" readonly></Input>
                </Form-item>
            </Form>
            <Alert v-show="modal.isLoading" type="warning" show-icon>数据提交中，请勿关闭此窗口...</Alert>
            <Alert v-show="modal.isErr" type="error" show-icon>{{modal.errMsg}}</Alert>
            <div slot="footer">
                <Button :disabled="modal.isLoading" @click="onModalCancel">取消</Button>
                <Button type="primary" :loading="modal.isLoading" @click="onModalOk">确定</Button>
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
    ip: string
    userID: string
    userName: string
}

interface modalModel extends BaseModalModel {
    data: modalDataModel
}

interface listModel extends BaseListModel {
    items: WhiteListModel[]
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
export default class MyWhiteList extends Vue {
    modal: modalModel = {
        isShow: false,
        isLoading: false,
        isErr: false,
        errMsg: '',
        data: {
            id: 0,
            ip: '',
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
            items: <WhiteListModel[]>[],
        },
        columns: <any[]>[
            {
                title: 'IP',
                key: 'ip'
            },
            {
                title: '申请人',
                key: 'userName'
            },
            {
                title: '申请时间',
                key: 'created',
                render: (h: Vue.CreateElement, params: any): Vue.VNode => {
                    return h('span', moment.unix(params.row.created).format('YYYY-MM-DD HH:mm:ss'))
                }
            },
            {
                title: '过期时间',
                key: 'expired',
                render: (h: Vue.CreateElement, params: any): Vue.VNode => {
                    return h('span', moment.unix(params.row.expired).format('YYYY-MM-DD HH:mm:ss'))
                }
            },
            {
                title: '操作',
                width: 160,
                render: (h: Vue.CreateElement, params: any): Vue.VNode => {
                    return h('div', [
                        h('Button', {
                            props: {
                                icon: 'edit',
                                type: 'info'
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
                                        type: 'error'
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

        Axios.get('/api/whitelist/list', { params: { pageIndex: pi, pageSize: ps } })
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

    onEdit(index: number) {
        let data = this.table.list.items[index]
        if (data) {
            this.onModalShow('', data)
        }
    }

    onDel(index: number) {
        let id: number = this.table.list.items[index].id!
        if (!id) {
            return
        }

        Axios.delete('/api/whitelist/' + id.toString())
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

    get getExpired(): string {
        return moment().add(1, 'd').format('YYYY-MM-DD HH:mm:ss')
    }

    onModalShow(event: any, data?: WhiteListModel) {
        this.modal.isShow = true
        this.modal.isLoading = false
        this.modal.isErr = false
        this.modal.errMsg = ''

        if (data) {
            this.modal.data = _.pick(data, _.keys(this.modal.data))
        } else {
            this.modal.data = {
                id: 0,
                ip: '',
                userID: '0',
                userName: ''
            }

            Axios.get('/api/whitelist/getip')
                .then((res: AxiosResponse) => {
                    if (res.data.code === 10000) {
                        this.modal.data.ip = res.data.data.ip
                    }
                })
        }
    }

    onModalOk() {
        this.modal.isErr = false
        this.modal.errMsg = ''

        let errMsg: string = ''
        if (this.modal.data.userName == '') {
            this.modal.errMsg = '申请者不能为空'
            this.modal.isErr = true
            return
        }

        this.modal.isLoading = true
        let id: string = ''
        if (this.modal.data.id > 0) {
            id = '/' + this.modal.data.id.toString()
        }

        Axios.put('/api/whitelist' + id, this.modal.data)
            .then((res: AxiosResponse) => {
                if (res.data.code === 10000) {
                    if (this.modal.data.id > 0) {
                        this.onLoad()
                    } else {
                        this.table.list.items.unshift(<WhiteListModel>res.data.data.whiteList)
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
