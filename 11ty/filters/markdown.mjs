import markdownIt from 'markdown-it';

const md = new markdownIt({
  html: true,
})
export const renderMarkdown = function(content) {
  return md.render(content);
}
