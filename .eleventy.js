const pluginRss = require('@11ty/eleventy-plugin-rss')
const syntaxHighlight = require('@11ty/eleventy-plugin-syntaxhighlight')

module.exports = function(config) {
  config.addPlugin(syntaxHighlight)
  config.addPlugin(pluginRss)

  config.setServerPassthroughCopyBehavior('passthrough')

  return {
    markdownTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    dir: {
      input: 'src',
      output: 'dist',
    },
  }
}
