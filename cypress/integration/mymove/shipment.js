/* global cy, Cypress */

describe('shipments', function() {
  describe('completing the hhg flow', function() {
    it('selects hhg and progresses thru form', function() {
      // sm_hhg@example.com
      cy.signInAsUser('4b389406-9258-4695-a091-0bf97b5a132f');
      serviceMemberAddsHHG();
      serviceMemberAddsMoveDates();
      serviceMemberAddsLocations();
      serviceMemberAddsWeight();
      serviceMemberAddsProGear();
      serviceMemberEditsDates();
      serviceMemberCancelsDateEdit();
      serviceMemberReviewsMove();
      serivceMemberSigns();
      moveIsSuccessfullyCreated();
    });
  });
  describe('continuing a hhg move after logging out before completion', function() {
    it('resumes where service member left off', function() {
      // sm_hhg_continue@example.com
      const serviceMemberId = '1256a6ea-27cc-4d60-92df-1bc2a5c39028';
      const firstIncompletePage = /^\/moves\/[^/]+\/hhg-weight/;

      cy.removeFetch();
      cy.signInAsUser(serviceMemberId);
      serviceMemberAddsHHG();
      serviceMemberAddsMoveDates();
      serviceMemberAddsLocations();
      cy.location().should(loc => {
        expect(loc.pathname).to.match(firstIncompletePage);
      });
      serviceMemberLogsOutThenContinues(serviceMemberId);

      cy.location().should(loc => {
        expect(loc.pathname).to.match(firstIncompletePage);
      });
    });
  });
});

function serviceMemberLogsOutThenContinues(serviceMemberId) {
  cy.logout();
  // wait returns after the 1st call to getShipments, so if multiple calls
  // placement of server and route are important
  cy.server();
  cy.route({ url: '**/api/v1/shipments/*' }).as('getShipments');
  cy.signInAsUser(serviceMemberId);
  cy.wait('@getShipments');
  cy.contains('Continue Move Setup').click();
}

function serviceMemberAddsHHG() {
  cy.contains('Continue Move Setup').click();
  cy
    .contains('Household Goods Move')
    .click()
    .nextPage();
}

function serviceMemberAddsMoveDates() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-start/);
  });
  cy.get('button.next').should('be.disabled');

  // Calendar move date

  // Try to get today, which should be disabled.  We may have to go back a month to find
  // today since the calendar scrolls to the month with the first available move date.
  cy
    .get('.DayPicker-Body')
    .then($body => {
      if ($body.find('.DayPicker-Day--today').length === 0) {
        cy.get('.DayPicker-NavButton--prev').click();
      }
    })
    .then(() => {
      cy
        .get('.DayPicker-Day--today')
        .first()
        .should('have.class', 'DayPicker-Day--disabled');
    });

  // We may or may not have an available date in the current month.  If not, then
  // we skip to the next month which should (at least at this point) have an available
  // date.
  cy
    .get('.DayPicker-Body')
    .then($body => {
      if ($body.find('[class=DayPicker-Day]').length === 0) {
        cy.get('.DayPicker-NavButton--next').click();
      }
    })
    .then(() => {
      cy
        .get('[class=DayPicker-Day]')
        .first()
        .click()
        .should('have.class', 'DayPicker-Day--pickup');
    });

  // Check for calendar move dates summary and color-coding of calendar.
  cy.contains('Movers Packing');
  cy.get('.DayPicker-Day.DayPicker-Day--pickup');
  cy.nextPage();
}

function serviceMemberAddsLocations() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-locations/);
  });
  // Note that we are not checking for a disabled save button because we
  // expect the pickup address to prefill with the SM residential address

  // Pickup address
  cy.get('input[name="pickup_address.street_address_1"]').should('have.value', '123 Any Street');
  cy.get('input[name="pickup_address.street_address_2"]').should('have.value', 'P.O. Box 12345');
  cy.get('input[name="pickup_address.city"]').should('have.value', 'Beverly Hills');
  cy.get('select[name="pickup_address.state"]').should('have.value', 'CA');
  cy.get('input[name="pickup_address.postal_code"]').should('have.value', '90210');

  // Pickup address
  cy
    .get('input[name="pickup_address.street_address_1"]')
    .clear({ force: true })
    .type('123 Elm Street');
  cy.get('input[name="pickup_address.street_address_2"]').clear({ force: true });
  cy
    .get('input[name="pickup_address.city"]')
    .clear()
    .type('Alameda');
  cy.get('select[name="pickup_address.state"]').select('CA');
  cy.get('input[name="pickup_address.postal_code"]').type('94607');

  // Check yes for radios
  cy.get('input[type="radio"]').check('yes', { force: true }); // checks yes for both radios on form

  // Secondary pickup address
  cy
    .get('input[name="secondary_pickup_address.street_address_1"]')
    .clear({ force: true })
    .type('543 Oak Street');
  cy
    .get('input[name="secondary_pickup_address.city"]')
    .clear()
    .type('Oakland');
  cy.get('select[name="secondary_pickup_address.state"]').select('CA');
  cy.get('input[name="secondary_pickup_address.postal_code"]').type('94609');

  // Destination address
  cy
    .get('input[name="delivery_address.street_address_1"]')
    .clear({ force: true })
    .type('678 Madrone Street');
  cy
    .get('input[name="delivery_address.city"]')
    .clear()
    .type('Fremont');
  cy.get('select[name="delivery_address.state"]').select('CA');
  cy
    .get('input[name="delivery_address.postal_code"]')
    .clear()
    .type('94567');

  cy.nextPage();
}

