import {link} from './link.mjs'
import {image} from './image.mjs'
import {thumbnail} from './thumbnail.mjs'

export const addShortcodes = function (config) {
  config.addShortcode("link", link)
  config.addShortcode("thumbnail", thumbnail)
  config.addShortcode("image", image)
}
