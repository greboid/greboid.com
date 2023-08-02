import path from 'path'
import Image from '@11ty/eleventy-img'

export const createThumbnail = async (url, outputFile, inputFile, assetName, width, alt = "") => {
  if (alt === "") {
    return Promise.reject("alt must not be blank")
  }
  inputFile = path.parse(inputFile)
  inputFile = path.parse(path.join(inputFile.dir, assetName))
  outputFile = path.parse(outputFile)
  const types = ['avif', 'webp', 'jpg']
  const images = await Image(path.join(inputFile.dir, inputFile.base), {
    widths: [width],
    formats: types,
    outputDir: outputFile.dir,
    urlPath: url,
  })
  let output = "<picture>"
  types.forEach(type => {
    if (images[type]) {
      if (type !== "jpg") {
        images[type].forEach(image => {
          output += `<source type="${image.sourceType}" decoding="async" srcset="${image.url}" width="${image.width}" height="${image.height}">`
        })
      }
    }
  })
  output += `<img alt=${alt} src="${images.jpeg[0].url}" decoding="async" width="${images.jpeg[0].width}" height="${images.jpeg[0].height}">`
  output += "</picture>"
  return output
}
