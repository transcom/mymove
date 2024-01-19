// @ts-check
import { test, expect } from '../../utils/my/customerTest';

test('A customer can go through onboarding', async ({ page, customerPage }) => {
  // Create new customer user
  await customerPage.signInAsNewCustomer();

  // CONUS/OCONUS section
  await customerPage.waitForPage.onboardingConus();
  await page.getByText('Starts and ends in the continental US').click();
  await customerPage.navigateForward();

  // Branch/DOD ID section
  await customerPage.waitForPage.onboardingDodId();
  await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Space Force' });
  await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Army' });
  await page.getByTestId('textInput').fill('1231231234');
  await page.getByRole('combobox', { name: 'Pay grade' }).selectOption({ label: 'E-7' });
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
  // Test changed duty location names
  const changedBaseNames = [
    { baseName: 'Fort Cavazos', baseAddress: 'Fort Cavazos, TX 76544', oldBaseAddress: 'Fort Hood, TX 76544' },
    { baseName: 'Fort Eisenhower', baseAddress: 'Fort Eisenhower, GA 30813', oldBaseAddress: 'Fort Gordon' },
    { baseName: 'Fort Novosel', baseAddress: 'Fort Novosel, AL 36362', oldBaseAddress: 'Fort Rucker, AL 36362' },
    { baseName: 'Fort Gregg-Adams', baseAddress: 'Fort Gregg-Adams, VA 23801', oldBaseAddress: 'Fort Lee' },
  ];

  for (const base of changedBaseNames) {
    await page.getByLabel('What is your current duty location?').fill(base.baseName);
    // click on the base name that pops up in result list
    await page.getByText(base.baseName, { exact: true }).click();
    // verify the duty location that populates the outlined box is the full base address
    const dutyLocationInBox = page.locator('span').filter({ hasText: base.baseAddress });
    await expect(dutyLocationInBox).toHaveText(base.baseAddress);
    // verify the duty location that appears underneath the outlined box is the full base address
    const dutyLocationUnderBox = await page.getByTestId('formGroup').getByRole('paragraph');
    await expect(dutyLocationUnderBox).toHaveText(base.baseAddress);
    await expect(dutyLocationUnderBox).not.toHaveText(base.oldBaseAddress);
  }

  await page.getByLabel('What is your current duty location?').fill('Scott AFB');
  await page.keyboard.press('Backspace'); // tests if backspace clears the duty location field
  await page.getByLabel('What is your current duty location?').fill('Scott AFB');
  // 'mark' is not yet supported by react testing library
  // https://github.com/testing-library/dom-testing-library/issues/1150
  // @ts-expect-error:next-line
  await page.getByRole('mark').nth(0).click();
  await customerPage.navigateForward();

  // Current address section
  await customerPage.waitForPage.onboardingCurrentAddress();
  await page.getByLabel('Address 1').fill('7 Q St');
  await page.getByLabel('City').fill('Atco');
  await page.getByLabel('State').selectOption({ label: 'NJ' });
  await page.getByLabel('ZIP').fill('08004');
  await page.getByLabel('ZIP').blur();
  await customerPage.navigateForward();

  // Backup mailing address section
  await customerPage.waitForPage.onboardingBackupAddress();
  await page.getByLabel('Address 1').fill('7 Q St');
  await page.getByLabel('City').fill('Atco');
  await page.getByLabel('State').selectOption({ label: 'NJ' });
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
