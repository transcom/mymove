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
      // eslint-disable-next-line no-param-reassign
      registration.onupdatefound = () => {
        const installingWorker = registration.installing;
        installingWorker.onstatechange = () => {
          if (installingWorker.state === 'installed') {
            if (navigator.serviceWorker.controller) {
              // RA Summary: eslint: no-console - System Information Leak: External
              // RA: The linter flags any console.
              // RA: This console in this file serves to indicate the status of the serviceWorker to a user.
              // RA: Given that the value displayed is a simple string with no interpolation
              // RA: nor variable names, SQL strings, system path information, or source or program code,
              // RA: this is not a finding.
              // RA Developer Status: Mitigated
              // RA Validator Status: Mitigated
              // RA Validator: jneuner@mitre.org
              // RA Modified Severity: CAT III
              // At this point, the old content will have been purged and
              // the fresh content will have been added to the cache.
              // It's the perfect time to display a "New content is
              // available; please refresh." message in your web app.
              console.log('New content is available; please refresh.'); // eslint-disable-line no-console
            } else {
              // RA Summary: eslint: no-console - System Information Leak: External
              // RA: The linter flags any console.
              // RA: This console in this file serves to indicate the status of the serviceWorker to a user.
              // RA: Given that the value displayed is a simple string with no interpolation
              // RA: nor variable names, SQL strings, system path information, or source or program code,
              // RA: this is not a finding.
              // RA Developer Status: Mitigated
              // RA Validator Status: Mitigated
              // RA Validator: jneuner@mitre.org
              // RA Modified Severity: CAT III
              // At this point, everything has been precached.
              // It's the perfect time to display a
              // "Content is cached for offline use." message.
              console.log('Content is cached for offline use.'); // eslint-disable-line no-console
            }
          }
        };
      };
    })
    .catch((error) => {
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message when registering a valid service worker fails.
      // RA: TODO: The possible values of this error need to be investigated further to determine mitigation actions.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5595
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      console.error('Error during service worker registration:', error); // eslint-disable-line no-console
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
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any console.
      // RA: This console in this file serves to indicate the status of internet connection to a user.
      // RA: Given that the value displayed is a simple string with no interpolation
      // RA: nor variable names, SQL strings, system path information, or source or program code,
      // RA: this is not a finding.
      // RA Developer Status: Mitigated
      // RA Validator Status: Mitigated
      // RA Validator: jneuner@mitre.org
      // RA Modified Severity: CAT III
      console.log('No internet connection found. App is running in offline mode.'); // eslint-disable-line no-console
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
