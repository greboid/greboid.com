const {moveAsset} = require('./assets.js')
const {thumbnail, picture} = require('./thumbnailer')

const configFunction = (config, _) => {
  config.addAsyncShortcode('link', async function(assetName, linkText) {
    let name = moveAsset(this.page.inputPath, this.page.outputPath, assetName)
    return `<a href="${this.page.url}${name}">${linkText}</a>`
  })
  config.addAsyncShortcode('image', async function(assetName, width, alt) {
    return picture(this.page.url, this.page.outputPath, this.page.inputPath, assetName, width, alt)
  })
  config.addAsyncShortcode('thumbnail', async function(assetName, width, alt) {
    return thumbnail(this.page.url, this.page.outputPath, this.page.inputPath, assetName, width, alt)
  })
}

module.exports = configFunction
