import {
  fillOutAboutPage,
  navigateFromCloseoutReviewPageToAboutPage,
  navigateFromCloseoutReviewPageToEditExpensePage,
  navigateFromCloseoutReviewPageToEditProGearPage,
  navigateFromCloseoutReviewPageToEditWeightTicketPage,
  navigateFromCloseoutReviewPageToExpensesPage,
  navigateFromCloseoutReviewPageToProGearPage,
  navigateFromPPMReviewPageToFinalCloseoutPage,
  setMobileViewport,
  signInAndNavigateToAboutPage,
  signInAndNavigateToPPMReviewPage,
  submitExpensePage,
  submitFinalCloseout,
  submitProgearPage,
  submitWeightTicketPage,
} from '../../../support/ppmCustomerShared';

describe('Entire PPM closeout flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('POST', '**/internal/mto_shipments').as('createShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('POST', '**/internal/ppm-shipments/**/weight-ticket').as('createWeightTicket');
    cy.intercept('POST', '**/internal/ppm-shipments/**/uploads**').as('uploadFile');
    cy.intercept('GET', '**/internal/moves/**/signed_certifications').as('signedCertifications');
    cy.intercept('POST', '**/internal/ppm-shipments/**/submit-ppm-shipment-documentation').as('submitCloseout');
  });

  const viewportType1 = [
    { viewport: 'desktop', isMobile: false, userId: '1dca189a-ca7e-4e70-b98e-be3829e4b6cc' }, // readyForCloseout@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: 'fe825617-a53a-49bf-bf2e-c271afee344d' }, // readyForCloseout2@ppm.approved
  ];

  viewportType1.forEach(({ viewport, isMobile, userId }) => {
    it(`flows through happy path for existing shipment - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      navigateHappyPath(userId, isMobile);
    });
  });

  const viewportType2 = [
    { viewport: 'desktop', isMobile: false, userId: '6f48be45-8ee0-4792-a961-ec6856e5435d' }, // closeoutHappyPathWithEdits@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '917da44e-7e44-41be-b912-1486a72b69d8' }, // closeoutHappyPathWithEditsMobile@ppm.approved
  ];

  viewportType2.forEach(({ viewport, isMobile, userId }) => {
    it(`happy path with edits and backs - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      navigateHappyPathWithEditsAndBacks(userId, isMobile);
    });
  });
});

function navigateHappyPath(userId, isMobile = false) {
  signInAndNavigateToAboutPage(userId);
  submitWeightTicketPage();
  navigateFromCloseoutReviewPageToProGearPage();
  submitProgearPage();
  navigateFromCloseoutReviewPageToExpensesPage();
  submitExpensePage();
  navigateFromPPMReviewPageToFinalCloseoutPage();
  submitFinalCloseout({
    totalNetWeight: '2,000 lbs',
    proGearWeight: '2,000 lbs',
    expensesClaimed: '675.99',
    finalIncentiveAmount: '$31,180.87',
  });
}

function navigateHappyPathWithEditsAndBacks(userId, isMobile = false) {
  signInAndNavigateToPPMReviewPage(userId);
  navigateFromCloseoutReviewPageToAboutPage();
  fillOutAboutPage();
  submitWeightTicketPage(); // temporary until about page routing is fixed, 2 weight tickets will exist for now
  navigateFromCloseoutReviewPageToEditWeightTicketPage();
  submitWeightTicketPage();
  navigateFromCloseoutReviewPageToEditProGearPage();
  submitProgearPage({ belongsToSelf: false });
  navigateFromCloseoutReviewPageToEditExpensePage();
  submitExpensePage({ isEditExpense: true, amount: '833.41' });
  navigateFromPPMReviewPageToFinalCloseoutPage();
  submitFinalCloseout({
    totalNetWeight: '4,000 lbs',
    proGearWeight: '500 lbs',
    expensesClaimed: '833.41',
    finalIncentiveAmount: '$62,363.15',
  });
}
