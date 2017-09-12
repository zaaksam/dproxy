declare namespace Model {
    interface Log {
        type: string
        content: string
        created: number
    }

    interface Logs {
        list: APIListModel<Log>
    }
}