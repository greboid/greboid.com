import {createHash} from 'crypto'
import {mkdirSync, copyFile, readFileSync} from 'fs'
import {parse, join, dirname} from 'path'

/**
 * @param {string} inputPath Input file path
 * @param {string} outputPath Output file path
 * @param {string} assetName Asset filename
 * @returns {string} New filename
 */
export function moveAsset(inputPath, outputPath, assetName) {
  const inputFile = parse(join(dirname(inputPath), assetName))
  const outputDir = parse(outputPath).dir
  if (inputFile.ext === '') {
    throw new Error('Cannot have a blank link to an asset or image')
  }
  const fileContents = readFileSync(join(inputFile.dir, inputFile.base))
  const hash = createHash('sha1').update(fileContents).digest('hex')
  const outputName = `${inputFile.name}-${hash}${inputFile.ext}`
  mkdirSync(outputDir, {recursive: true})
  copyFile(join(inputFile.dir, inputFile.base), join(outputDir, outputName), () => {})
  return outputName
}
