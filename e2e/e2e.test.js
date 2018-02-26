const { URL } = require('url');
const STAGING_BASE = new URL('https://app.staging.dp3.us/');

function buildStagingURL(path) {
  return new URL(path, STAGING_BASE);
}

var webdriver = require('selenium-webdriver'),
  By = webdriver.By,
  until = webdriver.until,
  promise = webdriver.promise,
  username = 'movemil',
  accessKey = process.env.SAUCE_ACCESS_KEY,
  driver;

// async/await do not work well when the promise manager is enabled.
promise.USE_PROMISE_MANAGER = false;

// jest.set_timeout() doesn't work in this circumstance, so using jasmine timeout.
jasmine.DEFAULT_TIMEOUT_INTERVAL = 1000 * 60 * 5;

beforeAll(async function() {
  driver = new webdriver.Builder()
    .withCapabilities({
      browserName: 'internet explorer',
      platform: 'Windows 8.1',
      version: '11.0',
      username: username,
      accessKey: accessKey,
    })
    .usingServer(
      'http://' +
        username +
        ':' +
        accessKey +
        '@ondemand.saucelabs.com:80/wd/hub',
    )
    .build();
});

describe('issue pages', async () => {
  beforeEach(async () => await driver.navigate().to(STAGING_BASE));

  it('loads Submit Feedback page', async () => {
    // When: Page is loaded, should display expected title
    await driver.wait(until.titleIs('Transcom PPP: Submit Feedback'), 2000);
  });

  it('allows issue submission and retrieval', async () => {
    // Given: A test issue and a feedback form on index page
    test_issue = 'Too much Alexi. Time: ' + Date.now();
    await driver.wait(
      until.elementLocated(By.css('[data-test="feedback-form"]')),
    );
    feedback_form = await driver.findElement(
      By.css('[data-test="feedback-form"]'),
    );
    feedback_form.clear();
    // When: Submit issue
    feedback_form.sendKeys(test_issue);
    await driver.findElement(By.css("input[type='submit']")).click();
    // Then: Visit submitted page
    await driver.get(buildStagingURL('submitted'));
    issue_cards = await driver.findElement(By.className('issue-cards'));
    // Expect: Submitted issue exists on page
    await driver.wait(until.elementTextContains(issue_cards, test_issue), 1000);
  });
});

describe('shipments pages', async () => {
  it('loads all shipments page', async () => {
    // When: Page is loaded, should display expected title
    await driver.navigate().to(buildStagingURL('shipments/all'));
    await driver.wait(until.titleIs('Transcom PPP: All Shipments'), 2000);
  });

  it('loads available shipments page', async () => {
    // When: Page is loaded, should display expected title
    await driver.navigate().to(buildStagingURL('shipments/available'));
    await driver.wait(until.titleIs('Transcom PPP: Available Shipments'), 2000);
  });

  it('loads awarded shipments page', async () => {
    // When: Page is loaded, should display expected title
    await driver.navigate().to(buildStagingURL('shipments/awarded'));
    await driver.wait(until.titleIs('Transcom PPP: Awarded Shipments'), 2000);
  });

  it('displays alert on incorrect url', async () => {
    await driver.navigate().to(buildStagingURL('shipments/dogs'));
    // Expect: Alert error exists on page
    await driver.wait(
      until.elementLocated(By.className('usa-alert-error')),
      2000,
    );
  });
});

describe('DD1299 page', async () => {
  beforeEach(async () => await driver.navigate().to(buildStagingURL('DD1299')));

  it('loads DD1299 page', async () => {
    // When: Page is loaded, should display expected title
    await driver.wait(until.titleIs('Transcom PPP: DD1299'), 2000);
  });
});

afterAll(async () => {
  await driver.quit();
});
