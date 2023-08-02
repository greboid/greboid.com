const Image = require('@11ty/eleventy-img')
const path = require('path')

const picture = async(url, outputFile, inputFile, assetName, width, alt = "") => {
  if (alt === "") {
    return Promise.reject("alt must not be blank")
  }
  inputFile = path.parse(inputFile)
  inputFile = path.parse(path.join(inputFile.dir, assetName))
  outputFile = path.parse(outputFile)
  const types = ['avif', 'webp', 'jpeg']
  const images = await Image(path.join(inputFile.dir, inputFile.base), {
    widths: [width, 'auto'],
    formats: types,
    outputDir: outputFile.dir,
    urlPath: url,
  })
  const imageAttributes = {
    alt: "",
    sizes: "100vw",
    loading: "lazy",
    decoding: "async",
  }
  return Image.generateHTML(images, imageAttributes);
}

const thumbnail = async (url, outputFile, inputFile, assetName, width, alt = "") => {
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

module.exports = { thumbnail, picture }
