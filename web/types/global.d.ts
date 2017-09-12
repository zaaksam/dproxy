//全局配置变量定义
interface globalConfigStatic {
    appName: string
    appVersion: string
    token: string
    prefixPath: string
}

declare var globalConfig: globalConfigStatic

interface BaseModalModel {
    isShow: boolean
    isLoading: boolean
    isErr: boolean
    errMsg: string
    isEdit: boolean
    title: string
    okBtnText: string
}