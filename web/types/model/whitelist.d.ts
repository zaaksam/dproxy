declare namespace Model {
    interface WhiteListQuery {
        ip?: string
        userName?: string
        isExpired?: '1' | '0'
        sortField: 'created' | 'expired'
        sortDesc: '1' | '0'
    }

    interface WhiteList {
        id: number
        ip: string
        userID: string
        userName: string
        expired: number
        created: number
        updated: number
    }

    interface WhiteListAlias {
        whiteList: WhiteList
    }

    interface WhiteLists {
        list: APIListModel<WhiteList>
    }
}