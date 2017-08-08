export interface BaseModalModel {
    isShow: boolean
    isLoading: boolean
    isErr: boolean
    errMsg: string
}

export interface BaseListModel {
    total: number
    pageIndex: number
    pageSize: number
    pageCount: number
}

export interface WhiteListModel {
    id: number
    ip: string
    userID: string
    userName: string
    expired: number
    created: number
    updated: number
}

export interface PortMapModel {
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

export interface LogModel {
    type: string
    content: string
    created: number
}