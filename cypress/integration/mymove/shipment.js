/* global cy */

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
    // Calendar move date
    cy
      .get('div[class="DayPicker-Day"][tabindex="0"]') // get's the 1st of the month
      .click()
      .should('have.class', 'DayPicker-Day--selected');

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

    // Weights
    cy
      .get('input[name="weight_estimate"]')
      .clear()
      .type('3000')
      .get('input[name="progear_weight_estimate"]')
      .clear()
      .type('250')
      .get('input[name="spouse_progear_weight_estimate"]')
      .clear()
      .type('158');

    // Review page
    cy.nextPage();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });

    // TODO: when shipment info is available on Review page, test edit of fields

    cy.nextPage();
    cy.contains('SIGNATURE');
    cy.get('input[name="signature"]').type('SM Signature');

    // Status summary page
    cy.nextPage();
    cy.contains('Success');
    cy.contains('Next Step: Awaiting approval');

    cy.resetDb();
  });
});
