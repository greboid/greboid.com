import pluginSyntaxHighlight from '@11ty/eleventy-plugin-syntaxhighlight'
import pluginRss  from '@11ty/eleventy-plugin-rss'
import pluginRev from 'eleventy-plugin-rev'
import pluginSass from 'eleventy-sass'

export const addPlugins = function(config) {
  config.addPlugin(pluginSyntaxHighlight)
  config.addPlugin(pluginRss)
  config.addPlugin(pluginRev)
  config.addPlugin(pluginSass, {
    rev: true
  })
}
