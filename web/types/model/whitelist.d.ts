declare namespace Model {
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