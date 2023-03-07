// @ts-check
import React from 'react';
import hoistNonReactStatics from 'hoist-non-react-statics';
import { trace, ROOT_CONTEXT, SpanKind } from '@opentelemetry/api';
import { SemanticAttributes } from '@opentelemetry/semantic-conventions';

/** @type { import('@opentelemetry/api').Tracer | undefined } */
let tracer;

// default to starting span on location change
let startSpanOnLocationChange = true;

let useEffect;
let createRoutesFromElements;
let matchRoutes;
let useLocation;
let useNavigationType;

/** @type {import('@opentelemetry/api').Span | undefined} */
let routingSpan;

/** @type {number | undefined} */
let routingSpanStartedAt;

/** @type { import('@opentelemetry/api').Context } */
let activeRoutingContext;

/**
 * ContextManager for OtelRoutes
 */
export class OtelRouteContextManager {
  /**
   * @private
   * @type {import('@opentelemetry/api').ContextManager}
   */
  base;

  /**
   * @private
   * @type {number}
   */
  activeRouteTimeout = 2000;

  /**
   * @param {import('@opentelemetry/api').ContextManager} baseContextManager
   * @param {Object} [options]
   * @param {number} [options.activeRouteTimeout=2000]
   */
  constructor(baseContextManager, options) {
    this.base = baseContextManager;
    const { activeRouteTimeout = 2000 } = options || {};
    this.activeRouteTimeout = activeRouteTimeout;
  }

  /**
   * Get the current active context
   *
   * @returns { import('@opentelemetry/api').Context }
   */
  active() {
    // Only use the custom routing context to set the current span if
    // this is the ROOT_CONTEXT
    if (this.base.active() !== ROOT_CONTEXT) {
      return this.base.active();
    }

    // if no routingSpan has been set up, return the existing context
    if (!routingSpan || !routingSpanStartedAt) {
      return this.base.active();
    }

    // if the span is older than the route timeout, use the existing
    // context, resetting the routingSpan information
    const sinceStarted = Date.now() - routingSpanStartedAt;
    if (sinceStarted > this.activeRouteTimeout) {
      routingSpan = undefined;
      routingSpanStartedAt = undefined;
      activeRoutingContext = undefined;
      return this.base.active();
    }

    // if the active routing context has the current routing span,
    // return that context
    if (activeRoutingContext) {
      const span = trace.getSpan(activeRoutingContext);
      if (span?.spanContext().spanId === routingSpan.spanContext().spanId) {
        return activeRoutingContext;
      }
    }

    // set the routing span as the active span
    activeRoutingContext = trace.setSpan(this.base.active(), routingSpan);
    return activeRoutingContext;
  }

  /**
   * Run the fn callback with object set as the current active context
   *
   * @template {unknown[]} A
   * @template {(...args: A) => ReturnType<F>} F
   *
   * @param { import('@opentelemetry/api').Context } context
   * @param {F} fn
   * @param {ThisParameterType<F>} [thisArg]
   * @param {A} args
   * @returns unknown
   */
  with(context, fn, thisArg, ...args) {
    return this.base.with(context, fn, thisArg, ...args);
  }

  /**
   * Bind an object as the current context (or a specific one)
   * @param {import('@opentelemetry/api').Context} context
   * @param {any} target
   */
  bind(context, target) {
    return this.base.bind(context, target);
  }

  /**
   * Enable context management
   *
   * @returns { import('@opentelemetry/api').ContextManager }
   */
  enable() {
    return this;
  }

  /**
   * Diable context management
   *
   * @returns { import('@opentelemetry/api').ContextManager }
   */
  disable() {
    return this;
  }
}
// inspired by @sentry/react

/**
 * configure the otel routes, passing in react functions
 *
 * @param {string} serviceName
 * @param {string} serviceVersion
 * @param {boolean} startSpanOnLocationChangeOption
 * @param {function}  useEffectOption,
 * @param {function}  createRoutesFromElementsOption,
 * @param {function}  matchRoutesOption,
 * @param {function}  useLocationOption,
 * @param {function}  useNavigationTypeOption,
 */
export function configureOtelRoutes(
  serviceName,
  serviceVersion,
  startSpanOnLocationChangeOption,
  useEffectOption,
  createRoutesFromElementsOption,
  matchRoutesOption,
  useLocationOption,
  useNavigationTypeOption,
) {
  startSpanOnLocationChange = startSpanOnLocationChangeOption;
  useEffect = useEffectOption;
  createRoutesFromElements = createRoutesFromElementsOption;
  matchRoutes = matchRoutesOption;
  useLocation = useLocationOption;
  useNavigationType = useNavigationTypeOption;

  tracer = trace.getTracer(serviceName, serviceVersion);
}

