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

jasmine.DEFAULT_TIMEOUT_INTERVAL = 1000 * 60 * 5;

describe('index page loads', () => {
  beforeEach(() => driver.navigate().to('https://app.staging.dp3.us/'));

  it('loads Submit Feedback page', () => {
    driver.wait(until.titleIs('Transcom PPP: Submit Feedback'), 2000);

    return driver.getTitle().then(title => {
      expect(title).toBe('Transcom PPP: Submit Feedback');
    });
  });

  it('allows issue submission', () => {
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

    return issue_cards.getText().then(issues => {
      expect(issues).toEqual(expect.stringContaining(test_issue));
    });
  });
});

afterAll(() => {
  driver.quit();
});
