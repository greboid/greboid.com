const pluginRss = require('@11ty/eleventy-plugin-rss')
const syntaxHighlight = require('@11ty/eleventy-plugin-syntaxhighlight')
const pageAssets = require('./page-assets')

module.exports = function(config) {
  config.addPlugin(syntaxHighlight)
  config.addPlugin(pluginRss)
  config.addPlugin(pageAssets)
  return {
    markdownTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    dir: {
      input: 'src',
      output: 'dist',
    },
  }
}
