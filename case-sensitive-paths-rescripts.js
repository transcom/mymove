/*  import/no-extraneous-dependencies */
const { appendWebpackPlugin } = require('@rescripts/utilities');
const CaseSensitivePathsPlugin = require('case-sensitive-paths-webpack-plugin');
/* eslint-enable import/no-extraneous-dependencies */

module.exports = appendWebpackPlugin(new CaseSensitivePathsPlugin());
