var webdriver = require('selenium-webdriver'),
  By = webdriver.By,
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
  driver.sleep(2000).then(function() {
    driver.getTitle().then(function(title) {
      if (title === 'Transcom PPP: Submit Feedback') {
        console.log('Test passed');
      } else {
        console.log('Test failed');
      }
    });
  });
  driver.quit();
}

load_app_test(driver);
