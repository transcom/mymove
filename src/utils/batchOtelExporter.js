//
// Inspired by https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-sdk-trace-base/src/export/BatchSpanProcessorBase.ts

import { context } from '@opentelemetry/api';
import { BindOnceFuture, suppressTracing, unrefTimer } from '@opentelemetry/core';

/**
 * @typedef {object} OtelLogEntry - log entry
 * @property {string} level - log level
 * @property {string} message - the message
 */

export class BatchOtelExporter {
  /** @type {string} */
  app;

  /** @type {string} */
  endpoint;

  /** @type {number} */
  maxQueueSize;

  /** @type {number} */
  maxExportBatchSize;

  /** @type {number} */
  scheduledDelayMillis;

  /** @type {number} */
  exportTimeoutMillis;

  /** @type {Array.<OtelLogEntry>} */
  logBuffer = [];

  /** @type {NodeJS.Timeout | undefined } */
  timer;

  /** @type {BindOnceFuture<void>} */
  shutdownOnce;

  /** @type {number} */
  droppedLogsCount = 0;

  /** @type {number} */
  failedSendCount = 0;

  /** @type {number} */
  failedTimerCount = 0;

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
    const {
      endpoint = '/client/log',
      maxQueueSize = 20,
      maxExportBatchSize = 10,
      scheduledDelayMillis = 1000,
      exportTimeoutMillis = 10000,
    } = options;

    this.app = app;
    this.endpoint = endpoint;
    this.maxQueueSize = maxQueueSize;
    this.maxExportBatchSize = maxExportBatchSize;
    this.scheduledDelayMillis = scheduledDelayMillis;
    this.exportTimeoutMillis = exportTimeoutMillis;

    this.shutdownOnce = new BindOnceFuture(this.privateShutdown, this);

    if (this.maxExportBatchSize > this.maxQueueSize) {
      this.maxExportBatchSize = this.maxQueueSize;
    }
  }

  /**
   * flush all logs
   *
   * @returns {Promise<void>}
   */
  forceFlush() {
    if (this.shutdownOnce.isCalled) {
      return this.shutdownOnce.promise;
    }
    return this.flushAll();
  }

  /**
   * Ensure the logger is shutdown a single time
   *
   * @returns {Promise<void>}
   */
  shutdown() {
    return this.shutdownOnce.call();
  }

  /**
   * private method for shutting down
   *
   * @returns <Promise<void>}
   */
  async privateShutdown() {
    return Promise.resolve().then(() => {
      return this.flushAll();
    });
  }

  /**
   * Add a log to the buffer.
   * @param {string} level - log level
   * @param {Array} args - args
   */
  addToBuffer(level, args) {
    if (this.logBuffer.length >= this.maxQueueSize) {
      // limit reached, drop log
      this.droppedLogsCount += 1;

      return;
    }

    this.logBuffer.push({ level, args });
    // only start a timer to send the logs if we have logs to send
    this.maybeStartTimer();
  }

  /**
   * Send all logs to the exporter respecting the batch size limit
   * This function is used only on forceFlush or shutdown,
   * for all other cases flush should be used
   *
   * @returns {Promise<void>}
   */
  flushAll() {
    return new Promise((resolve, reject) => {
      const promises = [];
      // calculate number of batches
      const count = Math.ceil(this.logBuffer.length / this.maxExportBatchSize);
      for (let i = 0, j = count; i < j; i += 1) {
        promises.push(this.flushOneBatch());
      }
      Promise.all(promises)
        .then(() => {
          resolve();
        })
        .catch(reject);
    });
  }

  /**
   * send batch of logs to the server
   *
   * @param {Array.<OtelLogEntry>} batch
   */
  async sendBatchToServer(batch) {
    // run fetch with the configured timeout
    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), this.exportTimeoutMillis);
    try {
      const data = {
        app: this.app,
        loggerStats: {
          droppedLogsCount: this.droppedLogsCount,
          failedSendCount: this.failedSendCount,
          failedTimerCount: this.failedTimerCount,
        },
        logEntries: batch,
      };
      await fetch(this.endpoint, {
        method: 'POST',
        mode: 'same-origin',
        cache: 'no-cache',
        credentials: 'same-origin',
        header: {
          // send the version of the data being sent so we can evolve
          // it in the future
          // Use a custom vendor mime type for versioning
          // https://en.wikipedia.org/wiki/Media_type#Vendor_tree
          'Content-Type': 'application/vnd.milmoveclientlog.v1+json',
        },
        body: JSON.stringify(data),
        signal: controller.signal,
      });
      // reset the stats here so they are only reset on successful delivery
      this.droppedLogsCount = 0;
      this.failedSendCount = 0;
      this.failedTimerCount = 0;
    } catch (err) {
      this.failedSendCount += 1;
      throw err;
    } finally {
      // clear the timeout
      clearTimeout(id);
    }
  }

  /**
   *
   * @returns {Promise<void>}
   */
  flushOneBatch() {
    this.clearTimer();
    if (this.logBuffer.length === 0) {
      return Promise.resolve();
    }
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        // don't wait anymore for export, this way the next batch can start
        reject(new Error('Timeout'));
      }, this.exportTimeoutMillis);
      // prevent downstream exporter calls from generating spans
      context.with(suppressTracing(context.active()), () => {
        // Reset the logBuffer here because the next invocations of
        // the _flush method could pass the same finished logs to the
        // exporter if the buffer is cleared outside the execution of
        // this callback.
        const batch = this.logBuffer.splice(0, this.maxExportBatchSize);

        this.sendBatchToServer(batch)
          .then(() => {
            clearTimeout(timer);
            resolve();
          })
          .catch((error) => {
            reject(error ?? new Error('BatchOtelExporter: log export failed'));
          });
      });
    });
  }

  /**
   * start a timer if one has not already been created
   */
  maybeStartTimer() {
    if (this.timer !== undefined) return;
    this.timer = setTimeout(() => {
      this.flushOneBatch()
        .then(() => {
          if (this.logBuffer.length > 0) {
            this.clearTimer();
            this.maybeStartTimer();
          }
        })
        .catch(() => {
          // not much we can do. We can't log an error
          this.failedTimerCount += 1;
        });
    }, this.scheduledDelayMillis);
    unrefTimer(this.timer);
  }

  clearTimer() {
    if (this.timer !== undefined) {
      clearTimeout(this.timer);
      this.timer = undefined;
    }
  }
}

export default BatchOtelExporter;
