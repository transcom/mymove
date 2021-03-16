/* eslint-disable no-console, security/detect-non-literal-fs-filename */
const fs = require('fs');

const lighthouse = require('lighthouse');
const puppeteer = require('puppeteer');
const reportGenerator = require('lighthouse/lighthouse-core/report/report-generator');
const { throttling } = require('lighthouse/lighthouse-core/config/constants');

const {
  networkProfiles,
  fileSizeMoveCodes,
  fileSizePaymentRequestIds,
  scenarios,
  measurementTypes,
} = require('./constants');

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
  return Promise.resolve(1);
};

const lighthouseFromPuppeteer = async (
  url,
  options,
  config = null,
  saveReports = false,
  networkProfileName = '',
  fileSizeName = '',
  measurement,
) => {
  // Run Lighthouse
  const { lhr, artifacts } = await lighthouse(url, options, config).catch((err) => {
    console.error('lighthouse audit error \n', err);
    process.exit(1);
  });

  if (!fs.existsSync('puppeteer/reports')) {
    fs.mkdirSync('puppeteer/reports');
  }

  // For debugging
  const json = reportGenerator.generateReport(lhr, 'json');
  if (saveReports) {
    const trace = reportGenerator.generateReport(artifacts, 'json');
    let lhReportPath;
    let pTracePath;

    const dateTime = new Date();
    const dateString = [
      dateTime.getFullYear(),
      dateTime.getMonth() + 1,
      dateTime.getDate(),
      dateTime.getHours(),
      dateTime.getMinutes(),
      dateTime.getSeconds(),
    ].join('');

    if (measurement === measurementTypes.networkComparison || measurement === measurementTypes.totalDuration) {
      lhReportPath = networkProfileName
        ? `puppeteer/reports/lighthouse-report-${networkProfileName}-${dateString}.json`
        : `puppeteer/reports/lighthouse-report-${dateString}.json`;
      pTracePath = networkProfileName
        ? `puppeteer/reports/performance-trace-${networkProfileName}-${dateString}.json`
        : `puppeteer/reports/performance-trace-${dateString}.json`;

      fs.writeFileSync(lhReportPath, json);
      fs.writeFileSync(pTracePath, trace);
    } else if (measurement === measurementTypes.fileDuration) {
      lhReportPath = fileSizeName
        ? `puppeteer/reports/lighthouse-report-${fileSizeName}-${dateString}.json`
        : `puppeteer/reports/lighthouse-report-${dateString}.json`;
      pTracePath = fileSizeName
        ? `puppeteer/reports/performance-trace-${fileSizeName}-${dateString}.json`
        : `puppeteer/reports/performance-trace-${dateString}.json`;

      fs.writeFileSync(lhReportPath, json);
      fs.writeFileSync(pTracePath, trace);
    }

    // eslint-disable-next-line no-console
    console.debug(`Generating lighthouse reports:
        ${lhReportPath}
        ${pTracePath}`);
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

const totalDuration = async ({ scenario, measurement, host, config, debug, saveReports, verbose }) => {
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

  // go to a document viewer
  let url;
  if (scenario === scenarios.tio) {
    url = `${host}/moves/${fileSizeMoveCodes[config.fileSize]}/payment-requests/${
      fileSizePaymentRequestIds[config.fileSize]
    }`;
  } else {
    // default to TOO
    url = `${host}/moves/${fileSizeMoveCodes[config.fileSize]}/orders`;
  }

  debug(`URL to gather metrics from: ${url}`);
  await page.goto(url, waitOptions).catch((err) => {
    console.error(`failed to navigate to page ${url}\n`, err);
    process.exit(1);
  });

  const docViewerContentSelector = 'div[data-testid="DocViewerContent"]';
  await page.waitForSelector(docViewerContentSelector).catch((err) => {
    console.error('waiting for document viewer selector timed out\n', err);
    process.exit(1);
  });

  // Will return all http requests and navigation performance on last navigation
  debug('Gathering performance timing metrics');
  const performanceEntries = await page
    .evaluate(() => {
      return JSON.stringify(performance.getEntries());
    })
    .catch((err) => {
      console.error('failed to fetch performance entries', err);
    });
  const navigationEntries = JSON.parse(performanceEntries);

  const lhOptions = {
    port: new URL(browser.wsEndpoint()).port,
    logLevel: verbose ? 'info' : 'error',
    chromeFlags: ['--disable-mobile-emulation'],
  };

  const lhConfig = {
    extends: 'lighthouse:default',
    settings: {
      maxWaitForLoad: 300000, // 300,000 ms. Increase max wait so lighthouse can gather metrics.
      formFactor: 'desktop', // TODO - Will need to change once we do device emulation
      screenEmulation: {
        mobile: false,
        width: deviceEmulationConfig.viewport.width,
        height: deviceEmulationConfig.viewport.height,
        deviceScaleFactor: 1,
        disabled: false,
      },
      emulatedUserAgent: deviceEmulationConfig.userAgent,
      throttlingMethod: 'devtools',
      throttling: {
        cpuSlowdownMultiplier: cpuRateConfig.rate,
        requestLatencyMs: 0,
        downloadThroughputKbps: 0,
        uploadThroughputKbps: 0,
      },
    },
  };

  if (networkConfig) {
    lhConfig.settings.throttling = {
      ...lhConfig.settings.throttling,
      // Lighthouse expects kilobits per sec instead of bytes per sec
      // 1 byte = 0.008 kilobits
      requestLatencyMs: networkConfig.latency * throttling.DEVTOOLS_RTT_ADJUSTMENT_FACTOR,
      downloadThroughputKbps: networkConfig.download * 0.008 * throttling.DEVTOOLS_THROUGHPUT_ADJUSTMENT_FACTOR,
      uploadThroughputKbps: networkConfig.upload * 0.008 * throttling.DEVTOOLS_THROUGHPUT_ADJUSTMENT_FACTOR,
    };
  }

  debug('Gathering lighthouse metrics');
  const lhResults = await lighthouseFromPuppeteer(
    url,
    lhOptions,
    lhConfig,
    saveReports,
    config.network,
    config.fileSize,
    measurement,
  ).catch((err) => {
    console.error('error running lighthouse\n', err);
  });

  const pfTimingResults = getTotalRequestTime(navigationEntries);

  await browser.close();

  return { '🏁 Peformance timing (seconds)': pfTimingResults, ...lhResults };
};

module.exports = { totalDuration, lighthouseFromPuppeteer };
