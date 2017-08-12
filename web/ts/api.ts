export interface APIStatic {
    initURL: (url: string) => string
}

const API: APIStatic = {
    initURL: (url: string) => {
        url += url.indexOf('?') == -1 ? '?' : '&'
        url += 'token=' + globalConfig.token

        return url
    }
}

export default API