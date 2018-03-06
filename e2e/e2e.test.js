const { URL } = require('url');

const E2E_BASE = process.env.E2E_BASE || 'https://app.staging.dp3.us/';
if (process.env.E2E_BASE) console.log('base url is ', E2E_BASE);
const BASE_URL = new URL(E2E_BASE);

function buildURL(path) {
  return new URL(path, BASE_URL);
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

const WEB_DRIVER_WAIT_TIME = 1000 * 10;

beforeAll(async function() {
  driver = new webdriver.Builder()
    .withCapabilities({
      browserName: 'internet explorer',
      platform: 'Windows 10',
      version: '11.103',
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
  beforeEach(async () => await driver.navigate().to(BASE_URL));

  it('loads Submit Feedback page', async () => {
    // When: Page is loaded, should display expected title
    await driver.wait(
      until.titleIs('Transcom PPP: Submit Feedback'),
      WEB_DRIVER_WAIT_TIME,
    );
  });

  it('allows issue submission and retrieval', async () => {
    // Given: A test issue and a feedback form on index page
    const descriptionTextArea = 'textarea[name="description"]';
    test_issue = 'Too few dogs. Time: ' + Date.now();
    await driver.wait(until.elementLocated(By.css(descriptionTextArea)));
    feedback_form = await driver.findElement(By.css(descriptionTextArea));
    feedback_form.clear();
    // When: Submit issue
    feedback_form.sendKeys(test_issue);
    await feedback_form.submit();

    // Then: Visit submitted page
    await driver.get(buildURL('submitted'));
    await driver.wait(
      until.titleIs('Transcom PPP: Submitted Feedback'),
      WEB_DRIVER_WAIT_TIME,
    );

    issue_cards = await driver.findElement(By.className('issue-cards'));
    // Expect: Submitted issue exists on page
    await driver.wait(
      until.elementTextContains(issue_cards, test_issue),
      WEB_DRIVER_WAIT_TIME,
    );
  });
});

describe('shipments pages', async () => {
  it('loads all shipments page', async () => {
    // When: Page is loaded, should display expected title
    await driver.navigate().to(buildURL('shipments/all'));
    await driver.wait(
      until.titleIs('Transcom PPP: All Shipments'),
      WEB_DRIVER_WAIT_TIME,
    );
  });

  it('loads available shipments page', async () => {
    // When: Page is loaded, should display expected title
    await driver.navigate().to(buildURL('shipments/available'));
    await driver.wait(
      until.titleIs('Transcom PPP: Available Shipments'),
      WEB_DRIVER_WAIT_TIME,
    );
  });

  it('loads awarded shipments page', async () => {
    // When: Page is loaded, should display expected title
    await driver.navigate().to(buildURL('shipments/awarded'));
    await driver.wait(
      until.titleIs('Transcom PPP: Awarded Shipments'),
      WEB_DRIVER_WAIT_TIME,
    );
  });

  it('displays alert on incorrect url', async () => {
    await driver.navigate().to(buildURL('shipments/dogs'));
    // Expect: Alert error exists on page
    await driver.wait(
      until.elementLocated(By.className('usa-alert-error')),
      WEB_DRIVER_WAIT_TIME,
    );
  });
});

describe('DD1299 page', async () => {
  beforeEach(async () => await driver.navigate().to(buildURL('DD1299')));

  it('loads DD1299 page', async () => {
    // When: Page is loaded, should display expected title
    await driver.wait(
      until.titleIs('Transcom PPP: DD1299'),
      WEB_DRIVER_WAIT_TIME,
    );
  });
});

afterAll(async () => {
  await driver.quit();
});
