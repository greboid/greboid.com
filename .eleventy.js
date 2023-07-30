const pluginRss = require('@11ty/eleventy-plugin-rss')
const syntaxHighlight = require('@11ty/eleventy-plugin-syntaxhighlight')
const pageAssets = require('./page-assets')
const pluginRev = require("eleventy-plugin-rev")
const pluginSass = require("eleventy-sass")
const markdownIt = require("markdown-it");
const {DateTime} = require("luxon");

module.exports = function(config) {
  const md = new markdownIt({
    html: true,
  })
  config.addPlugin(syntaxHighlight)
  config.addPlugin(pluginRss)
  config.addPlugin(pluginRev)
  config.addPlugin(pluginSass, { rev: true })
  config.addPlugin(pageAssets)
  config.addFilter("markdown", (content) => {
    return md.render(content);
  })
  config.addFilter("niceDate", (date) => { return DateTime.fromJSDate(date).toLocaleString(DateTime.DATE_MED) })
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
