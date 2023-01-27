/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, setMobileViewport, CustomerPage } = require('../../../utils/customerTest');

/**
 * CustomerPpmPage test fixture
 *
 * @extends CustomerPage
 */
export class CustomerPpmPage extends CustomerPage {
  /**
   * Create an CustomerPpmPage.
   * @param {CustomerPage} customerPage
   * @param {object} move
   */
  constructor(customerPage, move) {
    super(customerPage.page, customerPage.request);
    this.move = move;
    this.userId = move.Orders.ServiceMember.user_id;
    if (this.move.MTOShipments) {
      this.firstShipmentId = this.move.MTOShipments[0].ID;
    }
  }

  /**
   */
  async signInForPPM() {
    await this.signInAsExistingCustomer(this.userId);
  }

  /**
   * sign in for ppm
   */
  async signInAndClickOnUploadPPMDocumentsButton() {
    await this.signInForPPM();
    await expect(this.page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();

    await this.page.getByRole('button', { name: 'Upload PPM Documents' }).click();
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
   * @param {boolean} [options.isMoveSubmitted=false]
   */
  async signInAndNavigateFromHomePageToReviewPage(options = { isMoveSubmitted: false }) {
    await this.signInAsExistingCustomer(this.userId);

    await this.navigateFromHomePageToReviewPage(options);
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.selectAdvance=false]
   */
  async signInAndNavigateToAboutPage(options = { selectAdvance: false }) {
    await this.signInAndClickOnUploadPPMDocumentsButton();

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/about`;

    expect(this.page.url()).toContain(url);

    await expect(this.page.getByRole('heading', { name: 'About your PPM' })).toBeVisible();

    await this.fillOutAboutPage(options);
  }

  /**
   */
  async signInAndNavigateToPPMReviewPage() {
    await this.signInAndClickOnUploadPPMDocumentsButton();

    expect(this.page.url()).toContain('/moves/[^/]+/shipments/[^/]+/review/');

    await expect(this.page.getByRole('heading', { name: 'Review' })).toBeVisible();
  }

  /**
   */
  async navigateFromPPMReviewPageToFinalCloseoutPage() {
    await this.page.locator('a').getByText('Save & Continue').click();
    expect(this.page.url()).toContain('/moves/[^/]+/shipments/[^/]+/complete/');

    await expect(this.page.getByRole('heading', { name: 'Complete PPM' })).toBeVisible();
  }

  /**
   */
  async signInAndNavigateToFinalCloseoutPage() {
    await this.signInAndNavigateToPPMReviewPage();

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
    await this.page.locator('input[name="actualMoveDate"]').type('01 Feb 2022');
    await this.page.locator('input[name="actualMoveDate"]').blur();

    await this.page.locator('input[name="actualPickupPostalCode"]').type('90210');
    await this.page.locator('input[name="actualPickupPostalCode"]').blur();
    await this.page.locator('input[name="actualDestinationPostalCode"]').type('76127');
    await this.page.locator('input[name="actualDestinationPostalCode"]').blur();

    if (options?.selectAdvance) {
      await this.page.locator('label[for="hasRequestedAdvanceYes"]').click();
      await this.page.locator('input[name="advanceAmountReceived"]').clear();
      await this.page.locator('input[name="advanceAmountReceived"]').type('5000');
    } else {
      await this.page.locator('label[for="hasRequestedAdvanceNo"]').click();
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

    expect(this.page.url()).toContain('/moves/[^/]+/shipments/[^/]+/weight-tickets/');
  }

  /**
   */
  async signInAndNavigateToWeightTicketPage() {
    await this.signInAndClickOnUploadPPMDocumentsButton();

    await expect(this.page.getByRole('heading', { name: 'Weight Tickets' })).toBeVisible();
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/weight-tickets`;
    expect(this.page.url()).toContain(url);
  }

  /**
   * @param {object} options
   */
  async submitWeightTicketPage() {
    // TODO: DREW DEBUG FIXME
    // await this.fillOutWeightTicketPage(options);
    await this.navigateFromWeightTicketPage();
  }

  //   /**
  //    * @param {object} options
  //    */
  //  fillOutWeightTicketPage(options) {
  //   await this.page.locator('input[name="vehicleDescription"]').type('Kia Forte'); await this.page.locator('input[name="vehicleDescription"]').blur();

  //   if (options?.useConstructedWeight) {
  //     await this.page.locator('input[name="emptyWeight"]').type('1000'); await this.page.locator('input[name="emptyWeight"]').blur();
  //     await this.page.locator('input[name="missingEmptyWeightTicket"]').check();

  //     cy.upload_file('.emptyDocument.filepond--root', 'constructedWeight.xls');

  //     await expect(this.page.locator('[data-filepond-item-state="processing-complete"]')).toBeVisible();
  //     await this.page.locator('input[name="fullWeight"]').type('3000');
  //     await this.page.locator('input[name="missingFullWeightTicket"]').check();

  //     cy.upload_file('.fullDocument.filepond--root', 'constructedWeight.xlsx');
  //     cy.wait('@uploadFile');

  //     await expect(this.page.locator('[data-filepond-item-state="processing-complete"]')).toBeVisible();
  //   } else {
  //     await this.page.locator('input[name="emptyWeight"]').type('1000'); await this.page.locator('input[name="emptyWeight"]').blur();

  //     cy.upload_file('.emptyDocument.filepond--root', 'sampleWeightTicket.jpg');
  //     await expect(this.page.locator('[data-filepond-item-state="processing-complete"]')).toBeVisible();

  //     await page.locator('input[name="fullWeight"]').type('3000');

  //     cy.upload_file('.fullDocument.filepond--root', 'sampleWeightTicket.jpg');
  //     await expect(this.page.locator('[data-filepond-item-state="processing-complete"]')).toBeVisible();
  //   }

  //   await expect(this.page.locator('.tripWeightTotal')).toContainText('Trip weight: 2,000 lbs');

  //   if (options?.hasTrailer) {
  //     await this.page.locator('input[name="ownsTrailer"][value="true"]').check();
  //     if (options?.ownTrailer) {
  //       await this.page.locator('input[name="trailerMeetsCriteria"][value="true"]').check();

  //       cy.upload_file('.proofOfTrailerOwnershipDocument.filepond--root', 'trailerOwnership.pdf');
  //     await expect(this.page.locator('[data-filepond-item-state="processing-complete"]')).toBeVisible();
  //     } else {
  //       await this.page.locator('input[name="trailerMeetsCriteria"][value="false"]').check();
  //     }
  //   }
  // }

  /**
   */
  async navigateFromWeightTicketPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Save & Continue' }).click(),
    ]);

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/review`;
    expect(this.page.url()).toContain(url);
  }

  /**
   */
  async signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage() {
    await this.signInForPPM();

    await expect(this.page.getByRole('heading', { name: 'Time to submit your move' })).toBeVisible();

    await this.page.locator('[data-testid="shipment-list-item-container"] button').getByText('Edit').click();

    await expect(this.page.getByRole('heading', { name: 'PPM date & location' })).toBeVisible();
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/edit`;
    expect(this.page.url()).toContain(url);
  }

  /**
   * used for creating a new shipment
   */
  async submitsDateAndLocation() {
    await this.page.locator('input[name="pickupPostalCode"]').type('90210');
    await this.page.locator('input[name="pickupPostalCode"]').blur();

    await this.page.locator('input[name="destinationPostalCode"]').type('76127');
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

    expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  }

  /**
   */
  async submitsEstimatedWeightsAndProGear() {
    await this.page.locator('input[name="estimatedWeight"]').type('4000');
    await this.page.locator('input[name="estimatedWeight"]').blur();
    await this.page.locator('input[name="hasProGear"][value="true"]').check({ force: true });
    await this.page.locator('input[name="proGearWeight"]').type('500');
    await this.page.locator('input[name="proGearWeight"]').blur();
    await this.page.locator('input[name="spouseProGearWeight"]').type('400');
    await this.page.locator('input[name="spouseProGearWeight"]').blur();
    await this.page.getByRole('button', { name: 'Save & Continue' }).click();

    await this.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  }

  /**
   */
  async submitsEstimatedWeights() {
    await this.page.locator('input[name="estimatedWeight"]').type('4000');
    await this.page.locator('input[name="estimatedWeight"]').blur();
    await this.page.getByRole('button', { name: 'Save & Continue' }).click();

    await this.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  }

  /**
   */
  async navigateFromEstimatedWeightsPageToEstimatedIncentivePage() {
    await this.page.locator('button').getByText('Save & Continue').click();

    await expect(this.page.getByRole('heading', { name: 'Estimated incentive', exact: true })).toBeVisible();
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/estimated-incentive`;
    expect(this.page.url()).toContain(url);
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
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/advances`;
    expect(this.page.url()).toContain(url);
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

    expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  }

  /**
   */
  async navigateToAgreementAndSign() {
    await this.navigateForward();
    // TODO: DREW DEBUG FIXME!!!
    // await this.signAgreement();
  }

  // export function deleteShipment(selector, expectedLength) {
  //   await expect(this.page.locator(selector)).toContainText('Delete').click();
  //   await this.page.locator('[data-testid="modal"]').within(($modal) => {
  //     expect($modal).to.be.visible;
  //     await expect(this.page.locator('button')).toContainText('Yes, Delete').click();
  //   });
  //   cy.watest(['@deleteShipment', async ({page}) => {
  //   if (expectedLength > 0) {
  //     await this.page.locator(selector).should('have.length', expectedLength);
  //   } else {
  //     await this.page.locator(selector).should('not.exist');
  //   }
  //   await this.page.locator('[data-testid="alert"]').should('contain', 'The shipment was deleted.');
  // }

  /**
   */
  async signInAndNavigateToProgearPage() {
    await this.signInAndNavigateToPPMReviewPage();
    await this.navigateFromCloseoutReviewPageToProGearPage();
  }

  /**
   */
  async navigateFromCloseoutReviewPageToProGearPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      await this.page.getByRole('button', { name: 'Add Pro-gear Weight' }).click(),
    ]);
    expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToEditProGearPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      await this.page.locator('.progearSection a').getByText('Edit').click(),
    ]);
    expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToEditWeightTicketPage() {
    await Promise.all([
      this.page.waitForNavigation(),
      await this.page.locator('.reviewWeightTickets a').getByText('Edit').click(),
    ]);
    expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  }

  /**
   */
  async navigateFromCloseoutReviewPageToEditExpensePage() {
    await Promise.all([
      this.page.waitForNavigation(),
      this.page.locator('.reviewExpenses a').getByText('Edit').click(),
    ]);
    expect(this.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/expenses/);
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

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/review`;
    expect(this.page.url()).toContain(url);
  }

  // export function submitProgearPage(options) {
  //   fillOutProgearPage(options);
  //   navigateFromProgearPage();
  // }

  // export function fillOutProgearPage(options = { belongsToSelf: true }) {
  //   const progearTypeSelector = options?.belongsToSelf
  //     ? `[name="belongsToSelf"][value="true"]`
  //     : `[name="belongsToSelf"][value="false"]`;

  //   await this.page.locator(progearTypeSelector).click({ force: true });
  //   await this.page.locator('[name="description"]').type('Radio equipment');

  //   if (options?.belongsToSelf) {
  //     await this.page.locator('[name="weight"]')
  //       .clear()
  //       .type(options?.weight || '2000');
  //   } else {
  //     await this.page.locator('[name="weight"]')
  //       .clear()
  //       .type(options?.weight || '500');
  //   }

  //   if (options?.missingWeightTicket) {
  //     await this.page.locator('[name="missingWeightTicket"]').click();
  //     cy.upload_file('.filepond--root', 'constructedWeight.xls');
  //   } else {
  //     cy.upload_file('.filepond--root', 'sampleWeightTicket.jpg');
  //   }

  //   cy.wait('@uploadFile');
  // }

  /**
   */
  async navigateFromCloseoutReviewPageToExpensesPage() {
    await Promise.all([this.page.waitForNavigation(), this.page.getByRole('button', { name: 'Add Expense' }).click()]);
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/expenses/`;
    expect(this.page.url()).toContain(url);
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
    await this.page.locator('input[name="paidWithGTCC"][value="true"]').click();
    await this.page.locator('input[name="amount"]').clear();
    await this.page.locator('input[name="amount"]').type(options?.amount || '675.99');

    // TODO: really should be using a better locator
    await this.page
      .locator('.receiptDocument.filepond--root')
      .setInputFiles('playwright/fixtures/sampleWeightTicket.jpg');

    await this.page.locator('input[name="sitStartDate"]').type('14 Aug 2022');
    await this.page.locator('input[name="sitStartDate"]').blur();
    await this.page.locator('input[name="sitEndDate"]').type('20 Aug 2022');
    await this.page.locator('input[name="sitEndDate"]').blur();

    await Promise.all([
      this.page.waitForNavigation(),
      this.page.getByRole('button', { name: 'Save & Continue' }).click(),
    ]);

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/review`;
    expect(this.page.url()).toContain(url);

    await expect(this.page.getByText('Cloud storage')).toBeVisible();
    await expect(this.page.locator('dt').getByText('Days in storage:')).toBeVisible();
    await expect(this.page.locator('dd').getByText('7')).toBeVisible();
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

    await expect(this.page.locator('li')).toContainText(`${options?.totalNetWeight} total net weight`);

    // TODO: Once we get moving expenses and pro gear back, check for those here as well.

    await expect(this.page.locator('li')).toContainText(`${options?.proGearWeight} of pro-gear`);
    await expect(this.page.locator('li')).toContainText(`$${options?.expensesClaimed} in expenses claimed`);
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

  // export function submitFinalCloseout(options) {
  //   verifyFinalIncentiveAndTotals(options);
  //   signCloseoutAgreement();
  // }
}
export { expect, test, setMobileViewport };

export default CustomerPpmPage;
