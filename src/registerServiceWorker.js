import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

// In production, we register a service worker to serve assets from local cache.

// This lets the app load faster on subsequent visits in production, and gives
// it offline capabilities. However, it also means that developers (and users)
// will only see deployed updates on the "N+1" visit to a page, since previously
// cached resources are updated in the background.

// To learn more about the benefits of this model, read https://goo.gl/KwvDNy.
// This link also includes instructions on opting out of this behavior.

const isLocalhost = Boolean(
  window.location.hostname === 'localhost' ||
    // [::1] is the IPv6 localhost address.
    window.location.hostname === '[::1]' ||
    // milmovelocal is the default server name.
    window.location.hostname === 'milmovelocal' ||
    // RA Summary: eslint - security/detect-unsafe-regex - Denial of Service: Regular Expression
    // RA: Locates potentially unsafe regular expressions, which may take a very long time to run, blocking the event loop
    // RA: Per MilMove SSP, predisposing conditions are regex patterns from untrusted sources or unbounded matching.
    // RA: The regex pattern is a constant string set at compile-time and it is bounded to 15 characters (127.000.000.001).
    // RA Developer Status: Mitigated
    // RA Validator Status: Mitigated
    // RA Modified Severity: N/A
    // 127.0.0.1/8 is considered localhost for IPv4.
    // eslint-disable-next-line security/detect-unsafe-regex
    window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/),
);

function registerValidSW(swUrl) {
  navigator.serviceWorker
    .register(swUrl)
    .then((registration) => {
      // More info: https://spin.atomicobject.com/2011/04/10/javascript-don-t-reassign-your-function-arguments/
      // This is done to avoid an edge case when mutating the argument object. Although this is not an example of the edgecase.
      const serviceWorkerRegistration = registration;
      serviceWorkerRegistration.onupdatefound = () => {
        const installingWorker = registration.installing;
        installingWorker.onstatechange = () => {
          if (installingWorker.state === 'installed') {
            if (navigator.serviceWorker.controller) {
              milmoveLog(MILMOVE_LOG_LEVEL.LOG, 'New content is available; please refresh.');
            } else {
              milmoveLog(MILMOVE_LOG_LEVEL.LOG, 'Content is cached for offline use.');
            }
          }
        };
      };
    })
    .catch((error) => {
      milmoveLog(MILMOVE_LOG_LEVEL.ERROR, 'Error during service worker registration:', error);
    });
}

function checkValidServiceWorker(swUrl) {
  // Check if the service worker can be found. If it can't reload the page.
  fetch(swUrl)
    .then((response) => {
      // Ensure service worker exists, and that we really are getting a JS file.
      if (response.status === 404 || response.headers.get('content-type').indexOf('javascript') === -1) {
        // No service worker found. Probably a different app. Reload the page.
        navigator.serviceWorker.ready.then((registration) => {
          registration.unregister().then(() => {
            window.location.reload();
          });
        });
      } else {
        // Service worker found. Proceed as normal.
        registerValidSW(swUrl);
      }
    })
    .catch(() => {
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, 'No internet connection found. App is running in offline mode.');
    });
}

export function unregister() {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.ready.then((registration) => {
      registration.unregister();
    });
  }
}

export default function register() {
  if (process.env.NODE_ENV === 'production' && 'serviceWorker' in navigator) {
    // The URL constructor is available in all browsers that support SW.
    const publicUrl = new URL(process.env.PUBLIC_URL, window.location);
    if (publicUrl.origin !== window.location.origin) {
      // Our service worker won't work if PUBLIC_URL is on a different origin
      // from what our page is served on. This might happen if a CDN is used to
      // serve assets; see https://github.com/facebookincubator/create-react-app/issues/2374
      return;
    }

    window.addEventListener('load', () => {
      const swUrl = `${process.env.PUBLIC_URL}/service-worker.js`;

      if (isLocalhost) {
        // This is running on milmovelocal. Lets check if a service worker still exists or not.
        checkValidServiceWorker(swUrl);
      } else {
        // Is not local host. Just register service worker
        registerValidSW(swUrl);
      }
    });
  }
}
