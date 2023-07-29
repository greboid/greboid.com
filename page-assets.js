const path = require('path')
const fs = require('fs')
const {JSDOM} = require('jsdom')
const crypto = require('crypto')

let pluginOptions = {
  pagePrefix: ['./src/blog'],
}

async function transformParser(content, outputPath) {
  if (outputPath.endsWith('.html') && pluginOptions.pagePrefix.some(prefix => this.inputPath.startsWith(prefix)) && this.inputPath.endsWith('index.md')) {
      const templateDir = path.dirname(this.inputPath)
      const outputDir = path.dirname(outputPath)
      const dom = new JSDOM(content)
      await handleImages(dom, outputDir, templateDir)
      await handleFiles(dom, outputDir, templateDir)
      content = dom.serialize()
  }
  return content
}

async function handleFiles(dom, outputDir, templateDir) {
  const elms = Array.from(dom.window.document.querySelectorAll('a')).
      filter(elm => !elm.getAttribute('href').startsWith('http') && elm.getAttribute('href') !== ".")
  await Promise.all(elms.map(async (link) => {
    const assetPath = path.join(templateDir, link.getAttribute('href'))
    const file = path.parse(assetPath)
    const fileContents = fs.readFileSync(assetPath)
    const hash = crypto.createHash('sha1').update(fileContents).digest('hex')
    link.setAttribute('href', `${file.name}-${hash}${file.ext}`)
    fs.mkdirSync(outputDir, {recursive: true})
    await fs.promises.copyFile(assetPath, path.join(outputDir, `${file.name}-${hash}${file.ext}`))
  }))
}

async function handleImages(dom, outputDir, templateDir) {
  const elms = Array.from(dom.window.document.querySelectorAll('img')).
      filter(elm => !elm.getAttribute('src').startsWith('http'))
  await Promise.all(elms.map(async (img) => {
    const assetPath = path.join(templateDir, img.getAttribute('src'))
    const file = path.parse(assetPath)
    const fileContents = fs.readFileSync(assetPath)
    const hash = crypto.createHash('sha1').update(fileContents).digest('hex')
    img.setAttribute('integrity', `sha1-${hash}`)
    img.setAttribute('src', `${file.name}-${hash}${file.ext}`)
    fs.mkdirSync(outputDir, {recursive: true})
    await fs.promises.copyFile(assetPath, path.join(outputDir, `${file.name}-${hash}${file.ext}`))
  }))
}

// export plugin
module.exports = {
  configFunction(config, options) {
    Object.assign(pluginOptions, options)
    config.addTransform(`[Page-Assets]-transform-parser`, transformParser)
  },
}
