<template>
    <div>
        <Row type="flex" style="margin-top: 10px; margin-bottom: 15px">
            区域：
            <Input
                v-model="query.region"
                placeholder="模糊关键字"
                style="width: 160px"
            ></Input>
            &nbsp;&nbsp;&nbsp;&nbsp;源端口：
            <Input
                v-model="query.sourcePort"
                placeholder="模糊关键字"
                style="width: 160px"
            ></Input>
            &nbsp;&nbsp;&nbsp;&nbsp;目标IP：
            <Input
                v-model="query.targetIP"
                placeholder="模糊关键字"
                style="width: 160px"
            ></Input>
            &nbsp;&nbsp;&nbsp;&nbsp;目标端口：
            <Input
                v-model="query.targetPort"
                placeholder="模糊关键字"
                style="width: 160px"
            ></Input>
            &nbsp;&nbsp;&nbsp;&nbsp;排序字段：
            <RadioGroup v-model="query.sortField">
                <Radio label="created">
                    <span>创建时间</span>
                </Radio>
                <Radio label="sourcePort">
                    <span>源端口</span>
                </Radio>
            </RadioGroup>
            &nbsp;&nbsp;&nbsp;&nbsp;排序顺序：
            <RadioGroup v-model="query.sortDesc">
                <Radio label="1">
                    <span>降序</span>
                </Radio>
                <Radio label="0">
                    <span>升序</span>
                </Radio>
            </RadioGroup>
            &nbsp;&nbsp;&nbsp;&nbsp;
            <Button type="primary" @click="onLoad">查找</Button>
        </Row>
        <Table stripe border :columns="tableColumns" :data="tableData.items">
        </Table>
        <MyPage
            :total="tableData.total"
            :pageIndex="tableData.pageIndex"
            :pageSize="tableData.pageSize"
            @onChange="onLoad"
        ></MyPage>
        <Modal
            :title="modal.title"
            v-model="modal.isShow"
            :mask-closable="false"
            :closable="false"
        >
            <Form :label-width="80">
                <Form-item label="ID" v-show="modal.isEdit">
                    <Input v-model="modal.data.id" readonly></Input>
                </Form-item>
                <Form-item label="区域">
                    <Input v-model="modal.data.region"></Input>
                </Form-item>
                <Form-item label="标题">
                    <Input v-model="modal.data.title"></Input>
                </Form-item>
                <Form-item label="源IP">
                    <Input v-model="modal.data.sourceIP"></Input>
                </Form-item>
                <Form-item label="源端口">
                    <Input
                        v-model="modal.data.sourcePort"
                        @on-keydown="onKeyDown"
                        number
                    ></Input>
                </Form-item>
                <Form-item label="目标IP">
                    <Input v-model="modal.data.targetIP"></Input>
                </Form-item>
                <Form-item label="目标端口">
                    <Input
                        v-model="modal.data.targetPort"
                        @on-keydown="onKeyDown"
                        number
                    ></Input>
                </Form-item>
                <Form-item label="创建者ID">
                    <Input v-model="modal.data.userID"></Input>
                </Form-item>
                <Form-item label="创建者名称">
                    <Input v-model="modal.data.userName"></Input>
                </Form-item>
            </Form>
            <Alert v-show="modal.isLoading" type="warning" show-icon
                >数据提交中，请勿关闭此窗口...</Alert
            >
            <Alert v-show="modal.isErr" type="error" show-icon>{{
                modal.errMsg
            }}</Alert>
            <div slot="footer">
                <Button :disabled="modal.isLoading" @click="onModalCancel"
                    >取消</Button
                >
                <Button
                    type="primary"
                    :loading="modal.isLoading"
                    @click="onModalOK"
                    >{{ modal.okBtnText }}</Button
                >
            </div>
        </Modal>
    </div>
</template>

<script lang="ts">
import { CreateElement } from "vue";
import { Vue, Component } from "vue-property-decorator";
import _ from "lodash";
import moment from "moment";
import MyPage from "./page.vue";
import API from "../ts/api";

interface modalDataModel {
    id: number;
    region: string;
    title: string;
    sourceIP: string;
    sourcePort: number;
    targetIP: string;
    targetPort: number;
    userID: string;
    userName: string;
}

