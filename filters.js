const markdownIt = require("markdown-it");
const luxon = require('luxon')


const md = new markdownIt({
  html: true,
})

const renderMarkdown = (content) => {
  return md.render(content);
}

const niceDate = (date) => {
  return luxon.DateTime.fromJSDate(date).toLocaleString(luxon.DateTime.DATE_MED)
}

const configFunction = (config, _) => {
  config.addFilter("markdown", renderMarkdown)
  config.addFilter("niceDate", niceDate)
}

module.exports = configFunction
