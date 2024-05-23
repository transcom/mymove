import { diag, DiagConsoleLogger, DiagLogLevel } from '@opentelemetry/api';

import { configureLogger, configureGlobalLogger } from './milmoveLog';
import { OtelLogger } from './otelLogger';

describe('configureLogger', () => {
  let setLoggerMock;

  beforeEach(() => {
    setLoggerMock = jest.spyOn(diag, 'setLogger').mockImplementation(() => {});
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('should configure GlobalLogger', async () => {
    configureGlobalLogger();
  });

  it('should configure OtelLogger', async () => {
    configureLogger('test', { loggingType: 'otel' });
    expect(setLoggerMock).toBeCalled();
    // use mock.calls to use ToBeInstanceOf
    expect(setLoggerMock.mock.calls[0][0]).toBeInstanceOf(OtelLogger);
    expect(setLoggerMock.mock.calls[0][1]).toEqual({ logLevel: DiagLogLevel.ERROR });
  });

  it('should configure DiagConsoleLogger', async () => {
    configureLogger('test', { loggingType: 'console' });
    expect(setLoggerMock).toBeCalled();
    // use mock.calls to use ToBeInstanceOf
    expect(setLoggerMock.mock.calls[0][0]).toBeInstanceOf(DiagConsoleLogger);
    expect(setLoggerMock.mock.calls[0][1]).toEqual({ logLevel: DiagLogLevel.VERBOSE });
  });
});
