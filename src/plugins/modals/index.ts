import { reactive } from 'vue'
import type { App, Plugin } from 'vue'
import Modal from './Modal.vue'

const _current = reactive({
  name: '',
  resolve: null as any,
  reject: null as any
})
const api = {
  active() {
    return _current.name
  },
  show(name: string) {
    _current.name = name
    return new Promise((resolve = () => {}, reject = () => {}) => {
      _current.resolve = resolve
      _current.reject = reject
    })
  },
  accept(...args: any[]) {
    _current.resolve(...args)
    _current.name = ''
  },
  reject(...args: any[]) {
    _current.reject(...args)
    _current.name = ''
  }
}

const plugin: Plugin = {
  install(app: App, options: any) {
    app.component('Modal', Modal)
    app.provide('$modals', api)
  }
}

export default plugin
