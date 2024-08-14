// @ts-check
const { devices, defineConfig } = require('@playwright/test');

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// require('dotenv').config();

// mac arm machines are much faster than x86 and can handle more
// workers without increasing flakiness
const defaultWorkers = process.arch.startsWith('arm') ? 4 : 2;

/**
 * @see https://playwright.dev/docs/test-configuration
 * @type {import('@playwright/test').PlaywrightTestConfig}
 */
const config = {
  globalSetup: './playwright/tests/globalSetup.js',
  testDir: './playwright/tests',
  /* Maximum time one test can run for. */
  timeout: 30 * 1000,
  expect: {
    /**
     * Maximum time expect() should wait for the condition to be met.
     * For example in `await expect(locator).toHaveText();`
     */
    timeout: 5000,
  },
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 4 : 0,
  /* Opt out of parallel tests on CI, but default workers based on arch */
  workers: process.env.CI ? 1 : defaultWorkers,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [['html', { outputFolder: 'playwright/html-report' }]],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Maximum time each action such as `click()` can take. Defaults to 0 (no limit). */
    actionTimeout: 0,
    /* Base URL to use in actions like `await page.goto('/')`. */
    // baseURL: 'http://localhost:3000',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
  },

  /* Configure projects */
  projects: [
    {
      // ahobson 2022-12-08: for now, only test desktop chrome for admin
      name: 'admin',
      testMatch: 'admin/**/*.spec.js',
      use: {
        baseURL: process.env.PLAYWRIGHT_ADMIN_URL || 'http://adminlocal:3000',
        ...devices['Desktop Chrome'],
      },
    },

    {
      // ahobson 2022-12-08: for now, only test desktop chrome for my
      name: 'my',
      testMatch: 'my/**/*.spec.js',
      use: {
        baseURL: process.env.PLAYWRIGHT_MY_URL || 'http://milmovelocal:3000',
        ...devices['Desktop Chrome'],
      },
    },
    {
      // ahobson 2022-12-14: for now, only test desktop chrome for office
      name: 'office',
      testMatch: 'office/**/*.spec.js',
      use: {
        baseURL: process.env.PLAYWRIGHT_OFFICE_URL || 'http://officelocal:3000',
        ...devices['Desktop Chrome'],
      },
    },
    // ahobson 2022-12-08: leave examples for later
    // {
    //   name: 'firefox',
    //   use: {
    //     ...devices['Desktop Firefox'],
    //   },
    // },

    // {
    //   name: 'webkit',
    //   use: {
    //     ...devices['Desktop Safari'],
    //   },
    // },

    /* Test against mobile viewports. */
    // {
    //   name: 'Mobile Chrome',
    //   use: {
    //     ...devices['Pixel 5'],
    //   },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: {
    //     ...devices['iPhone 12'],
    //   },
    // },

    /* Test against branded browsers. */
    // {
    //   name: 'Microsoft Edge',
    //   use: {
    //     channel: 'msedge',
    //   },
    // },
    // {
    //   name: 'Google Chrome',
    //   use: {
    //     channel: 'chrome',
    //   },
    // },
  ],

  /* Folder for test artifacts such as screenshots, videos, traces, etc. */
  outputDir: 'playwright/results/',

  /* Run your local dev server before starting the tests */
  // webServer: {
  //   command: 'npm run start',
  //   port: 3000,
  // },
};

module.exports = defineConfig(config);