/**
 *
 * @param {string} url
 * @returns {number}
 */
function getNumberOfUrlSegments(url) {
  // split at '/' or at '\/' to split regex urls correctly
  return url.split(/\\?\//).filter((s) => s.length > 0 && s !== ',').length;
}

/**
 * @typedef {Object} Location
 * @property {string} pathname
 */

/**
 * @typedef {Object} RouteObject
 * @property {any} index
 * @property {string?} path
 */

/**
 * @typedef {Object} RouteMatch
 * @property {string} pathname
 * @property {RouteObject} route
 */

/**
 *
 * @param {RouteObject[]} routes
 * @param {Location} location
 * @param {RouteMatch[]} branches
 * @returns {[string,string]}
 */
function getNormalizedName(routes, location, branches) {
  if (!routes || routes.length === 0) {
    return [location.pathname, 'url'];
  }

  let pathBuilder = '';
  if (branches) {
    for (let x = 0; x < branches.length; x += 1) {
      const branch = branches[x];
      const { route } = branch;
      if (route) {
        // Early return if index route
        if (route.index) {
          return [branch.pathname, 'route'];
        }

        const { path } = route;
        if (path) {
          const newPath = path[0] === '/' || pathBuilder[pathBuilder.length - 1] === '/' ? path : `/${path}`;
          pathBuilder += newPath;
          if (branch.pathname === location.pathname) {
            if (
              // If the route defined on the element is something like
              // <Route path="/stores/:storeId/products/:productId" element={<div>Product</div>} />
              // We should check against the branch.pathname for the number of / seperators
              getNumberOfUrlSegments(pathBuilder) !== getNumberOfUrlSegments(branch.pathname) &&
              // We should not count wildcard operators in the url segments calculation
              pathBuilder.slice(-2) !== '/*'
            ) {
              return [newPath, 'route'];
            }
            return [pathBuilder, 'route'];
          }
        }
      }
    }
  }

  return [location.pathname, 'url'];
}

/**
 * @param {Location} location
 * @param {RouteObject[]} routes
 * @param {unknown} [matches]
 * @param {string} [basename]
 */
function updatePageloadTransaction(location, routes, matches, basename) {
  const branches = Array.isArray(matches) ? matches : matchRoutes(routes, location, basename);

  if (routingSpan && branches) {
    routingSpan.setAttribute(SemanticAttributes.HTTP_ROUTE, getNormalizedName(routes, location, branches)[0]);
  }
}

/**
 * @param {Location} location
 * @param {RouteObject[]} routes
 * @param {unknown} navigationType
 * @param {unknown} [matches]
 * @param {string} [basename]
 */
function handleNavigation(location, routes, navigationType, matches, basename) {
  const branches = Array.isArray(matches) ? matches : matchRoutes(routes, location, basename);

  if (startSpanOnLocationChange && (navigationType === 'PUSH' || navigationType === 'POP') && branches) {
    if (routingSpan && routingSpan.isRecording()) {
      routingSpan.end();
    }
    const [name, source] = getNormalizedName(routes, location, branches);
    routingSpan = tracer.startSpan('React Route Navigation', {
      kind: SpanKind.CLIENT,
      attributes: {
        [SemanticAttributes.HTTP_TARGET]: name,
        'client.react_router.source': source,
      },
    });
    routingSpanStartedAt = Date.now();
  }
}

export function withOtelReactRouterV6Routing(Routes) {
  if (!useEffect || !useLocation || !useNavigationType || !createRoutesFromElements || !matchRoutes) {
    return Routes;
  }

  let isMountRenderPass = true;

  const OtelRoutes = (props) => {
    const { children } = props;
    const location = useLocation();
    const navigationType = useNavigationType();

    useEffect(
      () => {
        const routes = createRoutesFromElements(children);

        if (isMountRenderPass) {
          updatePageloadTransaction(location, routes);
          isMountRenderPass = false;
        } else {
          handleNavigation(location, routes, navigationType);
          if (routingSpan) {
            routingSpan.end();
          }
        }
      },
      // RA Summary: eslint - react-hooks/exhaustive-deps - possible inconsistent update
      // RA: `props.children` is purpusely not included in the dependency
      // RA: array, because we do not want to re-run this effect
      // RA: when the children change. We only want to start transactions
      // RA: when the location or navigation type change.
      // RA Developer Status: Mitigated
      // RA Validator Status: False Positive
      // RA Modified Severity: N/A
      // eslint-disable-next-line react-hooks/exhaustive-deps
      [location, navigationType],
    );

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Routes {...props} />;
  };

  hoistNonReactStatics(OtelRoutes, Routes);

  return OtelRoutes;
}

export default withOtelReactRouterV6Routing;
