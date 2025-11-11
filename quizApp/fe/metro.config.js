const { getDefaultConfig } = require('expo/metro-config');

const config = getDefaultConfig(__dirname);

// GitHub Pages configuration
if (process.env.GITHUB_PAGES === 'true') {
  config.resolver.platforms = ['web', 'native'];
}

module.exports = config;