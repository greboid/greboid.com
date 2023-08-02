import {createPicture} from '../libs/createPicture.mjs'

export const image = async function(assetName, width, alt) {
  return createPicture(this.page.url, this.page.outputPath, this.page.inputPath, assetName, width, alt)
}
