import { OtelLogger } from './otelLogger';

describe('OtelLogger', () => {
  describe('constructor', () => {
    it('should create a OtelLogger with default config', () => {
      const logger = new OtelLogger('test');
      expect(logger.exporter.app).toEqual('test');
      logger.exporter.shutdown();
    });

    it('should create a OtelLogger with custom config', () => {
      const logger = new OtelLogger('test', {
        endpoint: '/endpoint',
        maxQueueSize: 40,
        maxExportBatchSize: 20,
        scheduledDelayMillis: 2000,
        exportTimeoutMillis: 20000,
      });
      expect(logger.exporter.app).toEqual('test');
      expect(logger.exporter.endpoint).toEqual('/endpoint');
      logger.exporter.shutdown();
    });
  });

  describe('log methods', () => {
    let logger;
    let addBufferMock;

    beforeEach(() => {
      logger = new OtelLogger('test');
      addBufferMock = jest.spyOn(logger.exporter, 'addToBuffer').mockImplementation(() => {});
    });

    afterEach(() => {
      jest.restoreAllMocks();
      logger.exporter.shutdown();
    });

    const methods = ['debug', 'error', 'info', 'warn', 'verbose'];

    it.each(methods)('should addToBuffer on "%s"', (methodName) => {
      const args = ['arg', 1];
      logger[methodName](...args);
      expect(addBufferMock).toBeCalled();
    });
  });
});
