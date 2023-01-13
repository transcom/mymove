/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect } = require('../../utils/officeTest');

test.describe('Customer Support User Flows', () => {
  test.describe('Customer support remarks', () => {
    test('is able to add, edit, and delete a remark', async ({ page, officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
      const moveLocator = move.locator;

      await officePage.signInAsNewQAECSRUser();
      await officePage.qaeCsrSearchForAndNavigateToMove(moveLocator);

      // Go to Customer support remarks
      await page.getByText('Customer support remarks').click();
      await officePage.waitForLoading();
      expect(page.url()).toContain(`/moves/${moveLocator}/customer-support-remarks`);
      await expect(page.getByText('Past remarks')).toBeVisible();

      // Validate remarks page content
      await expect(page.locator('h1')).toContainText('Customer support remarks');
      await expect(page.locator('h2').getByText('Remarks', { exact: true })).toBeVisible();
      await expect(page.locator('h3')).toContainText('Past remarks');
      await expect(page.locator('small')).toContainText(
        'Use this form to document any customer support provided for this move.',
      );
      await expect(page.locator('[data-testid="textarea"]')).toHaveAttribute('placeholder', 'Add your remarks here');

      await expect(page.locator('[data-testid=form] > [data-testid=button]')).toBeDisabled();

      // Should not have remarks (yet)
      await expect(page.getByText('No remarks yet')).toBeVisible();

      // Add a remark
      const testRemarkText = 'This is a test remark';
      const editString = '-edit';
      await page.locator('[data-testid="textarea"]').type(testRemarkText);
      await expect(page.locator('[data-testid=form] > [data-testid=button]')).toBeEnabled();
      await page.locator('[data-testid=form] > [data-testid=button]').click();
      await expect(page.getByText('No remarks yet')).toHaveCount(0);
      await expect(page.getByText(testRemarkText)).toBeVisible();

      // Open delete modal
      await expect(page.locator('[data-testid="modal"]')).toHaveCount(0);
      await page.locator('[data-testid="delete-remark-button"]').click();
      await expect(page.locator('[data-testid="modal"]')).toBeVisible();
      await expect(page.getByText('Are you sure you want to delete this remark')).toBeVisible();
      await expect(page.getByText('You cannot undo this action')).toBeVisible();
      await expect(page.getByText('Yes, Delete')).toBeVisible();
      await expect(page.getByText('No, keep it')).toBeVisible();

      // Exit modal with cancel button
      await page.locator('[data-testid=modalBackButton]').click();

      // Open the delete modal again
      await page.locator('[data-testid="delete-remark-button"]').click();

      // Exit modal with the X button
      await page.locator('[data-testid=modalCloseButton]').click();

      // Delete the remark for real
      await page.locator('[data-testid="delete-remark-button"]').click();
      await page.getByText('Yes, Delete').click();

      // Make sure success alert is shown
      await expect(page.getByText('Your remark has been deleted')).toHaveCount(1);

      // Validate that the deleted remark is not on the page
      await expect(page.getByText(testRemarkText)).toHaveCount(0);
      await expect(page.getByText('No remarks yet')).toBeVisible();

      // Add a new remark
      await page.locator('[data-testid="textarea"]').type(testRemarkText);
      await expect(page.locator('[data-testid=form] > [data-testid=button]')).toBeEnabled();
      await page.locator('[data-testid=form] > [data-testid=button]').click();

      // Open edit and cancel
      await page.locator('[data-testid="edit-remark-button"]').click();
      await page.locator('[data-testid="edit-remark-textarea"]').type(editString);
      await page.locator('[data-testid="edit-remark-cancel-button"]').click();

      // Validate remark was not edited
      await expect(page.getByText(testRemarkText)).toHaveCount(1);
      await expect(page.getByText(testRemarkText + editString)).not.toBeVisible();

      // Edit the remark
      await page.locator('[data-testid="edit-remark-button"]').click();
      await page.locator('[data-testid="edit-remark-textarea"]').type(testRemarkText + editString);

      // Save the remark edit
      await page.locator('[data-testid="edit-remark-save-button"]').click();

      // Validate remark was edited
      await expect(page.getByText(testRemarkText + editString)).toHaveCount(1);
      await expect(page.getByText('(edited)')).toBeVisible();

      // Change user
      await Promise.all([page.waitForNavigation(), await page.getByText('Sign out').click()]);

      await officePage.signInAsNewQAECSRUser();
      await officePage.qaeCsrSearchForAndNavigateToMove(moveLocator);

      // Go to Customer support remarks
      await page.getByText('Customer support remarks').click();
      await officePage.waitForLoading();

      // Edited remark should exist but no edit/delete buttons as I am a different user
      await expect(page.getByText(testRemarkText + editString)).toBeVisible();
      await expect(page.locator('[data-testid="edit-remark-button"]')).toHaveCount(0);
      await expect(page.locator('[data-testid="delete-remark-button"]')).toHaveCount(0);
    });
  });

  test.describe('Permission based access', () => {
    test('is able to see orders and form is read only', async ({ page, officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
      const moveLocator = move.locator;

      await officePage.signInAsNewQAECSRUser();
      await officePage.qaeCsrSearchForAndNavigateToMove(moveLocator);

      // Navigate to view orders page
      await page.locator('[data-testid="view-orders"]').getByText('View orders').click();

      await expect(page.locator('input[name="issueDate"]')).toBeDisabled();
      await expect(page.locator('input[name="reportByDate"]')).toBeDisabled();
      await expect(page.locator('select[name="departmentIndicator"]')).toBeDisabled();
      await expect(page.locator('input[name="ordersNumber"]')).toBeDisabled();
      await expect(page.locator('select[name="ordersType"]')).toBeDisabled();
      await expect(page.locator('select[name="ordersTypeDetail"]')).toBeDisabled();
      await expect(page.locator('input[name="tac"]')).toBeDisabled();
      await expect(page.locator('input[name="sac"]')).toBeDisabled();
      // no save button should exist
      await expect(page.getByRole('button', { name: 'Save' })).toHaveCount(0);
    });

    test('is able to see allowances and the form is read only', async ({ page, officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
      const moveLocator = move.locator;

      await officePage.signInAsNewQAECSRUser();
      await officePage.qaeCsrSearchForAndNavigateToMove(moveLocator);

      // Navigate to view allowances page
      await page.locator('[data-testid="view-allowances"]').getByText('View allowances').click();

      // read only pro-gear, pro-gear spouse, RME, SIT, and OCIE fields
      await expect(page.locator('input[name="proGearWeight"]')).toBeDisabled();
      await expect(page.locator('input[name="proGearWeightSpouse"]')).toBeDisabled();
      await expect(page.locator('input[name="requiredMedicalEquipmentWeight"]')).toBeDisabled();
      await expect(page.locator('input[name="storageInTransit"]')).toBeDisabled();
      await expect(page.locator('input[name="organizationalClothingAndIndividualEquipment"]')).toBeDisabled();

      // read only grade and authorized weight
      await expect(page.locator('select[name=agency]')).toBeDisabled();
      await expect(page.locator('select[name=agency]')).toBeDisabled();
      await expect(page.locator('select[name="grade"]')).toBeDisabled();
      await expect(page.locator('select[name="grade"]')).toBeDisabled();
      await expect(page.locator('input[name="authorizedWeight"]')).toBeDisabled();
      await expect(page.locator('input[name="dependentsAuthorized"]')).toBeDisabled();

      // no save button should exist
      await expect(page.getByRole('button', { name: 'Save' })).toHaveCount(0);
    });
  });
});
