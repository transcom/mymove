/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, CustomerPage } = require('../../../utils/customerTest');

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
    this.firstShipmentId = this.move.MTOShipments[0].ID;
  }

  /**
   * sign in for ppm
   */
  async signInAndClickOnUploadPPMDocumentsButton() {
    await this.signInAsExistingCustomer(this.userId);
    await expect(this.page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();

    await this.page.getByRole('button', { name: 'Upload PPM Documents' }).click();
  }

  /**
   */
  async customerStartsAddingAPPMShipment() {
    await this.page.locator('button[data-testid="shipment-selection-btn"]').click();
    await this.navigateForward();

    await this.page.locator('input[value="PPM"]').check();
    await this.navigateForward();
  }

  /**
   * @param {string} userId
   * @param {boolean} isMoveSubmitted
   */
  // async signInAndNavigateFromHomePageToReviewPage(userId, isMoveSubmitted = false) {
  //   await this.signInAsExistingCustomer(userId);

  //   await this.navigateFromHomePageToReviewPage(isMoveSubmitted);
  // }

  /**
   * @param {boolean} selectAdvance
   */
  async signInAndNavigateToAboutPage(selectAdvance) {
    await this.signInAndClickOnUploadPPMDocumentsButton();

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/about`;

    expect(this.page.url()).toContain(url);

    await expect(this.page.getByRole('heading', { name: 'About your PPM' })).toBeVisible();

    await this.fillOutAboutPage(selectAdvance);
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
   * @param {boolean} isMoveSubmitted
   */
  async navigateFromHomePageToReviewPage(isMoveSubmitted = false) {
    if (isMoveSubmitted) {
      await expect(this.page.getByRole('heading', { name: 'Next step: Your move gets approved' })).toBeVisible();

      await this.page.getByRole('button', { name: 'Review your request' }).click();
    } else {
      await expect(this.page.getByRole('heading', { name: 'Time to submit your move' })).toBeVisible();

      await this.page.getByRole('button', { name: 'Review and submit' }).click();
    }
  }

  /**
   * @param {boolean} selectAdvance
   */
  async fillOutAboutPage(selectAdvance) {
    // editing this field with the keyboard instead of the date picker runs async validators for pre-filled postal codes
    // this helps debounce the API calls that would be triggered in quick succession
    await this.page.locator('input[name="actualMoveDate"]').type('01 Feb 2022');
    await this.page.locator('input[name="actualMoveDate"]').blur();

    await this.page.locator('input[name="actualPickupPostalCode"]').type('90210');
    await this.page.locator('input[name="actualPickupPostalCode"]').blur();
    await this.page.locator('input[name="actualDestinationPostalCode"]').type('76127');
    await this.page.locator('input[name="actualDestinationPostalCode"]').blur();

    if (selectAdvance) {
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
  async submitWeightTicketPage(options) {
    await this.fillOutWeightTicketPage(options);
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
    await this.page.locator('button').getByText('Save & Continue').click();

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/review`;
    expect(this.page.url()).toContain(url);
  }

  /**
   */
  async signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage() {
    await this.signInAsExistingCustomer(this.userId);

    await expect(this.page.getByRole('heading', { name: 'Time to submit your move' })).toBeVisible();

    await this.page.locator('[data-testid="shipment-list-item-container"] button').getByText('Edit').click();

    await expect(this.page.getByRole('heading', { name: 'PPM date & location' })).toBeVisible();
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/edit`;
    expect(this.page.url()).toContain(url);
  }

  // // used for creating a new shipment
  // submitsDateAndLocation() {
  //   await this.page.locator('input[name="pickupPostalCode"]').type('90210'); await this.page.locator('input[name="pickupPostalCode"]').blur();

  //   await this.page.locator('input[name="destinationPostalCode"]').type('76127');
  //   await this.page.locator('input[name="expectedDepartureDate"]').type('01 Feb 2022'); await this.page.locator('input[name="expectedDepartureDate"]').blur();

  //   // Select closeout office
  //   this.selectDutyLocation('Fort Bragg', 'something');

  //   this.navigateFromDateAndLocationPageToEstimatedWeightsPage();
  // }

  /**
   */
  async navigateFromDateAndLocationPageToEstimatedWeightsPage() {
    await this.page.locator('button').getByText('Save & Continue').click();

    await expect(this.page.getByRole('heading', { name: 'Estimated weight', exact: true })).toBeVisible();

    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/estimated-weight`;
    expect(this.page.url()).toContain(url);
  }

  //  submitsEstimatedWeightsAndProGear() {
  //   await this.page.locator('input[name="estimatedWeight"]').type(4000); await this.page.locator('input[name="estimatedWeight"]').blur();
  //   await this.page.locator('input[name="hasProGear"][value="true"]').check({ force: true });
  //   await this.page.locator('input[name="proGearWeight"]').type(500); await this.page.locator('input[name="proGearWeight"]').blur();
  //   await this.page.locator('input[name="spouseProGearWeight"]').type(400); await this.page.locator('input[name="spouseProGearWeight"]').blur();
  //   await expect(this.page.locator('button').contains('Save & Continue')).toBeEnabled();

  //   navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  // }

  // export function submitsEstimatedWeights() {
  //   await this.page.locator('input[name="estimatedWeight"]').type(4000); await this.page.locator('input[name="estimatedWeight"]').blur();
  //   await expect(this.page.locator('button').contains('Save & Continue')).toBeEnabled();

  //   navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  // }

  async navigateFromEstimatedWeightsPageToEstimatedIncentivePage() {
    await this.page.locator('button').getByText('Save & Continue').click();

    await expect(this.page.getByRole('heading', { name: 'Estimated incentive', exact: true })).toBeVisible();
    const url = `/moves/${this.move.id}/shipments/${this.firstShipmentId}/estimated-incentive`;
    expect(this.page.url()).toContain(url);
  }

  // export function generalVerifyEstimatedIncentivePage(isMobile = false) {
  //   await this.page.locator('h1').should('contain', 'Estimated incentive');

  //   // checks the format of the incentive amount statement is `$<some comma-separated number without decimals> is`
  //   await expect(this.page.locator('.container h2')).toContainText(/\$\d{1,3}(?:,\d{3})*? is/);

  //   if (!isMobile) {
  //     await expect(this.page.locator('button')).toContainText('Next').should('not.be.disabled');
  //   } else {
  //     await expect(this.page.locator('button')).toContainText('Next').should('not.be.disabled').should('have.css', 'order', '1');
  //   }

  //   navigateFromEstimatedIncentivePageToAdvancesPage();
  // }

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
      await this.page.locator('label[for="yes-has-received-advance"]').click();
      await this.page.locator('input[name="advanceAmountRequested"]').type('4000');
      await this.page.locator('input[name="advanceAmountRequested"]').blur();

      await this.page.locator('input[name="agreeToTerms"]').check();
    } else {
      await this.page.locator('label[for="no-has-received-advance"]').click();
    }

    await this.navigateFromAdvancesPageToReviewPage({ addAdvance, isMobile });
  }

  /**
   * navigate from advances to review
   * @param {Object} options
   * @param {boolean} [options.addAdvance=false]
   * @param {boolean} [options.isMobile=false]
   */
  async navigateFromAdvancesPageToReviewPage(options = {}) {
    const { addAdvance = false, isMobile = false } = options;
    if (addAdvance) {
      await this.page.locator('input[name="agreeToTerms"]').check();
    }

    const saveButton = this.page.getByRole('button', { name: 'Save & Continue' });
    await expect(saveButton).toBeVisible();

    if (isMobile) {
      await expect(saveButton).toHaveCSS('order', '1');
    }

    await saveButton.click();

    await expect(this.page.getByRole('heading', { name: 'Review your details', exact: true })).toBeVisible();
    const url = `/moves/${this.move.id}/review/`;
    expect(this.page.url()).toContain(url);

    await expect(this.page.locator('.usa-alert__heading')).toContainText('Details saved');
    await expect(this.page.locator('.usa-alert__heading + span')).toContainText(
      'Review your info and submit your move request now, or come back and finish later.',
    );
  }

  // export function navigateFromReviewPageToHomePage() {
  //   await expect(this.page.locator('button')).toContainText('Return home').click();

  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.equal('/');
  //   });
  // }

  // export function navigateToFromHomePageToPPMCloseoutReview() {
  //   await expect(this.page.locator('[data-testid="stepContainer5"] button').contains('Upload PPM Documents')).toBeEnabled().click();

  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  //   });
  // }

  // export function navigateToAgreementAndSign() {
  //   cy.nextPage();
  //   signAgreement();
  // }

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

  // export function signInAndNavigateToProgearPage(userId) {
  //   signInAndNavigateToPPMReviewPage(userId);
  //   navigateFromCloseoutReviewPageToProGearPage();
  // }

  // export function navigateFromCloseoutReviewPageToProGearPage() {
  //   await expect(this.page.locator('a.usa-button')).toContainText('Add Pro-gear Weight').click();
  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  //   });
  // }

  // export function navigateFromCloseoutReviewPageToEditProGearPage() {
  //   await expect(this.page.locator('.progearSection a').eq(1)).toContainText('Edit').click();
  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/pro-gear/);
  //   });
  // }

  // export function navigateFromCloseoutReviewPageToEditWeightTicketPage() {
  //   await expect(this.page.locator('.reviewWeightTickets a').eq(1)).toContainText('Edit').click();
  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  //   });
  // }

  // export function navigateFromCloseoutReviewPageToEditExpensePage() {
  //   await expect(this.page.locator('.reviewExpenses a').eq(1)).toContainText('Edit').click();
  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/expenses/);
  //   });
  // }

  // export function navigateFromCloseoutReviewPageToAboutPage() {
  //   await expect(this.page.locator('[data-testid="aboutYourPPM"] a')).toContainText('Edit').click();
  // }

  // export function navigateFromProgearPage() {
  //   await expect(this.page.locator('button').contains('Save & Continue')).toBeEnabled().click();

  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  //   });
  // }

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

  // export function navigateFromCloseoutReviewPageToExpensesPage() {
  //   await expect(this.page.locator('a.usa-button')).toContainText('Add Expense').click();
  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/expenses/);
  //   });
  // }

  // export function submitExpensePage(options = { isEditExpense: false }) {
  //   await this.page.locator('select[name="expenseType"]').as('expenseType');
  //   if (!options?.isEditExpense) {
  //     await this.page.locator('@expenseType').should('have.value', '');
  //   }
  //   await this.page.locator('@expenseType').selectOption({ label: 'Storage'});

  //   await this.page.locator('input[name="description"]').type('Cloud storage');
  //   await this.page.locator('input[name="paidWithGTCC"][value="true"]').click({ force: true });
  //   await this.page.locator('input[name="amount"]')
  //     .clear()
  //     .type(options?.amount || '675.99');

  //   cy.upload_file('.receiptDocument.filepond--root', 'sampleWeightTicket.jpg');
  //   cy.wait('@uploadFile');

  //   await this.page.locator('input[name="sitStartDate"]').type('14 Aug 2022'); await this.page.locator('input[name="sitStartDate"]').blur();
  //   await this.page.locator('input[name="sitEndDate"]').type('20 Aug 2022'); await this.page.locator('input[name="sitEndDate"]').blur();

  //   await expect(this.page.locator('button').contains('Save & Continue')).toBeEnabled().click();
  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/review/);
  //   });

  //   await expect(page.getByText('Cloud storage')).toBeVisible();
  //   await expect(page.getByText('dt', 'Days in storage:')).toBeVisible();
  //   await expect(page.getByText('dd', '7')).toBeVisible();
  // }

  // export function verifyFinalIncentiveAndTotals(
  //   options = {
  //     totalNetWeight: '4,000 lbs',
  //     proGearWeight: '1,500 lbs',
  //     expensesClaimed: '450.00',
  //     finalIncentiveAmount: '$500,000.00',
  //   },
  // ) {
  //   await expect(this.page.locator('h2')).toContainText(`Your final estimated incentive: ${options?.finalIncentiveAmount}`);

  //   await expect(this.page.locator('li')).toContainText(`${options?.totalNetWeight} total net weight`);

  //   // TODO: Once we get moving expenses and pro gear back, check for those here as well.

  //   await expect(this.page.locator('li')).toContainText(`${options?.proGearWeight} of pro-gear`);
  //   await expect(this.page.locator('li')).toContainText(`$${options?.expensesClaimed} in expenses claimed`);
  // }

  // export function signCloseoutAgreement() {
  //   await this.page.locator('input[name="signature"]').type('Sofía Clark-Nuñez');
  //   await expect(this.page.locator('button').contains('Submit PPM Documentation')).toBeEnabled().click();
  //   cy.wait('@submitCloseout');

  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.eq('/');
  //   });

  //   await expect(this.page.locator('.usa-alert--success')).toContainText('You submitted documentation for review.');

  //   await this.page.locator('[data-testid="stepContainer5"]').within(() => {
  //     await expect(this.page.locator('button').contains('Download Incentive Packet')).toBeDisabled();
  //     await expect(page.getByText(/PPM documentation submitted: \d{2} \w{3} \d{4}/)).toBeVisible();
  //   });
  // }

  // export function submitFinalCloseout(options) {
  //   verifyFinalIncentiveAndTotals(options);
  //   signCloseoutAgreement();
  // }
}
export { expect, test };

export default CustomerPpmPage;
