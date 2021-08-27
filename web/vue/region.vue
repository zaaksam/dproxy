<template>
    <div>
        <Table stripe border :columns="tableColumns" :data="tableData.items">
        </Table>
        <MyPage
            :total="tableData.total"
            :pageIndex="tableData.pageIndex"
            :pageSize="tableData.pageSize"
            @onChange="onLoad"
        ></MyPage>
    </div>
</template>

<script lang="ts">
import { CreateElement } from "vue";
import { Vue, Component } from "vue-property-decorator";
import _ from "lodash";
import moment from "moment";
import MyPage from "./page.vue";
import API from "../ts/api";

@Component({
    components: {
        MyPage
    }
})
export default class MyRegion extends Vue {

    isActive: boolean = false

    tableData: APIListModel<Model.Region> = {
        pageIndex: 1,
        pageSize: 10,
        pageCount: 0,
        total: 0,
        items: []
    };

    tableColumns: any[] = <any>[
        {
            title: "区域",
            key: "name"
        },
        {
            title: "最后心跳时间",
            key: "lastHeartbeat",
            render: (h: CreateElement, params: any) => {
                return h(
                    "span",
                    moment.unix(params.row.lastHeartbeat).format("YYYY-MM-DD HH:mm:ss")
                );
            }
        },
        {
            title: "已超时",
            render: (h: CreateElement, params: any) => {
                return h(
                    "span",
                    this.getTimeout(params.row.lastHeartbeat)
                );
            }
        },
        {
            title: "状态",
            render: (h: CreateElement, params: any) => {
                return h(
                    "Tag", {
                    props: {
                        color: this.getStateColor(params.row.lastHeartbeat),
                    },
                },
                    this.getState(params.row.lastHeartbeat)
                )
            }
        }
    ];

    mounted() {
        this.isActive = true

        this.onLoad();

        setInterval(() => {
            if (this.isActive) {
                this.onLoad();
            }
        }, 5000)
    }

    destroyed() {
        this.isActive = false
    }

    getState(lastHeartbeat: number): string {
        return (moment().unix() - lastHeartbeat > 30) ? "离线" : "在线";
    }

    getStateColor(lastHeartbeat: number): string {
        return (moment().unix() - lastHeartbeat > 30) ? "red" : "green";
    }

    getTimeout(lastHeartbeat: number): string {
        return (moment().unix() - lastHeartbeat).toString() + " 秒"
    }

    async onLoad() {
        let pi = this.tableData!.pageIndex;
        let ps = this.tableData!.pageSize;
        if (typeof this.$route.query.pi === 'string') {
            pi = _.parseInt(this.$route.query.pi);
        }
        if (typeof this.$route.query.ps === 'string') {
            ps = _.parseInt(this.$route.query.ps);
        }

        let pms = new URLSearchParams();
        pms.set("pageIndex", pi.toString());
        pms.set("pageSize", ps.toString());

        let result = await API.get<Model.Regions>("/region/list/?" + pms.toString());
        if (result.code === 10000) {
            this.tableData = result.data!.list;
        } else {
            this.$Message.error({
                duration: 5,
                content: result.msg + "(" + result.code.toString() + ")"
            });
        }
    }
}
</script>
