var webdriver = require('selenium-webdriver'),
  By = webdriver.By,
  until = webdriver.until,
  username = 'movemil',
  accessKey = process.env.SAUCE_ACCESS_KEY,
  driver;

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

// jest.set_timeout() doesn't work in this circumstance, so using jasmine timeout.
jasmine.DEFAULT_TIMEOUT_INTERVAL = 1000 * 60 * 5;

describe('issue pages', () => {
  beforeEach(() => driver.navigate().to('https://app.staging.dp3.us/'));

  it('loads Submit Feedback page', () => {
    // When: Page is loaded, should display expected title
    driver.wait(until.titleIs('Transcom PPP: Submit Feedback'), 2000);
  });

  it('allows issue submission and retrieval', () => {
    // Given: A test issue and a feedback form on index page
    test_issue = 'Too few dogs. Time: ' + Date.now();
    driver.wait(until.elementLocated(By.css('[data-test="feedback-form"]')));
    feedback_form = driver.findElement(By.css('[data-test="feedback-form"]'));
    feedback_form.clear();
    // When: Submit issue
    feedback_form.sendKeys(test_issue);
    driver.findElement(By.css("input[type='submit']")).click();
    // Then: Visit submitted page
    driver.get('https://app.staging.dp3.us/submitted');
    issue_cards = driver.findElement(By.className('issue-cards'));
    // Expect: Submitted issue exists on page
    driver.wait(until.elementTextContains(issue_cards, test_issue), 1000);
  });
});

describe('shipments pages', () => {
  it('loads all shipments page', () => {
    // When: Page is loaded, should display expected title
    driver.navigate().to('https://app.staging.dp3.us/shipments/all');
    driver.wait(until.titleIs('Transcom PPP: All Shipments'), 2000);
  });

  it('loads available shipments page', () => {
    // When: Page is loaded, should display expected title
    driver.navigate().to('https://app.staging.dp3.us/shipments/available');
    driver.wait(until.titleIs('Transcom PPP: Available Shipments'), 2000);
  });

  it('loads awarded shipments page', () => {
    // When: Page is loaded, should display expected title
    driver.navigate().to('https://app.staging.dp3.us/shipments/awarded');
    driver.wait(until.titleIs('Transcom PPP: Awarded Shipments'), 2000);
  });

  it('displays alert on incorrect url', () => {
    driver.navigate().to('https://app.staging.dp3.us/shipments/dogs');
    // Expect: Alert error exists on page
    driver.wait(until.elementLocated(By.className('usa-alert-error')), 2000);
  });
});

describe('DD1299 page', () => {
  beforeEach(() => driver.navigate().to('https://app.staging.dp3.us/DD1299'));

  it('loads Submit Feedback page', () => {
    // When: Page is loaded, should display expected title
    driver.wait(until.titleIs('Transcom PPP: 1299'), 2000);
  });
});

afterAll(() => {
  driver.quit();
});
