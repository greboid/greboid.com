export const addAssets = function(config) {
  config.addPassthroughCopy({
    "./src/images/": "/images/"
  })
  config.addWatchTarget("./src/images/*")
}
