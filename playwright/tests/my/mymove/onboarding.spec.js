// @ts-check
import { test } from '../../utils/customerTest';

test('A customer can go through onboarding', async ({ page, customerPage }) => {
  // Create new customer user
  await customerPage.signInAsNewCustomer();

  // CONUS/OCONUS section
  await customerPage.waitForPage.onboardingConus();
  await page.getByText('Starts and ends in the continental US').click();
  await customerPage.navigateForward();

  // Branch/DOD ID/Rank section
  await customerPage.waitForPage.onboardingDodId();
  await page.getByRole('combobox', { name: 'Branch of service' }).selectOption('ARMY');
  await page.getByTestId('textInput').fill('1231231234');
  await page.getByRole('combobox', { name: 'Rank' }).selectOption('E_7');
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

  // Current duty location section
  await customerPage.waitForPage.onboardingDutyLocation();
  await page.getByLabel('What is your current duty location?').fill('Scott AFB');
  // 'mark' is not yet supported by react testing library
  // https://github.com/testing-library/dom-testing-library/issues/1150
  // @ts-expect-error:next-line
  await page.getByRole('mark').click();
  await customerPage.navigateForward();

  // Current mailing address section
  await customerPage.waitForPage.onboardingCurrentAddress();
  await page.getByLabel('Address 1').fill('7 Q St');
  await page.getByLabel('City').fill('Atco');
  await page.getByLabel('State').selectOption('NJ');
  await page.getByLabel('ZIP').fill('08004');
  await page.getByLabel('ZIP').blur();
  await customerPage.navigateForward();

  // Backup mailing address section
  await customerPage.waitForPage.onboardingBackupAddress();
  await page.getByLabel('Address 1').fill('7 Q St');
  await page.getByLabel('City').fill('Atco');
  await page.getByLabel('State').selectOption('NJ');
  await page.getByLabel('ZIP').fill('08004');
  await page.getByLabel('ZIP').blur();
  await customerPage.navigateForward();

  // Backup contact info section
  await customerPage.waitForPage.onboardingBackupContact();
  await page.getByLabel('Name').fill('Grace Griffin');
  await page.getByLabel('Email').fill('grace.griffin@example.com');
  await page.getByLabel('Phone').fill('2025553456');
  await customerPage.navigateForward();

  await customerPage.waitForPage.home();
});
