declare namespace Model {
    interface PortMapQuery {
        region?: string
        targetIP?: string
        targetPort?: string
        sourcePort?: string
        sortField: 'created' | 'sourcePort'
        sortDesc: '1' | '0'
    }

    interface PortMap {
        id: number
        region: string
        title: string
        targetIP: string
        targetPort: number
        sourceIP: string
        sourcePort: number
        userID: string
        userName: string
        created: number
        updated: number
        state: 'start' | 'startwait' | 'stop' | 'stopwait'
    }

    interface PortMapAlias {
        portMap: PortMap
    }

    interface PortMaps {
        list: APIListModel<PortMap>
    }
}