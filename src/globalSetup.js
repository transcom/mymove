// globalSetup.js
const dotenv = require('dotenv');

async function globalSetup() {
  dotenv.config({
    path: '.envrc',
    override: false,
  });
}

module.exports = globalSetup;
