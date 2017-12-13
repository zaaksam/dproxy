<template>
    <div>
        <Row type="flex" style="margin-top: 10px; margin-bottom:15px">
            类型：
            <Select v-model="query.type" style="width: 120px">
                <Option value="All"></Option>
                <Option value="Info"></Option>
                <Option value="Debug"></Option>
                <Option value="Error"></Option>
                <Option value="Warning"></Option>
                <Option value="Critical"></Option>
            </Select>
            &nbsp;&nbsp;&nbsp;&nbsp;指定时间：
            <Date-picker type="datetime" placement="bottom" placeholder="删除指定时间之前数据" style="width: 170px" @on-change="onDatetimeChange"></Date-picker>
            &nbsp;&nbsp;&nbsp;&nbsp;内容：
            <Input v-model="query.content" placeholder="模糊关键字" style="width: 160px"></Input>
            &nbsp;&nbsp;&nbsp;&nbsp;
            <Button type="primary" @click="onDel">删除指定</Button>
        </Row>
        <Table stripe border :columns="tableColumns" :data="tableData.items">
        </Table>
        <MyPage :total="tableData.total" :pageIndex="tableData.pageIndex" :pageSize="tableData.pageSize" @onChange="onLoad"></MyPage>
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
export default class MyLog extends Vue {
  query: Model.Log = {
    type: "All",
    content: "",
    created: 0
  };

  tableData: APIListModel<Model.Log> = {
    pageIndex: 1,
    pageSize: 10,
    pageCount: 0,
    total: 0,
    items: []
  };

  tableColumns: any[] = <any>[
    {
      title: "类型",
      key: "type"
    },
    {
      title: "内容",
      key: "content"
    },
    {
      title: "记录时间",
      key: "created",
      render: (h: CreateElement, params: any) => {
        return h(
          "span",
          moment.unix(params.row.created).format("YYYY-MM-DD HH:mm:ss")
        );
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
      pi = this.tableData.pageIndex;
    }
    if (_.isNaN(ps)) {
      ps = this.tableData.pageSize;
    }

    let pms = new URLSearchParams();
    pms.set("pageIndex", pi.toString());
    pms.set("pageSize", ps.toString());

    let result = await API.get<Model.Logs>("/log/list/?" + pms.toString());
    if (result.code === 10000) {
      this.tableData = result.data!.list;
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + "(" + result.code.toString() + ")"
      });
    }
  }

  async onDel() {
    if (this.query.created <= 0) {
      this.$Message.error({ duration: 5, content: "指定时间必须选择" });
      return;
    }

    let pms = new URLSearchParams();
    pms.set("created", this.query.created.toString());

    if (this.query.type != "All") {
      pms.set("type", this.query.type);
    }
    if (this.query.content != "") {
      pms.set("content", this.query.content);
    }

    let result = await API.delete<any>("/log/?" + pms.toString());
    if (result.code === 10000) {
      this.onLoad();
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + "(" + result.code.toString() + ")"
      });
    }
  }

  onDatetimeChange(expiredStr: string) {
    let unix = moment(expiredStr, "YYYY-MM-DD HH:mm:ss").unix();
    if (_.isNaN(unix)) {
      this.query.created = 0;
    } else {
      this.query.created = unix;
    }
  }
}
</script>
