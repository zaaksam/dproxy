declare module '*.vue' {
    import Vue from 'vue'
    export default typeof Vue
}

declare module 'iview' {
    const iview: any
    export default iview
}

declare interface BaseModalModel {
    isShow: boolean
    isLoading: boolean
    isErr: boolean
    errMsg: string
}

declare interface BaseListModel {
    total: number
    pageIndex: number
    pageSize: number
    pageCount: number
}

declare interface WhiteListModel {
    id: number
    ip: string
    userID: string
    userName: string
    expired: number
    created: number
    updated: number
}

declare interface PortMapModel {
    id: number
    title: string
    targetIP: string
    targetPort: number
    sourceIP: string
    sourcePort: number
    userID: string
    userName: string
    created: number
    updated: number
    isStart: boolean
}

declare interface LogModel {
    type: string
    content: string
    created: number
}

//全局配置变量定义
interface globalConfigStatic {
    appName: string
    appVersion: string
    token: string
    prefixPath: string
}

declare var globalConfig: globalConfigStatic