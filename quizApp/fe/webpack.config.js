const createExpoWebpackConfigAsync = require('@expo/webpack-config');

module.exports = async function (env, argv) {
  const config = await createExpoWebpackConfigAsync(env, argv);
  
  // Configure for GitHub Pages
  if (process.env.NODE_ENV === 'production') {
    config.output.publicPath = '/quizzie.cloud/';
  }
  
  return config;
};