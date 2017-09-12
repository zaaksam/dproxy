import Axios, { AxiosRequestConfig } from 'axios'

interface APIStatic {
    /**
     * get请求
     */
    get<T>(url: string): Promise<APIResult<T>>
    /**
     * delete请求
     */
    delete<T>(url: string): Promise<APIResult<T>>
    /**
     * put请求
     */
    put<T>(url: string, data?: any): Promise<APIResult<T>>
    /**
     * post请求
     */
    post<T>(url: string, data?: URLSearchParams): Promise<APIResult<T>>
}

async function request<T>(method: 'get' | 'delete' | 'post' | 'put', url: string, data?: any) {
    let result: APIResult<T> = { code: -1, msg: '' }

    if (globalConfig.token != '') {
        url += url.indexOf('?') == -1 ? '?' : '&'
        url += 'token=' + globalConfig.token
    }

    let conf: AxiosRequestConfig = {
        method: method,
        url: '/api' + url
    }
    if (data) {
        conf.data = data
    }

    try {
        let res = await Axios.request(conf)
        result.code = res.data.code
        result.msg = res.data.msg
        result.data = res.data.data

    } catch (err) {
        result.msg = err.message
    }

    return result
}

const API: APIStatic = {
    get: <T>(url: string) => {
        return request<T>('get', url)
    },
    delete: <T>(url: string) => {
        return request<T>('delete', url)
    },
    put: <T>(url: string, data?: T) => {
        return request<T>('put', url, data)
    },
    post: <T>(url: string, data?: URLSearchParams) => {
        return request<T>('post', url, data)
    }
}

export default API