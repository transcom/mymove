import { BatchOtelExporter } from './batchOtelExporter';

describe('BatchOtelExporter', () => {
  describe('constructor', () => {
    it('should create a BatchOtelExporter with default config', () => {
      const exporter = new BatchOtelExporter('test');
      expect(exporter.app).toEqual('test');
      exporter.shutdown();
    });

    it('should create a BatchOtelExporter with custom config', () => {
      const exporter = new BatchOtelExporter('test', {
        endpoint: '/endpoint',
        maxQueueSize: 40,
        maxExportBatchSize: 20,
        scheduledDelayMillis: 2000,
        exportTimeoutMillis: 20000,
      });
      expect(exporter.app).toEqual('test');
      expect(exporter.endpoint).toEqual('/endpoint');
      exporter.shutdown();
    });
  });

  describe('.shutdown/.forceFlush', () => {
    let exporter;
    let sendBatchMock;

    beforeEach(() => {
      exporter = new BatchOtelExporter('test');
      sendBatchMock = jest.spyOn(exporter, 'sendBatchToServer').mockImplementation(async () => {});
    });

    afterEach(() => {
      jest.restoreAllMocks();
      exporter.shutdown();
    });

    it('should do nothing after shutdown', async () => {
      await exporter.shutdown();
      exporter.addToBuffer('debug', 'some', 'arg');
      expect(sendBatchMock).not.toHaveBeenCalled();
    });

    it('should send batch on forceFlush', async () => {
      const level = 'debug';
      const args = ['some', 'arg', 1];
      exporter.addToBuffer(level, args);
      await exporter.forceFlush();
      expect(sendBatchMock).toBeCalled();
      expect(sendBatchMock).toHaveBeenCalledWith([{ level, args }]);
    });
  });

  describe('sendBatchToServer', () => {
    let exporter;
    let fetchMock;

    beforeEach(() => {
      exporter = new BatchOtelExporter('test');
      fetchMock = jest.spyOn(global, 'fetch').mockImplementation(async () => {});
      jest.useFakeTimers();
    });

    afterEach(() => {
      jest.restoreAllMocks();
      exporter.shutdown();
    });

    it('should send batch after timer', async () => {
      const level = 'debug';
      const args = ['some', 'arg', 1];
      exporter.addToBuffer(level, args);
      expect(fetchMock).not.toBeCalled();

      jest.runAllTimers();
      expect(fetchMock).toBeCalled();
      expect(fetchMock).toHaveBeenCalledWith(
        exporter.endpoint,
        expect.objectContaining({
          method: 'POST',
          mode: 'same-origin',
          cache: 'no-cache',
          credentials: 'same-origin',
          header: expect.objectContaining({
            'Content-Type': expect.stringContaining('application/vnd.milmoveclientlog'),
          }),
          signal: expect.any(AbortSignal),
          body: expect.stringContaining('logEntries'),
        }),
      );
    });
  });
});
