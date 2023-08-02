import {addShortcodes} from './shortcodes/main.mjs'
import {addFilters} from './filters/main.mjs'
import {addPlugins} from './plugins.mjs'
import {addAssets} from './assets.mjs'
import {updateFrontMatter} from './frontmatter.mjs'
import {returnSettings} from './config.mjs'
import pkg from '@11ty/eleventy'

const UserConfig = pkg

/** @param {UserConfig} config */
export const install = function (config) {
  addFilters(config)
  addShortcodes(config)
  addPlugins(config)
  addAssets(config)
  updateFrontMatter(config)
  return returnSettings(config)
}

export default install
