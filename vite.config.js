const goWasmPlugin = require('./serverPlugin/goWasmPlugin')
module.exports = {
  configureServer: [goWasmPlugin]
}