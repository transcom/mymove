// globalSetup.js
// Summary: @typescript-eslint/no-var-requires disallows the use of require statements except in import statements. This is for the most part only a syntactical rule. If you don't care about TypeScript module syntax, then you will not need this rule.
// RA: The linter flags any use of require that is not in an import.
// RA Developer Status: N/A
// RA Validator Status: Known Issue
// RA Modified Severity: N/A
/* eslint-disable @typescript-eslint/no-var-requires */
const dotenv = require('dotenv');

async function globalSetup() {
  dotenv.config({
    path: '.envrc',
    override: false,
  });
}

module.exports = globalSetup;
