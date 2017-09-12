declare namespace Model {
    interface PortMap {
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

    interface PortMapAlias {
        portMap: PortMap
    }

    interface PortMaps {
        list: APIListModel<PortMap>
    }
}