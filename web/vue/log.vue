<template>
    <div>
        <Row>
            <Col span="24">
            <MyPage :total="table.list.total" :pageIndex="table.list.pageIndex" :pageSize="table.list.pageSize" @onLoad="onLoad"></MyPage>
            </Col>
        </Row>
        <p>&nbsp;</p>
        <Table stripe border :columns="table.columns" :data="table.list.items">
        </Table>
    </div>
</template>

<script lang="ts">
import Vue from 'vue'
import VueRouter from 'vue-router'
import { Component } from 'vue-property-decorator'
import _ from 'lodash'
import moment from 'moment'
import Axios, { AxiosResponse, AxiosError } from 'axios'
import MyPage from './page.vue'
import API from '../ts/api'

interface listModel extends BaseListModel {
    items: LogModel[]
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
export default class MyLog extends Vue {
    table: tableModel = {
        list: {
            total: 0,
            pageIndex: 1,
            pageSize: 10,
            pageCount: 0,
            items: <LogModel[]>[],
        },
        columns: <any[]>[
            {
                title: '类型',
                key: 'type'
            },
            {
                title: '内容',
                key: 'content'
            },
            {
                title: '记录时间',
                key: 'created',
                render: (h: Vue.CreateElement, params: any): Vue.VNode => {
                    return h('span', moment.unix(params.row.created).format('YYYY-MM-DD HH:mm:ss'))
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

        Axios.get(API.initURL('/api/log/list'), { params: { pageIndex: pi, pageSize: ps } })
            .then((res: AxiosResponse) => {
                this.table.list = <listModel>res.data.data.list

                if (res.data.code === 90000) {
                    this.$Message.error({ duration: 5, content: res.data.msg + '(' + res.data.code.toString() + ')' })
                }
            })
            .catch((err: AxiosError) => {
                this.$Message.error({ duration: 5, content: err.message + '(' + err.code + ')' })
            })
    }
}
</script>
