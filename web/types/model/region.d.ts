declare namespace Model {
    interface Region {
        name: string
        lastHeartbeat: number
    }

    interface Regions {
        list: APIListModel<Region>
    }
}