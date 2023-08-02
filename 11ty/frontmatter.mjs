export const updateFrontMatter = function(config) {
  config.setFrontMatterParsingOptions({
    excerpt: true,
    excerpt_separator: "<!--more-->"
  });
}
