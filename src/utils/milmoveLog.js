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
// TODO There is an ADR in the works to figure out what the final solution will be.
// TODO https://github.com/transcom/mymove/pull/7007
// TODO for now removing all use of console.warn console.error console.debug console.log console.info
// TODO this new function to be used re-direct logs to the final destination
// eslint-disable-next-line no-unused-vars
export function milmoveLog(level, message) {}
