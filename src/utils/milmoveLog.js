// log levels to use with milmoveLog
export const MILMOVE_LOG_LEVEL = {
  ERROR: 'ERROR',
  WARN: 'WARN',
  INFO: 'INFO',
  DEBUG: 'DEBUG',
  LOG: 'LOG',
};

// milmoveLog: wrapper log to re-direct logging to something other than console.log/warn/error
// use the levels set by MILMOVE_LOG_LEVEL
//
// RA Summary: eslint - no-unused-vars - Should not have unused variables
// RA: The method parameters 'level' and 'message' are not used within the method
// RA: This method is a place holder so it is a no-op There is an ADR in the works to figure out what the final solution will be.
// RA: https://github.com/transcom/mymove/pull/7007
// RA: for now removing all use of console.warn console.error console.debug console.log console.info
// RA: this new function to be used re-direct logs to the final destination
// RA Developer Status: Known Issue
// RA Validator Status: Known Issue
// RA Validator: leodis.f.scott.civ@mail.mil
// RA Modified Severity: CAT III
// eslint-disable-next-line no-unused-vars
export function milmoveLog(level, message) {}
