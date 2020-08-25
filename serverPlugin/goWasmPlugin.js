const { isImportRequest } = require('vite')

module.exports = ({ app, config }) => {
  app.use((ctx, next) => {
    if (ctx.path.endsWith('.wasm') && isImportRequest(ctx)) {
      ctx.type = 'js'
      ctx.body = `import './wasm_exec.js'
      const go = new Go()
      Go.exports = {}
      export default (opts = go.importObject) => {
        return WebAssembly.instantiateStreaming(fetch(${JSON.stringify(
          ctx.path
        )}), opts)
          .then(obj => {
            go.run(obj.instance)
            return Go.exports
          })
      }`
      return
    }
    return next()
  })
}