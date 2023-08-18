// use the diag logger from opentelemetry
// This uses the logger configured via `diag.setLogger`, otherwise
// calls to diag are noops
import { diag, DiagConsoleLogger, DiagLogLevel } from '@opentelemetry/api';

import { OtelLogger } from 'utils/otelLogger';

/**
 * configure the logger
 * @param {string} app - my, office, or admin
 * @param {object} options
 * @param {string} [options.loggingType=none]
 * @param {string} [options.loggingLevel=unknown]
 */
export function configureLogger(app, options = {}) {
  const { loggingType = 'none', loggingLevel = 'unknown' } = options;

  if (loggingType === 'otel') {
    const logger = new OtelLogger(app);
    const logLevel = DiagLogLevel[loggingLevel] || DiagLogLevel.ERROR;
    diag.setLogger(logger, { logLevel });
  } else if (loggingType === 'console') {
    const logger = new DiagConsoleLogger();
    const logLevel = DiagLogLevel[loggingLevel] || DiagLogLevel.VERBOSE;
    diag.setLogger(logger, { logLevel });
  }
}

// Use the opentelemetry provided logger. If it has not been
// configured with `diag.setLogger`, nothing will be logged
export const milmoveLogger = diag;
