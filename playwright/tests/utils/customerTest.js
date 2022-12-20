const base = require('@playwright/test');

class CustomerPage {
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }
}

exports.test = base.test.extend({
  customerPage: async ({ page, request }, use) => {
    const customerPage = new CustomerPage(page, request);
    await use(customerPage);
  },
});

exports.expect = base.expect;
