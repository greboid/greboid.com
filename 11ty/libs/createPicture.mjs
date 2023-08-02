import path from 'path'
import Image from '@11ty/eleventy-img'

export const createPicture = async(url, outputFile, inputFile, assetName, width, alt = "") => {
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