function serviceMemberAddsWeight() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-weight/);
  });

  // Weight calculator
  cy
    .get('input[name="rooms"]')
    .clear()
    .type('9');

  cy.get('select[name="stuff"]').select('more');

  cy.get('input[name="weight_estimate"]').should('have.value', '13500');

  // Weight over entitlement
  cy
    .get('input[name="weight_estimate"]')
    .clear()
    .type('50000');

  cy.contains('Entitlement exceeded');

  // Weight
  cy
    .get('input[name="weight_estimate"]')
    .clear()
    .type('3000');

  cy.contains('Entitlement exceeded').should('not.exist');

  cy.nextPage();
}

function serviceMemberAddsProGear() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-progear/);
  });

  // Progear Weights

  // Check yes for radios
  cy.get('input[type="radio"]').check('yes', { force: true }); // checks yes for both radios on form

  cy
    .get('input[name="progear_weight_estimate"]')
    .clear()
    .type('2500');

  cy.contains('Entitlement exceeded');

  cy
    .get('input[name="progear_weight_estimate"]')
    .clear()
    .type('250');

  cy.contains('Entitlement exceeded').should('not.exist');

  cy
    .get('input[name="spouse_progear_weight_estimate"]')
    .clear()
    .type('1580');

  cy.contains('Entitlement exceeded');

  cy
    .get('input[name="spouse_progear_weight_estimate"]')
    .clear()
    .type('158');

  cy.contains('Entitlement exceeded').should('not.exist');
  cy.nextPage();
}

function serviceMemberReviewsMove() {
  // Review page
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
  cy.contains('Government moves all of your stuff (HHG)');

  cy.contains('123 Elm Street'); // pickup address
  cy.contains('543 Oak Street'); // secondary pickup address
  cy.contains('678 Madrone Street'); // destination address

  cy.contains('3,000 lbs + 250 lbs pro-gear + 158 lbs spouse pro-gear');
  cy.contains('Great! You appear within your weight allowance.');
  cy.nextPage();
}

function serviceMemberEditsDates() {
  // Go back and edit dates
  cy
    .get('.hhg-shipment-summary .hhg-dates a')
    .contains('Edit')
    .click();

  // Wait for valid move dates to be loaded from the server
  cy.waitForLoadingScreen();

  cy.contains('Pick a moving date');
  cy
    .get('.DayPicker-Body')
    .then($body => {
      if ($body.find('.DayPicker-Day--transit[aria-disabled=false]').length === 0) {
        cy.get('.DayPicker-NavButton--next').click();
      }
    })
    .then(() => {
      cy
        .get('.DayPicker-Day--transit[aria-disabled=false]') // pick the first transit day that is selectable
        .first()
        .click()
        .should('have.class', 'DayPicker-Day--pickup')
        .invoke('attr', 'aria-label')
        .as('dateLabel');

      cy.contains('Save').click();
      checkLoadingDate();
    });
}

function serviceMemberCancelsDateEdit() {
  // click edit again and then click cancel
  cy
    .get('.hhg-shipment-summary .hhg-dates a')
    .contains('Edit')
    .click();
  cy.contains('Cancel').click();
  checkLoadingDate();
}

function serivceMemberSigns() {
  cy.contains('SIGNATURE');
  cy.get('input[name="signature"]').type('SM Signature');

  // Status summary page
  cy.nextPage();
}

function moveIsSuccessfullyCreated() {
  cy.contains('Success');
  cy.contains('Government Movers and Packers');
}

// Verify that the date shown on the review page for loading is the same date
// that dateLabel was most recently set to.
function checkLoadingDate() {
  cy
    .get('.hhg-shipment-summary .hhg-dates tr')
    .contains('Loading Truck:')
    .next()
    .then(function($tr) {
      let actualDate = Cypress.moment($tr.text(), 'ddd, MMM DD');
      const expectedDate = Cypress.moment(this.dateLabel, 'ddd MMM DD YYYY');
      // since the actual date doesn't have year information, moment may fail to parse if it is next year
      if (!actualDate.isValid()) {
        actualDate = Cypress.moment($tr.text() + ' ' + expectedDate.year(), 'ddd, MMM DD YYYY');
      }
      expect(actualDate.toString()).to.equal(expectedDate.toString());
    });
}
