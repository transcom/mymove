import { BatchOtelExporter } from './batchOtelExporter';

// As for 2023-03-03 opentelemetry logging is still in development
// https://opentelemetry.io/docs/instrumentation/js/#status-and-releases
//
// Implement a logger that will upload
export class OtelLogger {
  /**
   * create a logger that sends information to the server periodically
   * and is compatibile with the @opentelemetry/api DiagLogger
   *
   * @param {string} app
   * @param {Object} options
   * @param {string} [options.endpoint=/client/log] The endpoint to send
   *   logs to
   * @param {number} [options.maxQueueSize=20] The maximum queue size.
   *   After the size is reached logs are dropped
   * @param {number} [options.maxExportBatchSize=10] The maximum batch
   *   size of every export. It must be smaller or equal to maxQueueSize
   * @param {number} [options.scheduledDelayMillis=1000] The interval
   *   between two consecutive imports
   * @param {number} [options.exportTimeoutMillis=10000] How log the
   *   export can run before it is cancelled
   * @returns {import{'@opentelemetry/api').DiagLogger}
   */
  constructor(app, options = {}) {
    // keep exporter hidden so the only exposed API are the log functions
    this.exporter = new BatchOtelExporter(app, options);
  }

  /**
   * log `args` at debug level
   *
   * @param {Array} args
   */
  debug(...args) {
    this.exporter.addToBuffer('debug', args);
  }

  /**
   * log `args` at error level
   *
   * @param {Array} args
   */
  error(...args) {
    return this.exporter.addToBuffer('error', args);
  }

  /**
   * log `args` at info level
   *
   * @param {Array} args
   */
  info(...args) {
    return this.exporter.addToBuffer('info', args);
  }

  /**
   * log `args` at warn level
   *
   * @param {Array} args
   */
  warn(...args) {
    return this.exporter.addToBuffer('warn', args);
  }

  /**
   * log `args` at verbose level
   *
   * @param {Array} args
   */
  verbose(...args) {
    return this.exporter.addToBuffer('verbose', args);
  }
}

export default OtelLogger;
