<template>
    <Page class="myPage" :total="total" :current="pageIndex" :page-size="pageSize" @on-change="onPageChange" @on-page-size-change="onPageSizeChange" placement="top" show-total show-sizer>
    </Page>
</template>

<style>
.myPage {
    height: 60px;
    line-height: 60px;
}
</style>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import VueRouter from 'vue-router'

@Component
export default class MyPage extends Vue {
    @Prop()
    pageIndex: number = 1
    @Prop()
    pageSize: number = 10
    @Prop()
    total: number = 0

    onPageChange(pi: number) {
        let url = this.$route.path + '?'
        if (this.$route.query.ps) {
            url += 'ps=' + this.$route.query.ps + '&'
        }
        url += 'pi=' + pi.toString()
        this.$router.replace(url)
        this.onChange()
    }

    onPageSizeChange(ps: number) {
        let url = this.$route.path + '?ps=' + ps.toString() + '&pi=1'
        this.$router.replace(url)
        this.onChange()
    }

    onChange() {
        this.$emit('onChange')
    }
}
</script>
