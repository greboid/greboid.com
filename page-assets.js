const path = require('path')
const fs = require('fs')
const crypto = require('crypto')

/**
 *
 * @param {string} inputPath Input file path
 * @param {string} outputPath Output file path
 * @param {string} assetName Asset filename
 * @returns {string} New filename
 */
function moveAsset(inputPath, outputPath, assetName) {
  const inputFile = path.parse(path.join(path.dirname(inputPath), assetName))
  const outputDir = path.parse(outputPath).dir
  if (inputFile.ext === '') {
    throw new Error('Cannot have a blank link to an asset or image')
  }
  const fileContents = fs.readFileSync(path.join(inputFile.dir, inputFile.base))
  const hash = crypto.createHash('sha1').update(fileContents).digest('hex')
  const outputName = `${inputFile.name}-${hash}${inputFile.ext}`
  fs.mkdirSync(outputDir, {recursive: true})
  fs.copyFile(path.join(inputFile.dir, inputFile.base), path.join(outputDir, outputName), () => {})
  return outputName
}

/**
 * @param {string} inputPath Input file path
 * @param {string} outputPath Output file path
 * @param {string} assetName Asset filename
 * @returns {Map<string, string>} Map of type to filename
 */
function findOtherImages(inputPath, outputPath, assetName) {
  const assetFile = path.parse(path.join(path.dirname(inputPath), assetName))
  let assets = new Map()
  Array.from(['avif', 'webp']).forEach(extension => {
    const newFile = path.join('./', assetFile.dir, `${assetFile.name}.${extension}`)
    if (fs.existsSync(newFile)) {
      let name = moveAsset(inputPath, outputPath, `${assetFile.name}.${extension}`)
      assets.set(extension, name)
    }
  })
  return assets
}

// export plugin
module.exports = {
  configFunction(config, _) {
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
  },
}
