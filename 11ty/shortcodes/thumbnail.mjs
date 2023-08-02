import {createThumbnail} from '../libs/createThumbnail.mjs'

export const thumbnail = async function(assetName, width, alt) {
  return createThumbnail(this.page.url, this.page.outputPath, this.page.inputPath, assetName, width, alt)
}
