const {moveAsset, findOtherImages} = require('./assets.js')

const configFunction = (config, _) => {
  config.addAsyncShortcode('link', async function(assetName, linkText) {
    let name = moveAsset(this.page.inputPath, this.page.outputPath, assetName)
    return `<a href="${this.page.url}${name}">${linkText}</a>`
  })
  config.addAsyncShortcode('image', async function(assetName) {
    let outputTag = '<picture>\n'
    findOtherImages(this.page.inputPath, this.page.outputPath, assetName).forEach((value, key) => {
      outputTag += `<source srcset="${this.page.url}${value}" type="image/${key}">\n`
    })
    let name = moveAsset(this.page.inputPath, this.page.outputPath, assetName)
    outputTag += `<img src="${this.page.url}${name}" alt="" loading="lazy">`
    outputTag += '</picture>\n'
    return outputTag
  })
}

module.exports = configFunction
