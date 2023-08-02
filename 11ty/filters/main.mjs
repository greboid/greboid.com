import {renderMarkdown} from './markdown.mjs'
import {niceDate} from './nicedate.mjs'

export const addFilters = function(config) {
  config.addFilter("markdown", renderMarkdown)
  config.addFilter("nicedate", niceDate)
}
