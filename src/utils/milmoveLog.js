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
export function milmoveLog(level, message) {}
