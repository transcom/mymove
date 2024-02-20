// globalSetup.js
/* eslint-disable @typescript-eslint/no-var-requires */
const dotenv = require('dotenv');

async function globalSetup() {
  dotenv.config({
    path: '.envrc',
    override: false,
  });
}

module.exports = globalSetup;
