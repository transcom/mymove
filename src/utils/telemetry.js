// @ts-check
import { SpanStatusCode } from '@opentelemetry/api';
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
import { useEffect } from 'react';
import { createRoutesFromElements, matchRoutes, useLocation, useNavigationType } from 'react-router-dom';

import { configureOtelRoutes, OtelRouteContextManager } from 'components/ThirdParty/OtelRoutes';
import { gitBranch, gitSha, isTelemetryEnabled, serviceName } from 'shared/constants';

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
 * Creates an open telemetry trace provider
 *
 * NOTE: changes to this function require restarting the client
 * server, the automatic reloading will not suffice
 */
export function configureTelemetry() {
  // Trace provider (Main application trace)

  if (isTelemetryEnabled) {
    const provider = new WebTracerProvider({
      resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: serviceName(),
        [SemanticResourceAttributes.SERVICE_VERSION]: serviceVersion,
      }),
    });

    // from https://www.npmjs.com/package/@opentelemetry/exporter-trace-otlp-http
    provider.addSpanProcessor(
      new BatchSpanProcessor(exporter, {
        // The maximum queue size. After the size is reached spans are dropped.
        maxQueueSize: 200,
        // The maximum batch size of every export. It must be smaller or equal to maxQueueSize.
        maxExportBatchSize: 20,
        // The interval between two consecutive exports
        scheduledDelayMillis: 500,
        // How long the export can run before it is cancelled
        exportTimeoutMillis: 30000,
      }),
    );

    const zoneContextManager = new ZoneContextManager();
    const routeContextManager = new OtelRouteContextManager(zoneContextManager.enable());
    provider.register({
      propagator: new W3CTraceContextPropagator(),
      contextManager: routeContextManager,
    });

    /**
     * Set custom fetch attributes
     * @param {import('@opentelemetry/api').Span} span
     * @param {Request} _
     * @param {Response} response
     */
    const applyCustomAttributesOnSpan = (span, _, response) => {
      if (response.status >= 400) {
        span.setStatus({
          code: SpanStatusCode.ERROR,
          message: response.statusText,
        });
      }
    };

    registerInstrumentations({
      // the SwaggerClient uses fetch under the hood
      instrumentations: [new DocumentLoadInstrumentation(), new FetchInstrumentation({ applyCustomAttributesOnSpan })],
    });

    configureOtelRoutes(
      serviceName(),
      serviceVersion,
      true,
      useEffect,
      createRoutesFromElements,
      matchRoutes,
      useLocation,
      useNavigationType,
    );
  }
}

export default configureTelemetry;
