import {moveAsset} from '../libs/moveAsset.mjs'

export const link = async function(assetName, linkText) {
  let name = moveAsset(this.page.inputPath, this.page.outputPath, assetName)
  return `<a href="${this.page.url}${name}">${linkText}</a>`
}
