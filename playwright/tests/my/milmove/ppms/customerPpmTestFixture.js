/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import {
  expect,
  test as customerTest,
  forEachViewport,
  useMobileViewport,
  CustomerPage,
} from '../../../utils/customerTest';

/**
 * CustomerPpmPage test fixture
 *
 * @extends CustomerPage
 */
export class CustomerPpmPage extends CustomerPage {
  /**
   * Create an CustomerPpmPage.
   * @param {CustomerPage} customerPage
   */
  constructor(customerPage) {
    super(customerPage.page, customerPage.request);
  }

  /**
   * @param {string} userId
   */
  async signInForPPM(userId) {
    await this.signInAsExistingCustomer(userId);
  }

  /**
   * @param {Object} move
   */
  async signInForPPMWithMove(move) {
    await this.signInAsExistingCustomer(move.Orders.ServiceMember.user_id);
  }

  /**
   * click on upload ppm documents
   */
  async clickOnUploadPPMDocumentsButton() {
    await expect(this.page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();

    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Upload PPM Documents' }).click(),
    ]);
  }

  /**
   */
  async customerStartsAddingAPPMShipment() {
    await this.page.getByRole('button', { name: 'Set up your shipments' }).click();
    await this.navigateForward();

    await this.page.locator('label[for="PPM"]').click();
    await this.navigateForward();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.selectAdvance=false]
   */
  async navigateToAboutPage(options = { selectAdvance: false }) {
    await this.clickOnUploadPPMDocumentsButton();

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/about/);

    await expect(this.page.getByRole('heading', { name: 'About your PPM' })).toBeVisible();

    await this.fillOutAboutPage(options);
  }

  /**
   */
  async navigateToPPMReviewPage() {
    await this.clickOnUploadPPMDocumentsButton();

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);

    await expect(this.page.getByRole('heading', { name: 'Review' })).toBeVisible();
  }

  /**
   */
  async navigateFromPPMReviewPageToFinalCloseoutPage() {
    await this.page.locator('a').getByText('Save & Continue').click();
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/complete/);

    await expect(this.page.getByRole('heading', { name: 'Complete PPM' })).toBeVisible();
  }

  /**
   */
  async navigateToFinalCloseoutPage() {
    await this.navigateToPPMReviewPage();

    await this.navigateFromPPMReviewPageToFinalCloseoutPage();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.isMoveSubmitted=false]
   */
  async navigateFromHomePageToReviewPage(options = { isMoveSubmitted: false }) {
    if (options?.isMoveSubmitted) {
      await expect(this.page.getByRole('heading', { name: 'Next step: Your move gets approved' })).toBeVisible();

      await this.page.getByRole('button', { name: 'Review your request' }).click();
    } else {
      await expect(this.page.getByRole('heading', { name: 'Time to submit your move' })).toBeVisible();

      await this.page.getByRole('button', { name: 'Review and submit' }).click();
    }
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.selectAdvance=false]
   */
  async fillOutAboutPage(options = { selectAdvance: false }) {
    // editing this field with the keyboard instead of the date picker runs async validators for pre-filled postal codes
    // this helps debounce the API calls that would be triggered in quick succession
    await this.page.locator('input[name="actualMoveDate"]').clear();
    await this.page.locator('input[name="actualMoveDate"]').type('01 Feb 2022');
    await this.page.locator('input[name="actualMoveDate"]').blur();

    await this.page.locator('input[name="actualPickupPostalCode"]').clear();
    await this.page.locator('input[name="actualPickupPostalCode"]').type('90210');
    await this.page.locator('input[name="actualPickupPostalCode"]').blur();
    await this.page.locator('input[name="actualDestinationPostalCode"]').clear();
    await this.page.locator('input[name="actualDestinationPostalCode"]').type('76127');
    await this.page.locator('input[name="actualDestinationPostalCode"]').blur();

    if (options?.selectAdvance) {
      await this.page.locator('label[for="yes-has-received-advance"]').click();
      await this.page.locator('input[name="advanceAmountReceived"]').clear();
      await this.page.locator('input[name="advanceAmountReceived"]').type('5000');
    } else {
      await this.page.locator('label[for="no-has-received-advance"]').click();
    }

    await this.page.locator('input[name="w2Address.streetAddress1"]').clear();
    await this.page.locator('input[name="w2Address.streetAddress1"]').type('1819 S Cedar Street');

    await this.page.locator('input[name="w2Address.city"]').clear();
    await this.page.locator('input[name="w2Address.city"]').type('Yuma');
    await this.page.locator('select[name="w2Address.state"]').selectOption({ label: 'AZ' });
    await this.page.locator('input[name="w2Address.postalCode"]').clear();
    await this.page.locator('input[name="w2Address.postalCode"]').type('85369');
    await this.page.locator('input[name="w2Address.postalCode"]').blur();

    await this.page.getByRole('button', { name: 'Save & Continue' }).click();
  }

  /**
   */
  async navigateFromAboutPageToWeightTicketPage() {
    await this.page.getByRole('button', { name: 'Save & Continue' }).click();

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  }

  /**
   */
  async navigateToWeightTicketPage() {
    await this.clickOnUploadPPMDocumentsButton();

    await expect(this.page.getByRole('heading', { name: 'Weight Tickets' })).toBeVisible();
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.hasTrailer=false]
   * @param {boolean} [options.ownTrailer=false]
   * @param {boolean} [options.useConstructedWeight=false]
   */
  async submitWeightTicketPage(options = {}) {
    await this.fillOutWeightTicketPage(options);
    await this.navigateFromWeightTicketPage();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.hasTrailer=false]
   * @param {boolean} [options.ownTrailer=false]
   * @param {boolean} [options.useConstructedWeight=false]
   */
  async fillOutWeightTicketPage(options) {
    const { hasTrailer = false, ownTrailer = false, useConstructedWeight = false } = options;
    await this.page.locator('input[name="vehicleDescription"]').type('Kia Forte');
    await this.page.locator('input[name="vehicleDescription"]').blur();

    await this.page.getByLabel('Empty weight').clear();
    await this.page.getByLabel('Empty weight').type('1000');
    await this.page.getByLabel('Empty weight').blur();
    if (useConstructedWeight) {
      // this page has multiple labels with the same text, grab the
      // first one for Empty Weight. Not sure why getByLabel does not work
      await this.page.locator('label').getByText("I don't have this weight ticket").first().click();

      // this page has multiple labels with the same text, grab the
      // first one for Empty Weight
      const emptyConstructed = this.page.locator('label').getByText('Upload constructed weight spreadsheet').first();

      await expect(emptyConstructed).toBeVisible();
      let filepond = emptyConstructed.locator('../..').locator('.filepond--wrapper');
      await expect(filepond).toBeVisible();

      await this.uploadFileViaFilepond(filepond, 'constructedWeight.xls');

      // wait for the file to be visible in the uploads
      await expect(
        filepond.locator('../..').locator('p').getByText('constructedWeight.xls', { exact: true }),
      ).toBeVisible();

      await this.page.getByLabel('Full weight').clear();
      await this.page.getByLabel('Full weight').type('3000');

      // this page has multiple labels with the same text, grab the
      // second one for Full Weight. Not sure why getByLabel does not work
      await this.page.locator('label').getByText("I don't have this weight ticket").nth(1).click();

      // this page has multiple labels with the same text, grab the
      // second one for Full Weight
      const fullConstructed = this.page.locator('label').getByText('Upload constructed weight spreadsheet').nth(1);

      await expect(fullConstructed).toBeVisible();
      filepond = fullConstructed.locator('../..').locator('.filepond--wrapper');
      await expect(filepond).toBeVisible();

      await this.uploadFileViaFilepond(filepond, 'constructedWeight.xlsx');

      // wait for the file to be visible in the uploads
      await expect(
        filepond.locator('../..').locator('p').getByText('constructedWeight.xlsx', { exact: true }),
      ).toBeVisible();
    } else {
      // find the label, then find the filepond wrapper. Not sure why
      // getByLabel doesn't work
      const emptyWeightLabel = this.page.locator('label').getByText('Upload empty weight ticket', { exact: true });
      await expect(emptyWeightLabel).toBeVisible();
      const emptyFilepond = emptyWeightLabel.locator('../..').locator('.filepond--wrapper');
      await expect(emptyFilepond).toBeVisible();

      await this.uploadFileViaFilepond(emptyFilepond, 'sampleWeightTicket.jpg');

      // wait for the file to be visible in the uploads
      await expect(
        emptyFilepond.locator('../..').locator('p').getByText('sampleWeightTicket.jpg', { exact: true }),
      ).toBeVisible();

      await this.page.getByLabel('Full Weight').clear();
      await this.page.getByLabel('Full Weight').type('3000');

      // find the label, then find the filepond wrapper. Not sure why
      // getByLabel doesn't work
      const fullWeightLabel = this.page.locator('label').getByText('Upload full weight ticket', { exact: true });
      await expect(fullWeightLabel).toBeVisible();
      const fullFilepond = fullWeightLabel.locator('../..').locator('.filepond--wrapper');
      await expect(fullFilepond).toBeVisible();

      await this.uploadFileViaFilepond(fullFilepond, 'sampleWeightTicket.jpg');
      // wait for the file to be visible in the uploads
      await expect(
        fullFilepond.locator('../..').locator('p').getByText('sampleWeightTicket.jpg', { exact: true }),
      ).toBeVisible();
    }

    await expect(this.page.locator('.tripWeightTotal')).toContainText('Trip weight: 2,000 lbs');

    if (hasTrailer) {
      // the page design makes it hard to click without using a css locator
      await this.page.locator('label[for="yesOwnsTrailer"]').click();
      if (ownTrailer) {
        // the page design makes it hard to click without using a css locator
        await this.page.locator('label[for="yestrailerMeetsCriteria"]').click();

        // find the label, then find the filepond wrapper, not sure
        // why getByLabel does not work
        const ownershipLabel = this.page.locator('label').getByText('Upload proof of ownership', { exact: true });
        await expect(ownershipLabel).toBeVisible();
        const ownershipFilepond = ownershipLabel.locator('../..').locator('.filepond--wrapper');
        await expect(ownershipFilepond).toBeVisible();

        await this.uploadFileViaFilepond(ownershipFilepond, 'trailerOwnership.pdf');

        // wait for the file to be visible in the uploads
        await expect(
          ownershipFilepond.locator('../..').locator('p').getByText('trailerOwnership.pdf', { exact: true }),
        ).toBeVisible();
      } else {
        // the page design makes it hard to click without using a css locator
        await this.page.locator('label[for="notrailerMeetsCriteria"]').click();
      }
    }
  }

  /**
   */
  async navigateFromWeightTicketPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Save & Continue' }).click(),
    ]);

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  }

  /**
   */
  async navigateFromHomePageToExistingPPMDateAndLocationPage() {
    await expect(this.page.getByRole('heading', { name: 'Time to submit your move' })).toBeVisible();

    await this.page.locator('[data-testid="shipment-list-item-container"] button').getByText('Edit').click();

    await expect(this.page.getByRole('heading', { name: 'PPM date & location' })).toBeVisible();
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/edit/);
  }

  /**
   * used for creating a new shipment
   */
  async submitsDateAndLocation() {
    await this.page.locator('input[name="pickupPostalCode"]').clear();
    await this.page.locator('input[name="pickupPostalCode"]').type('90210');
    await this.page.locator('input[name="pickupPostalCode"]').blur();

    await this.page.locator('input[name="destinationPostalCode"]').clear();
    await this.page.locator('input[name="destinationPostalCode"]').type('76127');

    await this.page.locator('input[name="expectedDepartureDate"]').clear();
    await this.page.locator('input[name="expectedDepartureDate"]').type('01 Feb 2022');
    await this.page.locator('input[name="expectedDepartureDate"]').blur();

    // Select closeout office
    await this.selectDutyLocation('Fort Bragg', 'closeoutOffice');

    await this.navigateFromDateAndLocationPageToEstimatedWeightsPage();
  }

  /**
   */
  async navigateFromDateAndLocationPageToEstimatedWeightsPage() {
    await this.page.getByRole('button', { name: 'Save & Continue' }).click();

    await expect(this.page.getByRole('heading', { name: 'Estimated weight', exact: true })).toBeVisible();

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  }

  /**
   */
  async submitsEstimatedWeightsAndProGear() {
    await this.page.locator('input[name="estimatedWeight"]').clear();
    await this.page.locator('input[name="estimatedWeight"]').type('4000');

    await this.page.locator('label[for="hasProGearYes"]').click();

    // seems to need to click then clear
    await this.page.locator('input[name="proGearWeight"]').click();
    await this.page.locator('input[name="proGearWeight"]').clear();
    await this.page.locator('input[name="proGearWeight"]').type('500');

    await this.page.locator('input[name="spouseProGearWeight"]').clear();
    await this.page.locator('input[name="spouseProGearWeight"]').type('400');

    await expect(this.page.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await this.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  }

  /**
   */
  async submitsEstimatedWeights() {
    await this.page.locator('input[name="estimatedWeight"]').clear();
    await this.page.locator('input[name="estimatedWeight"]').type('4000');
    await expect(this.page.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await this.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  }

  /**
   */
  async navigateFromEstimatedWeightsPageToEstimatedIncentivePage() {
    await this.page.locator('button').getByText('Save & Continue').click();

    await expect(this.page.getByRole('heading', { name: 'Estimated incentive', exact: true })).toBeVisible();
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/estimated-incentive/);
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.isMobile=false]
   */
  async generalVerifyEstimatedIncentivePage(options = {}) {
    const { isMobile = false } = options;
    await expect(this.page.getByRole('heading', { name: 'Estimated incentive', exact: true })).toBeVisible();

    // checks the format of the incentive amount statement is `$<some
    // comma-separated number without decimals> is`
    await expect(this.page.locator('.container h2')).toContainText(/\$[0-9,]+ is/);

    await expect(this.page.getByRole('button', { name: 'Next' })).toBeEnabled();
    if (isMobile) {
      await expect(this.page.getByRole('button', { name: 'Next' })).toHaveCSS('order', '1');
    }

    await this.navigateFromEstimatedIncentivePageToAdvancesPage();
  }

  /**
   */
  async navigateFromEstimatedIncentivePageToAdvancesPage() {
    await this.page.getByRole('button', { name: 'Next', exact: true }).click();

    await expect(this.page.getByRole('heading', { name: 'Advances', exact: true })).toBeVisible();
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/advances/);
  }

  /**
   * submit ppm advance page
   * @param {Object} options
   * @param {boolean} [options.addAdvance=false]
   * @param {boolean} [options.isMobile=false]
   */
  async submitsAdvancePage(options = {}) {
    const { addAdvance = false, isMobile = false } = options;
    if (addAdvance) {
      await this.page.locator('label[for="hasRequestedAdvanceYes"]').click();

      // not sure why, but need to click then clear
      await this.page.locator('input[name="advanceAmountRequested"]').click();
      await this.page.locator('input[name="advanceAmountRequested"]').clear();
      await this.page.locator('input[name="advanceAmountRequested"]').type('4000');
      await this.page.locator('input[name="advanceAmountRequested"]').blur();

      await this.page.locator('label[for="agreeToTerms"]').click();
    } else {
      await this.page.locator('label[for="hasRequestedAdvanceNo"]').click();
    }

    await this.navigateFromAdvancesPageToReviewPage({ isMobile });
  }

  /**
   * navigate from advances to review
   * @param {Object} options
   * @param {boolean} [options.isMobile=false]
   */
  async navigateFromAdvancesPageToReviewPage(options = {}) {
    const { isMobile = false } = options;

    const saveButton = this.page.getByRole('button', { name: 'Save & Continue' });
    await expect(saveButton).toBeVisible();

    if (isMobile) {
      await expect(saveButton).toHaveCSS('order', '1');
    }

    // when navigating through an existing PPM that requested an
    // advance, we must agree to the terms again to proceed
    const hasAdvance = await this.page.locator('label[for="hasRequestedAdvanceYes"]').isChecked();
    if (hasAdvance) {
      // only look for this if hasAdvance
      const agreedToTerms = await this.page.locator('label[for="agreeToTerms"]').isChecked();
      if (!agreedToTerms) {
        await this.page.locator('label[for="agreeToTerms"]').click();
      }
    }
    await saveButton.click();

    await expect(this.page.getByRole('heading', { name: 'Review your details', exact: true })).toBeVisible();
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/review/);

    await expect(this.page.locator('.usa-alert__heading')).toContainText('Details saved');
    await expect(this.page.locator('.usa-alert__heading + p')).toContainText(
      'Review your info and submit your move request now, or come back and finish later.',
    );
  }

  /**
   */
  async navigateToHomePage() {
    await this.page.getByRole('link', { name: 'Home' }).click();
  }

  /**
   */
  async navigateFromReviewPageToHomePage() {
    await Promise.all([this.page.waitForNavigation(), this.page.getByRole('button', { name: 'Return home' }).click()]);

    expect(new URL(this.page.url()).pathname).toEqual('/');
  }

  /**
   */
  async navigateToFromHomePageToPPMCloseoutReview() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Upload PPM Documents' }).click(),
    ]);

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  }

  /**
   */
  async navigateToAgreementAndSign() {
    await this.navigateForward();
    await this.signAgreement();
  }

  /**
   */
  async signAgreement() {
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/agreement/);
    await expect(this.page.getByRole('heading', { name: 'Now for the official part…' })).toBeVisible();

    await this.page.locator('input[name="signature"]').type('Sofía Clark-Nuñez');
    await expect(this.page.getByRole('button', { name: 'Complete' })).toBeEnabled();
  }

  /**
   */
  async submitMove() {
    await this.page.getByRole('button', { name: 'Complete' }).click();

    await expect(this.page.locator('.usa-alert--success')).toContainText('You’ve submitted your move request.');

    await expect(this.page.getByRole('heading', { name: 'Next step: Your move gets approved' })).toBeVisible();

    // ensure that shipment list doesn't have a button to edit or delete
    await expect(this.page.locator('[data-testid="shipment-list-item-container"] button')).not.toBeVisible();
  }

  /**
   * @param {import('@playwright/test').Locator} locator
   * @param {number} expectedLength
   */
  async deleteShipment(locator, expectedLength) {
    await locator.getByText('Delete').click();
    const modal = this.page.locator('[data-testid="modal"]');
    await expect(modal).toBeVisible();
    await modal.getByRole('button', { name: 'Yes, Delete' }).click();
    await expect(locator).toHaveCount(expectedLength);
    await expect(this.page.locator('[data-testid="alert"]')).toContainText('The shipment was deleted.');
  }

  /**
   * @param {number} ticketIndex
   * @param {boolean} isLastWeightTicket
   */
  async deleteWeightTicket(ticketIndex, isLastWeightTicket) {
    const weightMoved = this.page.getByRole('heading', { name: 'Weight moved' });
    await expect(weightMoved).toBeVisible();
    const weightMovedContainer = weightMoved.locator('../../..');
    await expect(weightMovedContainer).toBeVisible();
    await weightMovedContainer.getByRole('button', { name: 'Delete' }).nth(ticketIndex).click();
    await expect(this.page.getByText(`You are about to delete Trip ${1 + ticketIndex}`)).toBeVisible();
    await this.page.getByRole('button', { name: 'Yes, Delete' }).click();
    if (isLastWeightTicket) {
      await expect(
        this.page.getByText('No weight moved documented. At least one trip is required to continue.'),
      ).toBeVisible();
      await expect(
        this.page.getByText(
          'There are items below that are missing required information. Please select “Edit” to enter all required information or “Delete” to remove the item.',
        ),
      ).toBeVisible();
    }
    await expect(this.page.getByText(`Trip ${1 + ticketIndex} successfully deleted.`)).toBeVisible();
  }

  /**
   * @param {number} index
   * @param {boolean} isLastProGear
   */
  async deleteProGearExpense(index, isLastProGear) {
    const proGearExpense = this.page.getByRole('heading', { name: 'Pro-gear' });
    await expect(proGearExpense).toBeVisible();
    const proGearExpenseContainer = proGearExpense.locator('../../..');
    await expect(proGearExpenseContainer).toBeVisible();
    await proGearExpenseContainer.getByRole('button', { name: 'Delete' }).nth(index).click();
    await expect(this.page.getByText(`You are about to delete Set ${1 + index}`)).toBeVisible();
    await this.page.getByRole('button', { name: 'Yes, Delete' }).click();
    await expect(this.page.getByText(`Set ${1 + index} successfully deleted.`)).toBeVisible();
    if (isLastProGear) {
      await expect(this.page.getByText('No pro-gear weight documented.')).toBeVisible();
    }
  }

  /**
   * @param {number} index
   * @param {boolean} isLastMovingExpense
   */
  async deleteMovingExpense(index, isLastMovingExpense) {
    const moveExpense = this.page.getByRole('heading', { name: 'Expenses' });
    await expect(moveExpense).toBeVisible();
    const moveExpensesContainer = moveExpense.locator('../../..');
    await expect(moveExpensesContainer).toBeVisible();
    await moveExpensesContainer.getByRole('button', { name: 'Delete' }).nth(index).click();
    await expect(this.page.getByText(`You are about to delete Receipt ${index + 1}`)).toBeVisible();
    await this.page.getByRole('button', { name: 'Yes, Delete' }).click();
    await expect(this.page.getByText(`Receipt ${index + 1} successfully deleted.`)).toBeVisible();
    if (isLastMovingExpense) {
      await expect(this.page.getByText('No receipts uploaded.')).toBeVisible();
    }
  }

  /**
   * @param {string[][]} shipmentCardFields
   * @param {Object} options
   * @param {boolean} options.isEditable=false
   */
  async verifyPPMShipmentCard(shipmentCardFields, options = { isEditable: false }) {
    const { isEditable = false } = options;
    // get first div after the move setup heading
    const ppm1 = this.page.locator(':text("Move setup") + div');
    await expect(ppm1).toBeVisible();

    if (isEditable) {
      await expect(ppm1.getByRole('button', { name: 'Edit' })).toBeVisible();
      await expect(ppm1.getByRole('button', { name: 'Delete' })).toBeVisible();
    } else {
      await expect(ppm1.locator('[data-testid="ShipmentContainer"]').locator('button')).not.toBeVisible();
    }

    await expect(ppm1.locator('dt')).toHaveCount(shipmentCardFields.length);
    await expect(ppm1.locator('dd')).toHaveCount(shipmentCardFields.length);

    shipmentCardFields.forEach(async (shipmentField, index) => {
      await expect(ppm1.locator('dt').nth(index)).toContainText(shipmentField[0]);
      await expect(ppm1.locator('dd').nth(index)).toContainText(shipmentField[1]);
    });
  }

  /**
   */
  async navigateToProgearPage() {
    await this.navigateToPPMReviewPage();
    await this.navigateFromCloseoutReviewPageToProGearPage();
  }

  /**
   */
  async navigateFromCloseoutReviewPageToProGearPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('link', { name: 'Add Pro-gear Weight' }).click(),
    ]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToEditProGearPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.locator('.progearSection a').getByText('Edit').click(),
    ]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToAddProGearPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('link', { name: 'Add Pro-gear Weight' }).click(),
    ]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToEditWeightTicketPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.locator('.reviewWeightTickets a').getByText('Edit').click(),
    ]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToAddWeightTicketPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('link', { name: 'Add More Weight' }).click(),
    ]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  }

  /**
   */
  async cancelAddLineItemAndReturnToCloseoutReviewPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Return to Homepage' }).click(),
    ]);
    expect(new URL(this.page.url()).pathname).toBe('/');
    await this.navigateToPPMReviewPage();
  }

  /**
   */
  async navigateFromCloseoutReviewPageToEditExpensePage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.locator('.reviewExpenses a').getByText('Edit').click(),
    ]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/expenses/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToAddExpensePage() {
    await Promise.all([this.page.waitForNavigation(), this.page.getByRole('link', { name: 'Add Expenses' }).click()]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/expenses/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToAboutPage() {
    await this.page.locator('[data-testid="aboutYourPPM"] a').getByText('Edit').click();
  }

  /**
   */
  async navigateFromProgearPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Save & Continue' }).click(),
    ]);

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  }

  async verifySaveAndContinueDisabled() {
    await expect(this.page.getByRole('link', { name: 'Save & Continue' })).toBeDisabled();
    await expect(
      this.page.getByText(
        'There are items below that are missing required information. Please select “Edit” to enter all required information or “Delete” to remove the item.',
      ),
    ).toBeVisible();
  }

  async verifySaveAndContinueEnabled() {
    await expect(this.page.getByRole('link', { name: 'Save & Continue' })).toBeEnabled();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.belongsToSelf=true]
   * @param {string} [options.weight]
   * @param {boolean} [options.missingWeightTicket]
   */
  async submitProgearPage(options = { belongsToSelf: true }) {
    await this.fillOutProgearPage(options);
    await this.navigateFromProgearPage();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.belongsToSelf=true]
   * @param {string} [options.weight]
   * @param {boolean} [options.missingWeightTicket]
   */
  async fillOutProgearPage(options) {
    const belongs = options?.belongsToSelf ? 'Self' : 'Spouse';
    const progearTypeSelector = `label[for="ownerOfProGear${belongs}"]`;
    await this.page.locator(progearTypeSelector).click();

    await this.page.locator('[name="description"]').clear();
    await this.page.locator('[name="description"]').type('Radio equipment');

    // need to click before clear for some reason
    await this.page.locator('[name="weight"]').click();
    await this.page.locator('[name="weight"]').clear();
    if (options?.belongsToSelf) {
      await this.page.locator('[name="weight"]').type(options?.weight || '2000');
    } else {
      await this.page.locator('[name="weight"]').type(options?.weight || '500');
    }

    let uploadFilename = 'sampleWeightTicket.jpg';
    let upload;

    if (options?.missingWeightTicket) {
      await this.page.locator('label').getByText("I don't have this weight ticket").click();
      uploadFilename = 'constructedWeight.xls';
      upload = this.page.locator('label').getByText('Upload constructed weight spreadsheet');
    } else {
      upload = this.page.locator('label').getByText("Upload your pro-gear's weight tickets");
    }

    await expect(upload).toBeVisible();
    const filepond = upload.locator('../..').locator('.filepond--wrapper');
    await expect(filepond).toBeVisible();
    await this.uploadFileViaFilepond(filepond, uploadFilename);

    // wait for the file to be visible in the uploads
    await expect(filepond.locator('../..').locator('p').getByText(uploadFilename, { exact: true })).toBeVisible();
  }

  /**
   */
  async navigateFromCloseoutReviewPageToExpensesPage() {
    await Promise.all([this.page.waitForNavigation(), this.page.getByRole('link', { name: 'Add Expense' }).click()]);
    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/expenses/);
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.isEditExpense=false]
   * @param {string} [options.amount]
   */
  async submitExpensePage(options = { isEditExpense: false }) {
    const expenseType = this.page.locator('select[name="expenseType"]');
    if (!options?.isEditExpense) {
      await expect(expenseType).toHaveValue('');
    }
    await expenseType.selectOption({ label: 'Storage' });

    await this.page.locator('input[name="description"]').type('Cloud storage');
    await this.page.locator('label[for="yes-used-gtcc"]').click();
    await this.page.locator('input[name="amount"]').clear();
    await this.page.locator('input[name="amount"]').type(options?.amount || '675.99');

    // find the label, then find the filepond wrapper. Not sure why
    // getByLabel doesn't work
    const fullWeightLabel = this.page.locator('label').getByText('Upload receipt', { exact: true });
    await expect(fullWeightLabel).toBeVisible();
    const fullFilepond = fullWeightLabel.locator('../..').locator('.filepond--wrapper');
    await expect(fullFilepond).toBeVisible();

    await this.uploadFileViaFilepond(fullFilepond, 'sampleWeightTicket.jpg');
    // wait for the file to be visible in the uploads
    await expect(
      fullFilepond.locator('../..').locator('p').getByText('sampleWeightTicket.jpg', { exact: true }),
    ).toBeVisible();

    await this.page.locator('input[name="sitStartDate"]').type('14 Aug 2022');
    await this.page.locator('input[name="sitStartDate"]').blur();
    await this.page.locator('input[name="sitEndDate"]').type('20 Aug 2022');
    await this.page.locator('input[name="sitEndDate"]').blur();

    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Save & Continue' }).click(),
    ]);

    await expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);

    const cloudStorage = this.page.getByText('Cloud storage');
    await expect(cloudStorage).toBeVisible();
    const receiptContainer = cloudStorage.locator('../..');
    await expect(receiptContainer.locator('dt').getByText('Days in storage:')).toBeVisible();
    await expect(receiptContainer.locator('dd').getByText('7', { exact: true })).toBeVisible();
  }

  /**
   * @param {Object} options
   * @param {string} [options.totalNetWeight='4,000 lbs']
   * @param {string} [options.proGearWeight='1,500 lbs']
   * @param {string} [options.expensesClaimed='450.00']
   * @param {string} [options.finalIncentiveAmount='$500,000.00']
   */
  async verifyFinalIncentiveAndTotals(
    options = {
      totalNetWeight: '4,000 lbs',
      proGearWeight: '1,500 lbs',
      expensesClaimed: '450.00',
      finalIncentiveAmount: '$500,000.00',
    },
  ) {
    await expect(
      this.page.getByRole('heading', { name: `Your final estimated incentive: ${options?.finalIncentiveAmount}` }),
    ).toBeVisible();

    await expect(this.page.locator('li').getByText(`${options?.totalNetWeight} total net weight`)).toBeVisible();

    // TODO: Once we get moving expenses and pro gear back, check for those here as well.

    await expect(this.page.locator('li').getByText(`${options?.proGearWeight} of pro-gear`)).toBeVisible();
    await expect(this.page.locator('li').getByText(`$${options?.expensesClaimed} in expenses claimed`)).toBeVisible();
  }

  /**
   */
  async signCloseoutAgreement() {
    await this.page.locator('input[name="signature"]').type('Sofía Clark-Nuñez');
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Submit PPM Documentation' }).click(),
    ]);

    expect(new URL(this.page.url()).pathname).toEqual('/');

    await expect(this.page.locator('.usa-alert--success')).toContainText('You submitted documentation for review.');

    const stepContainer = this.page.locator('[data-testid="stepContainer5"]');
    await expect(stepContainer.getByRole('button', { name: 'Download Incentive Packet' })).toBeDisabled();
    await expect(stepContainer.getByText(/PPM documentation submitted: \d{2} \w{3} \d{4}/)).toBeVisible();
  }

  /**
   * @param {Object} options
   * @param {string} [options.totalNetWeight='4,000 lbs']
   * @param {string} [options.proGearWeight='1,500 lbs']
   * @param {string} [options.expensesClaimed='450.00']
   * @param {string} [options.finalIncentiveAmount='$500,000.00']
   */
  async submitFinalCloseout(options) {
    await this.verifyFinalIncentiveAndTotals(options);
    await this.signCloseoutAgreement();
  }
}

/**
 * @typedef {object} CustomerPpmPageTestArgs - customer ppm page test args
 * @property {CustomerPpmPage} customerPpmPage    - customer ppm page
 */

/** @type {import('@playwright/test').Fixtures<CustomerPpmPageTestArgs, {}, import('../../../utils/customerTest').CustomerPageTestArgs, import('@playwright/test').PlaywrightWorkerArgs>} */
const customerPpmFixtures = {
  customerPpmPage: async ({ customerPage }, use) => {
    const customerPpmPage = new CustomerPpmPage(customerPage);
    await use(customerPpmPage);
  },
};

export const test = customerTest.extend(customerPpmFixtures);

export { expect, forEachViewport, useMobileViewport };

export default CustomerPpmPage;
