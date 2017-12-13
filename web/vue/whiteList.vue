<template>
    <div>
        <Table stripe border :columns="tableColumns" :data="tableData.items">
        </Table>
        <MyPage :total="tableData.total" :pageIndex="tableData.pageIndex" :pageSize="tableData.pageSize" @onChange="onLoad"></MyPage>
        <Modal :title="modal.title" v-model="modal.isShow" :mask-closable="false" :closable="false">
            <Form :label-width="100">
                <Form-item label="白名单ID" v-show="modal.isEdit">
                    <Input v-model="modal.data.id" readonly></Input>
                </Form-item>
                <Form-item label="申请者ID">
                    <Input v-model="modal.data.userID"></Input>
                </Form-item>
                <Form-item label="申请者名称">
                    <Input v-model="modal.data.userName"></Input>
                </Form-item>
                <Form-item label="申请者IP">
                    <Input v-model="modal.data.ip"></Input>
                </Form-item>
                <Form-item label="过期时间">
                    <Date-picker type="datetime" :value="modal.expiredStr" placement="top" :options="dateOptions" :editable="false" placeholder="放空默认24小时后过期" style="width: 200px" @on-change="onModalDatetimeChange"></Date-picker>
                </Form-item>
            </Form>
            <Alert v-show="modal.isLoading" type="warning" show-icon>数据提交中，请勿关闭此窗口...</Alert>
            <Alert v-show="modal.isErr" type="error" show-icon>{{modal.errMsg}}</Alert>
            <div slot="footer">
                <Button :disabled="modal.isLoading" @click="onModalCancel">取消</Button>
                <Button type="primary" :loading="modal.isLoading" @click="onModalOk">{{modal.okBtnText}}</Button>
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
  ip: string;
  userID: string;
  userName: string;
  expired: number;
}

interface modalModel extends BaseModalModel {
  expiredStr: string;
  data: modalDataModel;
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
    errMsg: "",
    isEdit: false,
    title: "",
    okBtnText: "",
    expiredStr: "",
    data: {
      id: 0,
      ip: "",
      userID: "",
      userName: "",
      expired: 0
    }
  };
  dateOptions: any = {
    disabledDate: (date: Date) => {
      // - 86400000
      return date && date.valueOf() < Date.now();
    }
  };

  tableData: APIListModel<Model.WhiteList> = {
    pageIndex: 1,
    pageSize: 10,
    pageCount: 0,
    total: 0,
    items: []
  };

  tableColumns: any[] = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "IP",
      key: "ip"
    },
    {
      title: "申请人ID",
      key: "userID"
    },
    {
      title: "申请人名称",
      key: "userName"
    },
    {
      title: "申请时间",
      key: "created",
      render: (h: CreateElement, params: any) => {
        return h(
          "span",
          moment.unix(params.row.created).format("YYYY-MM-DD HH:mm:ss")
        );
      }
    },
    {
      title: "过期时间",
      key: "expired",
      render: (h: CreateElement, params: any) => {
        return h(
          "span",
          moment.unix(params.row.expired).format("YYYY-MM-DD HH:mm:ss")
        );
      }
    },
    {
      width: 160,
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
              icon: "edit",
              type: "info"
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
                  type: "error"
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
    let pi = _.parseInt(this.$route.query.pi);
    let ps = _.parseInt(this.$route.query.ps);
    if (_.isNaN(pi)) {
      pi = this.tableData!.pageIndex;
    }
    if (_.isNaN(ps)) {
      ps = this.tableData!.pageSize;
    }

    let pms = new URLSearchParams();
    pms.append("pageIndex", pi.toString());
    pms.append("pageSize", ps.toString());

    let result = await API.get<Model.WhiteLists>(
      "/whitelist/list/?" + pms.toString()
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

  onEdit(index: number) {
    let data = this.tableData!.items[index];
    if (data) {
      this.onModalShow(data);
    }
  }

  async onDel(index: number) {
    let id: number = this.tableData!.items[index].id!;
    if (!id) {
      return;
    }

    let result = await API.delete<any>("/whitelist/" + id.toString());
    if (result.code === 10000) {
      this.onLoad();
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + "(" + result.code.toString() + ")"
      });
    }
  }

  onModalShow(data?: Model.WhiteList) {
    this.modal.isShow = true;
    this.modal.isLoading = false;
    this.modal.isErr = false;
    this.modal.errMsg = "";

    if (data) {
      this.modal.isEdit = true;
      this.modal.title = "白名单修改";
      this.modal.okBtnText = "修改";

      this.modal.data = <modalDataModel>_.pick(data, _.keys(this.modal.data));
      this.modal.expiredStr = moment
        .unix(this.modal.data.expired)
        .format("YYYY-MM-DD HH:mm:ss");
    } else {
      this.modal.isEdit = false;
      this.modal.title = "白名单申请";
      this.modal.okBtnText = "申请";

      this.modal.data = {
        id: 0,
        ip: "",
        userID: "",
        userName: "",
        expired: 0
      };
      // this.modal.expiredDate = moment.unix(this.modal.data.expired).toDate()
      this.modal.expiredStr = "";
    }
  }

  onModalDatetimeChange(expiredStr: string) {
    this.modal.expiredStr = expiredStr;

    let unix = moment(expiredStr, "YYYY-MM-DD HH:mm:ss").unix();
    if (_.isNaN(unix)) {
      this.modal.data.expired = 0;
    } else {
      this.modal.data.expired = unix;
    }
  }

  async onModalOk() {
    this.modal.isErr = false;
    this.modal.errMsg = "";

    let errMsg: string = "";
    if (this.modal.data.ip == "") {
      errMsg = "IP不能为空";
    } else if (this.modal.data.userID == "") {
      errMsg = "申请者ID不能为空";
    } else if (this.modal.data.userName == "") {
      errMsg = "申请者名称不能为空";
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

    let result = await API.put<Model.WhiteListAlias>(
      "/whitelist" + id,
      this.modal.data
    );
    if (result.code === 10000) {
      if (this.modal.data.id > 0) {
        this.onLoad();
      } else {
        this.tableData!.items.unshift(result.data!.whiteList);
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
