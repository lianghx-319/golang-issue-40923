import initGoWasm from '../assets/main.wasm'
import { ref, reactive } from 'vue'


export function useGoWasm() {
  const loading = ref(true)
  const ret = reactive({})

  initGoWasm().then((exports) => {
    loading.value = false
    Object.keys(exports).forEach(key => {
      if (typeof exports[key] === 'function') {
        ret[key] = (...args) => new Promise((resolve, reject) => {
          exports[key](args, (error, result) => {
            if (error) {
              reject(error)
            }
            resolve(resolve)
          })
        })
      } else {
        ret[key] = exports[key]
      }
    })
  })

  return {
    loading,
    exports: ret
  }
}