const pluginRss = require('@11ty/eleventy-plugin-rss')
const syntaxHighlight = require('@11ty/eleventy-plugin-syntaxhighlight')
const pluginRev = require("eleventy-plugin-rev")
const pluginSass = require("eleventy-sass")
const pluginShortcodes = require('./shortcodes')
const pluginFilters = require('./filters')

const configFunction = (config) => {
  config.addPlugin(syntaxHighlight)
  config.addPlugin(pluginRss)
  config.addPlugin(pluginRev)
  config.addPlugin(pluginSass, {
    rev: true
  })
  config.addPlugin(pluginShortcodes)
  config.addPlugin(pluginFilters)
  config.setFrontMatterParsingOptions({
    excerpt: true,
    excerpt_separator: "<!--more-->"
  });
  config.addPassthroughCopy({
    "./src/images/": "/images/"
  })
  config.addWatchTarget("./src/images/*")
  return {
    markdownTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    dir: {
      input: 'src',
      output: 'dist',
    },
  }
}

module.exports = configFunction
