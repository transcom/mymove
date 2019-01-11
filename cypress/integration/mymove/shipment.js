/* global cy, Cypress */

describe('completing the hhg flow', function() {
  beforeEach(() => {
    // sm_hhg@example.com
    cy.signInAsUser('4b389406-9258-4695-a091-0bf97b5a132f');
  });

  it('selects hhg and progresses thru form', function() {
    cy.contains('Continue Move Setup').click();
    cy
      .contains('Household Goods Move')
      .click()
      .nextPage();

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

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-locations/);
    });
    cy.get('button.next').should('be.disabled');

    // Pickup address
    cy
      .get('input[name="pickup_address.street_address_1"]')
      .clear({ force: true })
      .type('123 Elm Street');
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
    cy.get('input[name="delivery_address.postal_code"]').type('94567');

    cy.nextPage();

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

    // Review page
    cy.nextPage();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });
    cy.contains('Government moves all of your stuff (HHG)');

    cy.contains('123 Elm Street'); // pickup address
    cy.contains('543 Oak Street'); // secondary pickup address
    cy.contains('678 Madrone Street'); // destination address

    cy.contains('3,000 lbs + 250 lbs pro-gear + 158 lbs spouse pro-gear');
    cy.contains('Great! You appear within your weight allowance.');

    // Go back and edit dates
    cy
      .get('.hhg-shipment-summary .hhg-dates a')
      .contains('Edit')
      .click();

    cy.contains('Pick a moving date');
    cy.get('.DayPicker-Body').then($body => {
      if ($body.find('.DayPicker-Day--transit[aria-disabled=false]').length === 0) {
        cy.get('.DayPicker-NavButton--next').click();
      }
    });
    cy
      .get('.DayPicker-Day--transit[aria-disabled=false]') // pick the first transit day that is selectable
      .first()
      .click()
      .should('have.class', 'DayPicker-Day--pickup')
      .invoke('attr', 'aria-label')
      .as('dateLabel');

    cy.contains('Save').click();

    // verify new date is shown
    checkLoadingDate();

    // click edit again and then click cancel
    cy
      .get('.hhg-shipment-summary .hhg-dates a')
      .contains('Edit')
      .click();
    cy.contains('Cancel').click();

    // verify loading date hasn't changed
    checkLoadingDate();

    cy.nextPage();
    cy.contains('SIGNATURE');
    cy.get('input[name="signature"]').type('SM Signature');

    // Status summary page
    cy.nextPage();
    cy.contains('Success');
    cy.contains('Government Movers and Packers');
  });
});

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
