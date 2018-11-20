/* global cy */

describe('service member adds a ppm to an hhg', function() {
  it('service member clicks on Add PPM Shipment', function() {
    serviceMemberSignsIn('f83bc69f-10aa-48b7-b9fe-425b393d49b8');
    serviceMemberAddsPPMToHHG();
    serviceMemberCancelsAddPPMToHHG();
    serviceMemberAddsPPMToHHG();
    serviveMemberFillsInDatesAndLocations();
    serviceMemberSelectsWeightRange();
    serviceMemberCanCustomizeWeight();
    serviceMemberCanReviewMoveSummary();
    serviceMemberCanSignAgreement();
    serviceMemberViewsUpdatedHomePage();
  });
});

function serviceMemberSignsIn(uuid) {
  cy.signInAsUser(uuid);
}

function serviceMemberAddsPPMToHHG() {
  cy
    .get('.sidebar > div > a')
    .contains('Add PPM Shipment')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-start/);
  });

  // does not have a back button on first flow page
  cy
    .get('button')
    .contains('Back')
    .should('not.be.visible');
}

function serviceMemberCancelsAddPPMToHHG() {
  cy
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\//);
  });
}

function serviveMemberFillsInDatesAndLocations() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-start/);
  });

  cy
    .get('input[name="planned_move_date"]')
    .first()
    .type('9/2/2018{enter}')
    .blur();

  cy
    .get('input[name="pickup_postal_code"]')
    .clear()
    .type('80913');

  cy
    .get('input[name="destination_postal_code"]')
    .clear()
    .type('76127');

  cy.nextPage();
}

function serviceMemberSelectsWeightRange() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-size/);
  });

  //todo verify entitlement
  cy.contains('A trailer').click();

  cy.nextPage();
}

function serviceMemberCanCustomizeWeight() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-weight/);
  });

  cy.get('.rangeslider__handle').click();

  cy.get('.incentive').contains('$');

  cy.nextPage();
}

function serviceMemberCanReviewMoveSummary() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });

  cy.get('body').should($div => expect($div.text()).not.to.include('Government moves all of your stuff (HHG)'));
  cy.get('.ppm-container').should($div => {
    const text = $div.text();
    expect(text).to.include('Shipment - You move your stuff (PPM)');
    expect(text).to.include('Move Date: 09/02/2018');
    expect(text).to.include('Pickup ZIP Code:  80913');
    expect(text).to.include('Delivery ZIP Code:  76127');
    expect(text).to.include('Storage: Not requested');
    expect(text).to.include('Estimated Weight:  1,50');
    expect(text).to.include('Estimated PPM Incentive:  $2,032.89 - 2,246.87');
  });

  cy.nextPage();
}

function serviceMemberCanSignAgreement() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
  });

  cy
    .get('body')
    .should($div =>
      expect($div.text()).to.include(
        'Before officially booking your move, please carefully read and then sign the following.',
      ),
    );

  cy.get('input[name="signature"]').type('Jane Doe');
  cy.nextPage();
}

function serviceMemberViewsUpdatedHomePage() {
  cy.location().should(loc => {
    expect(loc.pathname).to.eq('/');
  });
}
