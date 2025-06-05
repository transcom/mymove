// @ts-check
// import { test, expect } from '../../utils/my/customerTest';
import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;
const LocationLookup = 'ATCO, NJ 08004 (CAMDEN)';

test.describe('Onboarding', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  test('A customer can go through onboarding', async ({ page, customerPage }) => {
    // Create new customer user
    await customerPage.signInAsNewCustomer();

    // Branch/DOD ID section
    await customerPage.waitForPage.onboardingDodId();
    await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Space Force' });
    await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Army' });
    await page.getByTestId('textInput').fill('1231231234');
    await customerPage.navigateForward();

    // Name secton
    await customerPage.waitForPage.onboardingName();
    await page.getByLabel('First name').fill('Leo');
    await page.getByLabel('Last name').fill('Spacemen');
    await customerPage.navigateForward();

    // Contact info section
    await customerPage.waitForPage.onboardingContactInfo();
    await page.getByLabel('Best contact phone').fill('2025552345');
    await page.getByText('Email', { exact: true }).click();
    await customerPage.navigateForward();

    // Current address section
    await customerPage.waitForPage.onboardingCurrentAddress();
    await page.getByLabel('Address 1').fill('7 Q St');
    await page.locator('input[id="current_residence-input"]').fill('08004');
    await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Backup mailing address section
    await customerPage.waitForPage.onboardingBackupAddress();
    await page.getByLabel('Address 1').fill('7 Q St');
    await page.locator('input[id="backup_mailing_address-input"]').fill('08004');
    await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Backup contact info section
    await customerPage.waitForPage.onboardingBackupContact();
    await page.getByLabel('First Name').fill('Grace');
    await page.getByLabel('Last Name').fill('Griffin');
    await page.getByLabel('Email').fill('grace.griffin@example.com');
    await page.getByLabel('Phone').fill('2025553456');
    await customerPage.navigateForward();

    await customerPage.waitForPage.home();
  });
});

test.describe('(MultiMove) Onboarding', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('A customer can go through onboarding', async ({ page, customerPage }) => {
    // Create new customer user
    await customerPage.signInAsNewCustomer();

    // Branch/DOD ID/Rank section
    await customerPage.waitForPage.onboardingDodId();
    await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Space Force' });
    await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Army' });
    await customerPage.navigateForward();

    // Name secton
    await customerPage.waitForPage.onboardingName();
    await page.getByLabel('First name').fill('Leo');
    await page.getByLabel('Last name').fill('Spacemen');
    await customerPage.navigateForward();

    // Contact info section
    await customerPage.waitForPage.onboardingContactInfo();
    await page.getByLabel('Best contact phone').fill('2025552345');
    await page.getByText('Email', { exact: true }).click();
    await customerPage.navigateForward();

    // Current address section
    await customerPage.waitForPage.onboardingCurrentAddress();
    await page.getByLabel('Address 1').fill('7 Q St');
    await page.getByLabel('Address 1').blur();
    await page.locator('input[id="current_residence-input"]').fill('08004');
    await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Backup mailing address section
    await customerPage.waitForPage.onboardingBackupAddress();
    await page.getByLabel('Address 1').fill('7 Q St');
    await page.getByLabel('Address 1').blur();
    await page.locator('input[id="backup_mailing_address-input"]').fill('08004');
    await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Backup contact info section
    await customerPage.waitForPage.onboardingBackupContact();
    await page.getByLabel('First Name').fill('Grace');
    await page.getByLabel('Last Name').fill('Griffin');
    await page.getByLabel('Email').fill('grace.griffin@example.com');
    await page.getByLabel('Phone').fill('2025553456');
    await customerPage.navigateForward();

    await customerPage.waitForPage.multiMoveDashboard();
  });
});
