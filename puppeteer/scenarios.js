/* eslint-disable no-console */
const fs = require('fs');

const lighthouse = require('lighthouse');
const puppeteer = require('puppeteer');
const reportGenerator = require('lighthouse/lighthouse-core/report/report-generator');
const { throttling } = require('lighthouse/lighthouse-core/config/constants');

const { networkProfiles, fileSizeMoveCodes } = require('./constants');

// Gets the total request time based on the responseEnd - requestStart
// in secs
const getTotalRequestTime = (navigationEntries = []) => {
  if (navigationEntries.length === 0) {
    return null;
  }

  let firstRequestStart = null;
  let lastResponseEnd = null;

  navigationEntries.forEach((entry) => {
    if (entry.entryType === 'resource' || entry.entryType === 'navigation') {
      // requestStart is null on some entries
      if (firstRequestStart === null || (entry.requestStart || entry.fetchStart) < firstRequestStart) {
        firstRequestStart = entry.requestStart;
      }

      if (lastResponseEnd === null || entry.responseEnd > lastResponseEnd) {
        lastResponseEnd = entry.responseEnd;
      }
    }
  });

  return `${((lastResponseEnd - firstRequestStart) / 1000).toFixed(1)} s`;
};

// Gets the request time for the orders image
const getOrdersImageRequestTime = (navigationEntries = []) => {
  if (navigationEntries.length === 0) {
    return null;
  }

  const entry = navigationEntries.find((o) => o.name.includes('png'));

  return {
    totalTime: `${(entry.duration / 1000).toFixed(1)} s`,
    downloading: `${((entry.responseEnd - entry.fetchStart) / 1000).toFixed(1)} s`,
    waiting: `${((entry.responseStart - entry.fetchStart) / 1000).toFixed(1)} s`,
  };
};

const setupEmulation = (config, page, userAgent) => {
  const { device } = config;

  const emulator = device ? puppeteer.devices[`${device}`] : config.emulate;

  if (device && !emulator) {
    console.debug(`Skipping page emulation device '${device}' is not defined`);
  }

  if (emulator) {
    if (device) {
      console.debug(`Emulating page using ${emulator.name}`);
    } else {
      emulator.userAgent = emulator.userAgent || userAgent;
    }
    return page
      .emulate(emulator)
      .then(() => {
        return Promise.resolve(emulator);
      })
      .catch((err) => {
        console.error(`Error setting emulation of page`, err);
      });
  }

  return Promise.resolve();
};

const setupNetwork = (config, page) => {
  const networkType = networkProfiles[config.network] || config.throttling;

  if (networkType) {
    return page.emulateNetworkConditions(networkType).then(() => {
      return Promise.resolve(networkType);
    });
  }

  return Promise.resolve();
};

const setupCPU = async (config, page) => {
  const { cpuSlowdownRate } = config;

  if (cpuSlowdownRate > 1) {
    const client = await page.target().createCDPSession();
    const rate = { rate: cpuSlowdownRate };
    return client.send('Emulation.setCPUThrottlingRate', rate).then(() => {
      return Promise.resolve(rate);
    });
  }
  return Promise.resolve();
};

const lighthouseFromPuppeteer = async (url, options, config = null, saveReports = false, networkProfileName = '') => {
  // Run Lighthouse
  const { lhr, artifacts } = await lighthouse(url, options, config);

  // For debugging
  const json = reportGenerator.generateReport(lhr, 'json');
  if (saveReports) {
    const trace = reportGenerator.generateReport(artifacts, 'json');
    const lhReportPath = `puppeteer/lighthouse-report-${networkProfileName}.json`;
    const pTracePath = `puppeteer/performance-trace-${networkProfileName}.json`;
    // eslint-disable-next-line no-console
    console.debug(`\n
    Generating lighthouse reports:
      ${lhReportPath}
      ${pTracePath}`);
    // Need literal string to work around the detect-non-literal-fs-filename
    switch (networkProfileName) {
      case 'fast':
        fs.writeFileSync('puppeteer/lighthouse-report-fast.json', json);
        fs.writeFileSync('puppeteer/performance-trace-fast.json', trace);
        break;
      case 'medium':
        fs.writeFileSync('puppeteer/lighthouse-report-medium.json', json);
        fs.writeFileSync('puppeteer/performance-trace-medium.json', trace);
        break;
      case 'slow':
        fs.writeFileSync('puppeteer/lighthouse-report-slow.json', json);
        fs.writeFileSync('puppeteer/performance-trace-slow.json', trace);
        break;
      default:
        fs.writeFileSync('puppeteer/lighthouse-report.json', json);
        fs.writeFileSync('puppeteer/performance-trace.json', trace);
    }
  }

  const { audits } = JSON.parse(json); // Lighthouse audits
  const largestContentfulPaint = audits['largest-contentful-paint'].displayValue;
  const totalBlockingTime = audits['total-blocking-time'].displayValue;
  const timeToInteractive = audits.interactive.displayValue;

  return {
    '🎨 Largest Contentful Paint (seconds)': largestContentfulPaint,
    '👆 Time To Interactive (seconds)': timeToInteractive,
    '⌛️ Total Blocking Time (milliseconds)': totalBlockingTime,
  };
};

