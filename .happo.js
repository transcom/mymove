const { RemoteBrowserTarget } = require('happo.io');
const happoPluginStorybook = require('happo-plugin-storybook');

require('dotenv').config();

module.exports = {
  apiKey: process.env.HAPPO_API_KEY,
  apiSecret: process.env.HAPPO_API_SECRET,
  targets: {
    chrome: new RemoteBrowserTarget('chrome', {
      viewport: '1024x768',
    }),
    // TODO - IE is failing because Storybook causes syntax error. Need to investigate
    /*
    'internet explorer': new RemoteBrowserTarget('internet explorer', {
      viewport: '1024x768',
    }),
    */
    // TODO - Safari is failing because Storybook causes syntax error. Need to investigate
    // 'ios-safari': new RemoteBrowserTarget('ios-safari', {
    //   viewport: '375x667',
    //   scrollStitch: true,
    // }),
  },
  plugins: [
    happoPluginStorybook({
      outputDir: 'storybook-static',
    }),
  ],
};
