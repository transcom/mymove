/* eslint-disable no-console */
const puppeteer = require('puppeteer');

const { networkProfiles } = require('./constants');

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
      if (firstRequestStart === null || entry.requestStart < firstRequestStart) {
        firstRequestStart = entry.requestStart;
      }

      if (lastResponseEnd === null || entry.responseEnd > lastResponseEnd) {
        lastResponseEnd = entry.responseEnd;
      }
    }
  });

  return (lastResponseEnd - firstRequestStart) / 1000;
};

const setupEmulation = (config, page, userAgent) => {
  const { device } = config;

  const emulator = device ? puppeteer.devices[`${device}`] : config.emulate;

  if (device && !emulator) {
    console.log(`Skipping page emulation device '${device}' is not defined`);
  }

  if (emulator) {
    if (device) {
      console.log(`Emulating page using ${emulator.name}`);
    } else {
      emulator.userAgent = emulator.userAgent || userAgent;
    }
    return page.emulate(emulator).catch((err) => {
      console.error(`Error setting emulation of page`, err);
    });
  }

  return Promise.resolve();
};

const setupNetwork = (config, page) => {
  const networkType = networkProfiles[config.network] || config.throttling;

  if (networkType) {
    return page.emulateNetworkConditions(networkType);
  }

  return Promise.resolve();
};

const totalDuration = async (host, config, debug) => {
  const waitOptions = { timeout: 0, waitUntil: 'networkidle0' };

  const browser = await puppeteer.launch(config.launch);
  const userAgent = await browser.userAgent();

  debug(`browser version ${await browser.version()}`);

  const page = await browser.newPage();

  await setupEmulation(config, page, userAgent);

  await page.goto(`${host}/devlocal-auth/login`, waitOptions);

  // Login by clicking button for existing user
  const loginBtnSelector = 'button[value="9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"]';
  await page.waitForSelector(loginBtnSelector);
  await Promise.all([page.click(loginBtnSelector), page.waitForNavigation(waitOptions)]);

  // grab first table data for locator
  const locatorSelector = 'td[data-testid="locator-0"]';
  await page.waitForSelector(locatorSelector);
  const element = await page.$(locatorSelector);
  const locatorValue = await page.evaluate((el) => el.textContent, element);

  await setupNetwork(config, page);

  // go to a document viewer, orders
  await page.goto(`${host}/moves/${locatorValue}/orders`, waitOptions);

  const docViewerContentSelector = 'div[data-testid="DocViewerContent"]';
  await page.waitForSelector(docViewerContentSelector);

  // Will return all http requests and navigation performance on last navigation
  const navigationEntries = JSON.parse(
    await page.evaluate(() => {
      return JSON.stringify(performance.getEntries());
    }),
  );

  await browser.close();

  return getTotalRequestTime(navigationEntries);
};

module.exports = { totalDuration };
