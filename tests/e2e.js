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

function load_app_test(driver) {
  driver.get('https://app.staging.dp3.us/');
  driver.wait(until.titleIs('Transcom PPP: Submit Feedback'));
}

function submit_issue_test(driver) {
  test_issue = 'Too few dogs. Time: ' + Date.now();
  driver.get('https://app.staging.dp3.us/');
  driver.wait(until.elementLocated(By.css('[data-test="feedback-form"]')));
  feedback_form = driver.findElement(By.css('[data-test="feedback-form"]'));
  feedback_form.clear();
  feedback_form.sendKeys(test_issue);
  driver.findElement(By.css("input[type='submit']")).click();
  driver.get('https://app.staging.dp3.us/submitted');
  issue_cards = driver.findElement(By.className('issue-cards'));
  driver.wait(until.elementTextContains(issue_cards, test_issue), 1000);
}

load_app_test(driver);
submit_issue_test(driver);
driver.quit();
