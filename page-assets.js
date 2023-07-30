const path = require('path')
const fs = require('fs')
const crypto = require('crypto')

let pluginOptions = {
  hash: true,
}

async function moveAsset(inputPath, outputPath, assetName) {
  const inputFile = path.parse(path.join(path.dirname(inputPath), assetName))
  const outputDir = path.parse(outputPath).dir
  if (inputFile.ext === "" ) {
    throw new Error("Cannot have a blank link to an asset or image")
  }
  let outputName
  if (pluginOptions.hash) {
    const fileContents = fs.readFileSync(path.join(inputFile.dir, inputFile.base))
    const hash = crypto.createHash('sha1').update(fileContents).digest('hex')
    outputName = `${inputFile.name}-${hash}${inputFile.ext}`
  } else {
    outputName = `${inputFile.name}${inputFile.ext}`
  }
  fs.mkdirSync(outputDir, {recursive: true})
  await fs.promises.copyFile(path.join(inputFile.dir, inputFile.base), path.join(outputDir, outputName))
  return outputName
}

// export plugin
module.exports = {
  configFunction(config, options) {
    Object.assign(pluginOptions, options)
    config.addAsyncShortcode("link", async function(assetName, linkText){
      let name = await moveAsset(this.page.inputPath, this.page.outputPath, assetName)
      return `<a href="${this.page.url}${name}">${linkText}</a>`
    })
    config.addAsyncShortcode("image", async function(assetName){
      let name = await moveAsset(this.page.inputPath, this.page.outputPath, assetName)
      return `<img alt="" src="${this.page.url}${name}"/>`
    })
  },
}
