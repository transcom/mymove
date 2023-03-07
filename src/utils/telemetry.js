import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { W3CTraceContextPropagator } from '@opentelemetry/core';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { registerInstrumentations } from '@opentelemetry/instrumentation';

import { gitBranch, gitSha, isTelemetryEnabled } from 'shared/constants';

const serviceVersion = `${gitSha}@${gitBranch}`;

const traceCollectorOptions = {
  url: '/client/collector',
  headers: {
    'Content-Type': 'application/json',
  },
  concurrencyLimit: 10,
};

// Exporter (opentelemetry collector hidden behind proxy)
const exporter = new OTLPTraceExporter(traceCollectorOptions);
/**
 * creates an open telemetry trace provider
 *
 * @param {string} serviceName
 */
export function configureTelemetry(serviceName) {
  // Trace provider (Main application trace)

  if (isTelemetryEnabled) {
    const provider = new WebTracerProvider({
      resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: serviceName,
        [SemanticResourceAttributes.SERVICE_VERSION]: serviceVersion,
      }),
    });

    // from https://www.npmjs.com/package/@opentelemetry/exporter-trace-otlp-http
    provider.addSpanProcessor(new BatchSpanProcessor(exporter), {
      // The maximum queue size. After the size is reached spans are dropped.
      maxQueueSize: 100,
      // The maximum batch size of every export. It must be smaller or equal to maxQueueSize.
      maxExportBatchSize: 10,
      // The interval between two consecutive exports
      scheduledDelayMillis: 500,
      // How long the export can run before it is cancelled
      exportTimeoutMillis: 30000,
    });

    provider.register({
      propagator: new W3CTraceContextPropagator(),
      contextManager: new ZoneContextManager(),
    });

    registerInstrumentations({
      // the SwaggerClient uses fetch under the hood
      instrumentations: [new DocumentLoadInstrumentation(), new FetchInstrumentation()],
    });
  }
}

export default configureTelemetry;
