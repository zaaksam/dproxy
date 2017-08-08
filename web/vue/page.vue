<template>
    <Page :total="total" :current="pageIndex" :page-size="pageSize" @on-change="onPageChange" @on-page-size-change="onPageSizeChange" show-total show-sizer>
    </Page>
</template>

<script lang="ts">
import Vue from 'vue'
import VueRouter from 'vue-router'
import { Component, Prop } from 'vue-property-decorator'

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
        this.onLoad()
    }

    onPageSizeChange(ps: number) {
        let url = this.$route.path + '?ps=' + ps.toString() + '&pi=1'
        this.$router.replace(url)
        this.onLoad()
    }

    onLoad() {
        this.$emit('onLoad')
    }
}
</script>