const totalDuration = async ({ host, config, debug, saveReports, verbose }) => {
  const waitOptions = { timeout: 0, waitUntil: 'networkidle0' };

  const browser = await puppeteer.launch(config.launch);
  const userAgent = await browser.userAgent();

  debug(`browser version ${await browser.version()}`);

  const page = await browser.newPage();

  const deviceEmulationConfig = await setupEmulation(config, page, userAgent);

  await page.goto(`${host}/devlocal-auth/login`, waitOptions).catch(() => {
    console.error(`Unable to reach host ${host}. Make sure your server and client are already running`);
    return Promise.reject();
  });

  // Login by clicking button for existing user
  const loginBtnSelector = 'button[value="9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"]';
  await page.waitForSelector(loginBtnSelector).catch(() => {
    console.error(`Unable to reach host ${host}. Make sure your server and client are already running.`);
    return Promise.reject();
  });

  await Promise.all([page.click(loginBtnSelector), page.waitForNavigation(waitOptions)]);

  const networkConfig = await setupNetwork(config, page);
  const cpuRateConfig = await setupCPU(config, page);

  // go to a document viewer, orders
  const url = `${host}/moves/${fileSizeMoveCodes[config.fileSize]}/orders`;
  debug(`URL to gather metrics from: ${url}`);
  await page.goto(url, waitOptions);

  const docViewerContentSelector = 'div[data-testid="DocViewerContent"]';
  await page.waitForSelector(docViewerContentSelector);

  // Will return all http requests and navigation performance on last navigation
  debug('Gathering performance timing metrics');
  const navigationEntries = JSON.parse(
    await page.evaluate(() => {
      return JSON.stringify(performance.getEntries());
    }),
  );

  const lhOptions = {
    port: new URL(browser.wsEndpoint()).port,
    logLevel: verbose ? 'info' : 'error',
    chromeFlags: ['--disable-mobile-emulation'],
  };
  const lhConfig = {
    extends: 'lighthouse:default',
    settings: {
      maxWaitForLoad: 200000, // 200,000 ms. Increase max wait so lighthouse can gather metrics.
      formFactor: 'desktop', // TODO - Will need to change once we do device emulation
      throttlingMethod: 'devtools',
      throttling: {
        // Lighthouse expects kilobits per sec instead of bytes per sec
        // 1 byte = 0.008 kilobits
        cpuSlowdownMultiplier: cpuRateConfig.rate,
        requestLatencyMs: networkConfig.latency * throttling.DEVTOOLS_RTT_ADJUSTMENT_FACTOR,
        downloadThroughputKbps: networkConfig.download * 0.008 * throttling.DEVTOOLS_THROUGHPUT_ADJUSTMENT_FACTOR,
        uploadThroughputKbps: networkConfig.upload * 0.008 * throttling.DEVTOOLS_THROUGHPUT_ADJUSTMENT_FACTOR,
      },
      screenEmulation: {
        mobile: false,
        width: deviceEmulationConfig.viewport.width,
        height: deviceEmulationConfig.viewport.height,
        deviceScaleFactor: 1,
        disabled: false,
      },
      emulatedUserAgent: deviceEmulationConfig.userAgent,
    },
  };

  debug('Gathering lighthouse metrics');
  const lhResults = await lighthouseFromPuppeteer(url, lhOptions, lhConfig, saveReports, config.network);
  const pfTimingResults = getTotalRequestTime(navigationEntries);

  await browser.close();

  return { '🏁 Peformance timing (seconds)': pfTimingResults, ...lhResults };
};

const fileDownloadDuration = async ({ host, config, debug, fileSize }) => {
  const waitOptions = { timeout: 0, waitUntil: 'networkidle0' };

  const browser = await puppeteer.launch(config.launch);

  debug(`browser version ${await browser.version()}`);

  const page = await browser.newPage();

  await page.goto(`${host}/devlocal-auth/login`, waitOptions).catch(() => {
    console.error(`Unable to reach host ${host}. Make sure your server and client are already running`);
    return Promise.reject();
  });

  // Login by clicking button for existing user
  const loginBtnSelector = 'button[value="9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"]';
  await page.waitForSelector(loginBtnSelector).catch(() => {
    console.error(`Unable to reach host ${host}. Make sure your server and client are already running.`);
    return Promise.reject();
  });

  await Promise.all([page.click(loginBtnSelector), page.waitForNavigation(waitOptions)]);

  // go to a document viewer, orders
  const fileSizeMoveCode = fileSizeMoveCodes[`${fileSize}`];
  const url = `${host}/moves/${fileSizeMoveCode}/orders`;
  debug(`URL to gather metrics from: ${url}`);
  await page.goto(url, waitOptions);

  const docViewerContentSelector = 'div[data-testid="DocViewerContent"]';
  await page.waitForSelector(docViewerContentSelector);

  // Will return all http requests and navigation performance on last navigation
  debug('Gathering performance timing metrics');
  const navigationEntries = JSON.parse(
    await page.evaluate(() => {
      return JSON.stringify(performance.getEntries());
    }),
  );
  const pfTimingResults = getOrdersImageRequestTime(navigationEntries);

  await browser.close();

  return {
    '⌛  Wait time (seconds': pfTimingResults.waiting,
    '📥  Downloading (seconds)': pfTimingResults.downloading,
    '🏁  Total time (seconds)': pfTimingResults.totalTime,
  };
};

module.exports = { totalDuration, lighthouseFromPuppeteer, fileDownloadDuration };
