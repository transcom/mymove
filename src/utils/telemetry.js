import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor, RandomIdGenerator } from '@opentelemetry/sdk-trace-base';
import { CompositePropagator, W3CTraceContextPropagator } from '@opentelemetry/core';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { AWSXRayPropagator } from '@opentelemetry/propagator-aws-xray';
import { AWSXRayIdGenerator } from '@opentelemetry/id-generator-aws-xray';

import { gitBranch, gitSha, isTelemetryEnabled, isXrayEnabled } from 'shared/constants';

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
    /** @type {import('@opentelemetry/sdk-trace-base').IdGenerator} */
    let idGenerator;
    /** @type {import('@opentelemetry/api').Attributes} */
    const attributes = {
      [SemanticResourceAttributes.SERVICE_NAME]: serviceName,
      [SemanticResourceAttributes.SERVICE_VERSION]: serviceVersion,
    };
    if (isXrayEnabled) {
      idGenerator = new AWSXRayIdGenerator();
      attributes[SemanticResourceAttributes.CLOUD_PROVIDER] = 'aws';
    } else {
      idGenerator = new RandomIdGenerator();
    }
    /** @type {import('@opentelemetry/sdk-trace-base').TracerConfig} */
    const tracerConfig = {
      resource: new Resource(attributes),
      idGenerator,
    };
    const provider = new WebTracerProvider(tracerConfig);

    // The following is inspired by
    // https://www.npmjs.com/package/@opentelemetry/exporter-trace-otlp-http
    //
    /** @type {import('@opentelemetry/sdk-trace-base').BatchSpanProcessorBrowserConfig} */
    const batchSpanProcessorConfig = {
      // The maximum queue size. After the size is reached spans are dropped.
      maxQueueSize: 100,
      // The maximum batch size of every export. It must be smaller or equal to maxQueueSize.
      maxExportBatchSize: 10,
      // The interval between two consecutive exports
      scheduledDelayMillis: 500,
      // How long the export can run before it is cancelled
      exportTimeoutMillis: 30000,
      // no need to send info when not active
      disableAutoFlushOnDocumentHide: true,
    };
    provider.addSpanProcessor(new BatchSpanProcessor(exporter, batchSpanProcessorConfig));

    /** @type {import('@opentelemetry/core').CompositePropagatorConfig} */
    let propagatorConfig;
    if (isXrayEnabled) {
      propagatorConfig = {
        propagators: [new AWSXRayPropagator(), new W3CTraceContextPropagator()],
      };
    } else {
      propagatorConfig = {
        propagators: [new W3CTraceContextPropagator()],
      };
    }

    const propagator = new CompositePropagator(propagatorConfig);

    provider.register({
      propagator,
      contextManager: new ZoneContextManager(),
    });

    registerInstrumentations({
      // the SwaggerClient uses fetch under the hood
      instrumentations: [new DocumentLoadInstrumentation(), new FetchInstrumentation()],
    });
  }
}

export default configureTelemetry;
