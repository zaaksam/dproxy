interface APIResult<T> {
    code: number
    msg: string
    data?: T
}

interface APIListModel<T> {
    pageIndex: number
    pageSize: number
    pageCount: number
    total: number
    items: T[]
}