/* eslint-disable import/no-extraneous-dependencies */
const { RemoteBrowserTarget } = require('happo.io');
const happoPluginStorybook = require('happo-plugin-storybook');

require('dotenv').config();

module.exports = {
  apiKey: process.env.HAPPO_API_KEY,
  apiSecret: process.env.HAPPO_API_SECRET,
  targets: {
    chrome: new RemoteBrowserTarget('chrome', {
      viewport: '1024x768',
      project: 'mymove-integration',
    }),
  },
  plugins: [
    happoPluginStorybook({
      outputDir: 'storybook-static',
    }),
  ],
};
