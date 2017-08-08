declare module '*.vue' {
    import Vue from 'vue'
    export default typeof Vue
}

declare module 'iview' {
    const iview: any
    export default iview
}