interface modalModel extends BaseModalModel {
    data: modalDataModel;
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
        errMsg: "",
        isEdit: false,
        title: "",
        okBtnText: "",
        data: {
            id: 0,
            region: "",
            title: "",
            targetIP: "",
            targetPort: 0,
            sourceIP: "",
            sourcePort: 0,
            userID: "",
            userName: ""
        }
    };

    query: Model.PortMapQuery = {
        sortField: 'created',
        sortDesc: '1'
    };

    tableData: APIListModel<Model.PortMap> = {
        pageIndex: 1,
        pageSize: 10,
        pageCount: 0,
        total: 0,
        items: []
    };

    tableColumns: any[] = [
        {
            title: "区域",
            key: "region"
        },
        {
            title: "标题",
            key: "title"
        },
        {
            title: "源IP",
            key: "sourceIP"
        },
        {
            title: "源端口",
            key: "sourcePort"
        },
        {
            title: "目标IP",
            key: "targetIP"
        },
        {
            title: "目标端口",
            key: "targetPort"
        },
        {
            title: "创建者ID",
            key: "userID"
        },
        {
            title: "创建者名称",
            key: "userName"
        },
        {
            title: "创建时间",
            key: "created",
            render: (h: CreateElement, params: any) => {
                return h(
                    "span",
                    moment.unix(params.row.created).format("YYYY-MM-DD HH:mm:ss")
                );
            }
        },
        {
            width: 220,
            renderHeader: (h: CreateElement, params: any) => {
                return h("Button", {
                    props: {
                        size: "small",
                        type: "primary",
                        icon: "plus-round"
                    },
                    on: {
                        click: () => {
                            this.onModalShow();
                        }
                    }
                });
            },
            render: (h: CreateElement, params: any) => {
                return h("div", [
                    h("Button", {
                        props: {
                            size: "small",
                            type: this.getStartOrStopType(params.index),
                            icon: this.getStartOrStopIcon(params.index)
                        },
                        style: {
                            marginRight: "10px"
                        },
                        on: {
                            click: () => {
                                this.onStartOrStop(params.index);
                            }
                        }
                    }),
                    h("Button", {
                        props: {
                            size: "small",
                            icon: "edit",
                            type: "info",
                            disabled: this.getDisable(params.index)
                        },
                        style: {
                            marginRight: "10px"
                        },
                        on: {
                            click: () => {
                                this.onEdit(params.index);
                            }
                        }
                    }),
                    h(
                        "Poptip",
                        {
                            props: {
                                confirm: true,
                                placement: "left",
                                title: "您确认要删除这条内容吗？"
                            },
                            on: {
                                "on-ok": () => {
                                    this.onDel(params.index);
                                }
                            }
                        },
                        [
                            h("Button", {
                                props: {
                                    size: "small",
                                    icon: "close-round",
                                    type: "error",
                                    disabled: this.getDisable(params.index)
                                }
                            })
                        ]
                    )
                ]);
            }
        }
    ];

    mounted() {
        this.onLoad();
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
        pms.append("pageIndex", pi.toString());
        pms.append("pageSize", ps.toString());

        if (this.query.region) {
            pms.append("region", this.query.region)
        }

        if (this.query.sourcePort) {
            pms.append("sourcePort", this.query.sourcePort)
        }

        if (this.query.targetIP) {
            pms.append("targetIP", this.query.targetIP)
        }

        if (this.query.targetPort) {
            pms.append("targetPort", this.query.targetPort)
        }

        if (this.query.sortField) {
            pms.append("sortField", this.query.sortField)
        }

        if (this.query.sortDesc) {
            pms.append("sortDesc", this.query.sortDesc)
        }

        let result = await API.get<Model.PortMaps>(
            "/portmap/list/?" + pms.toString()
        );
        if (result.code === 10000) {
            this.tableData = result.data!.list;
        } else {
            this.$Message.error({
                duration: 5,
                content: result.msg + "(" + result.code.toString() + ")"
            });
        }
    }

    getStartOrStopType(index: number): string {
        let typ = "ghost";
        let state = this.tableData.items[index].state;
        switch (state) {
            case "start":
                typ = "warning";
                break;
            case "stop":
                typ = "success";
                break;
            case "startwait":
            case "stopwait":
                typ = "dashed";
                break;
        }

        return typ;
    }

    getStartOrStopIcon(index: number): string {
        let icon = "more";
        let state = this.tableData.items[index].state;
        switch (state) {
            case "start":
            case "startwait":
                icon = "stop";
                break;
            case "stop":
            case "stopwait":
                icon = "play";
                break;
        }

        return icon;
    }

    getDisable(index: number): boolean {
        let isDisable = this.tableData.items[index].state !== "stop";
        return isDisable;
    }

    onKeyDown(event: any) {
        if (_.isNaN(_.parseInt(event.key))) {
            event.returnValue = false;
        }
    }

    async onStartOrStop(index: number) {
        let state = this.tableData.items[index].state;
        let str = "";
        let waitState = "";
        let newState = "";
        if (state === "start") {
            str = "stop";
            waitState = "startwait";
            newState = "stop";
        } else if (state === "stop") {
            str = "start";
            waitState = "stopwait";
            newState = "start";
        }
        if (str === "") {
            return;
        }

        this.tableData.items[index].state = <any>waitState;
        let result = await API.get<Model.PortMaps>(
            "/proxy/" + str + "/" + this.tableData.items[index].id.toString()
        );
        if (result.code === 10000) {
            this.tableData.items[index].state = <any>newState;
        } else {
            this.tableData.items[index].state = state;
            this.$Message.error({
                duration: 5,
                content: result.msg + "(" + result.code.toString() + ")"
            });
        }
    }

    onEdit(index: number) {
        let data = this.tableData.items[index];
        if (data) {
            this.onModalShow(data);
        }
    }

    async onDel(index: number) {
        let id: number = this.tableData.items[index].id;

        let result = await API.delete<any>("/portmap/" + id.toString());
        if (result.code === 10000) {
            this.onLoad();
        } else {
            this.$Message.error({
                duration: 5,
                content: result.msg + "(" + result.code.toString() + ")"
            });
        }
    }

    onModalShow(data?: Model.PortMap) {
        this.modal.isShow = true;
        this.modal.isLoading = false;
        this.modal.isErr = false;
        this.modal.errMsg = "";

        if (data) {
            this.modal.isEdit = true;
            this.modal.title = "端口映射修改";
            this.modal.okBtnText = "修改";

            this.modal.data = <modalDataModel>_.pick(data, _.keys(this.modal.data));
        } else {
            this.modal.isEdit = false;
            this.modal.title = "端口映射添加";
            this.modal.okBtnText = "添加";

            this.modal.data = {
                id: 0,
                region: "",
                title: "",
                targetIP: "127.0.0.1",
                targetPort: 80,
                sourceIP: "127.0.0.1",
                sourcePort: 80,
                userID: "",
                userName: ""
            };
        }
    }

    async onModalOK() {
        this.modal.isErr = false;
        this.modal.errMsg = "";

        let errMsg: string = "";
        if (this.modal.data.sourceIP == "") {
            errMsg = "源IP不能为空";
        } else if (
            this.modal.data.sourcePort <= 0 ||
            this.modal.data.sourcePort > 65536
        ) {
            errMsg = "源端口数值错误，不能为0，最大不能超过65536";
        } else if (this.modal.data.targetIP == "") {
            errMsg = "目标IP不能为空";
        } else if (
            this.modal.data.targetPort <= 0 ||
            this.modal.data.targetPort > 65536
        ) {
            errMsg = "目标端口数值错误，不能为0，最大不能超过65536";
        } else if (this.modal.data.userName == "") {
            errMsg = "创建者不能为空";
        }
        if (errMsg != "") {
            this.modal.errMsg = errMsg;
            this.modal.isErr = true;
            return;
        }

        this.modal.isLoading = true;
        let id: string = "";
        if (this.modal.data.id > 0) {
            id = "/" + this.modal.data.id.toString();
        }

        let result = await API.put<Model.PortMapAlias>(
            "/portmap" + id,
            this.modal.data
        );
        if (result.code === 10000) {
            if (this.modal.data.id > 0) {
                this.onLoad();
            } else {
                this.tableData!.items.unshift(result.data!.portMap);
                this.tableData!.total += 1;
            }
            this.modal.isShow = false;
        } else {
            this.$Message.error({
                duration: 5,
                content: result.msg + "(" + result.code.toString() + ")"
            });
            this.modal.errMsg = result.msg + "(" + result.code.toString() + ")";
            this.modal.isErr = true;
        }
        this.modal.isLoading = false;
    }

    onModalCancel() {
        this.modal.isShow = false;
    }
}
</script>
