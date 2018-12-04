/* global cy */

describe('service member adds a ppm to an hhg', function() {
  it('service member clicks on Add PPM Shipment', function() {
    serviceMemberSignsIn('f83bc69f-10aa-48b7-b9fe-425b393d49b8');
    serviceMemberAddsPPMToHHG();
    serviceMemberCancelsAddPPMToHHG();
    serviceMemberContinuesPPMSetup();
    serviceMemberFillsInDatesAndLocations();
    serviceMemberSelectsWeightRange();
    serviceMemberCanCustomizeWeight();
    serviceMemberCanReviewMoveSummary();
    serviceMemberCanSignAgreement();
    serviceMemberViewsUpdatedHomePage();
  });
  it('service member edits an HHG_PPM Shipment', function() {
    serviceMemberSignsIn('f83bc69f-10aa-48b7-b9fe-425b393d49b8');
    serviceMemberClicksEditMove();
    serviceMemberVerifiesHHGPPMSummary();
  });
});

function serviceMemberVerifiesHHGPPMSummary() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/edit/);
  });

  cy.get('.review-content').should($div => {
    const text = $div.text();
    // Profile section
    expect(text).to.include('Name: HHG Ready  For PPM');
    expect(text).to.include('Branch:Army');
    expect(text).to.include('Rank/Pay Grade: E-1');
    expect(text).to.include('DoD ID#: 7777567890');
    expect(text).to.include('Current Duty Station: Yuma AFB');

    // Orders section
    expect(text).to.include('Orders Type: Permanent Change Of Station');
    expect(text).to.include('Orders Date:  05/20/2018');
    expect(text).to.include('Report-by Date: 08/01/2018');
    expect(text).to.include('New Duty Station:  Yuma AFB');
    expect(text).to.include('Dependents?:  Yes');
    expect(text).to.include('Spouse Pro Gear?: Yes');
    expect(text).to.include('Orders Uploaded: 1');

    // Contact Info
    expect(text).to.include('Best Contact Phone: 555-555-5555');
    expect(text).to.include('Alt. Phone:');
    expect(text).to.include('Personal Email: hhgforppm@award.ed');
    expect(text).to.include('Preferred Contact Method: Email');
    expect(text).to.include('Current Mailing Address: 123 Any Street');
    expect(text).to.include('P.O. Box 12345');
    expect(text).to.include('Beverly Hills, CA 90210');
    expect(text).to.include('Backup Mailing Address: 123 Any Street');
    expect(text).to.include('P.O. Box 12345');
    expect(text).to.include('Beverly Hills, CA 90210');

    // Backup Contact info
    expect(text).to.include('Backup Contact: name');
    expect(text).to.include('Email:  email@example.com');
    expect(text).to.include('Phone:  555-555-5555');
  });

  cy.get('body').should($div => expect($div.text()).not.to.include('Government moves all of your stuff (HHG)'));
  cy.get('.ppm-container').should($div => {
    const text = $div.text();
    expect(text).to.include('Shipment - You move your stuff (PPM)');
    expect(text).to.include('Move Date: 05/20/2018');
    expect(text).to.include('Pickup ZIP Code:  90210');
    expect(text).to.include('Delivery ZIP Code:  50309');
    expect(text).not.to.include('Storage: Not requested');
    expect(text).to.include('Estimated Weight:  1,50');
    expect(text).to.include('Estimated PPM Incentive:  $4,255.80 - 4,703.78');
  });
}

function serviceMemberClicksEditMove() {
  cy
    .get('.usa-button-secondary')
    .contains('Edit Move')
    .click();
}

function serviceMemberSignsIn(uuid) {
  cy.signInAsUser(uuid);
}

function serviceMemberAddsPPMToHHG() {
  cy
    .get('.sidebar > div > button')
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

function serviceMemberContinuesPPMSetup() {
  cy
    .get('button')
    .contains('Continue Move Setup')
    .click();
}

function serviceMemberFillsInDatesAndLocations() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-start/);
  });

  cy
    .get('input[name="planned_move_date"]')
    .should('have.value', '5/20/2018')
    .clear()
    .first()
    .type('9/2/2018{enter}')
    .blur();

  cy.get('input[name="pickup_postal_code"]').should('have.value', '90210');

  cy.get('input[name="destination_postal_code"]').should('have.value', '50309');

  cy.nextPage();
}

function serviceMemberSelectsWeightRange() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-size/);
  });

  cy.get('.entitlement-container p:nth-child(2)').should($div => {
    const text = $div.text();
    expect(text).to.include('Estimated 2,000 lbs entitlement remaining (10,500 lbs - 8,500 lbs estimated HHG weight).');
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
    expect(text).to.include('Move Date: 05/20/2018');
    expect(text).to.include('Pickup ZIP Code:  90210');
    expect(text).to.include('Delivery ZIP Code:  50309');
    expect(text).not.to.include('Storage: Not requested');
    expect(text).to.include('Estimated Weight:  1,50');
    expect(text).to.include('Estimated PPM Incentive:  $4,255.80 - 4,703.78');
  });

  cy.nextPage();
}
function serviceMemberCanSignAgreement() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-agreement/);
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

  cy.get('body').should($div => {
    expect($div.text()).to.include('Government Movers and Packers');
    expect($div.text()).to.include('Move your own stuff');
    expect($div.text()).to.not.include('Add PPM Shipment');
  });

  cy.get('.usa-width-three-fourths').should($div => {
    const text = $div.text();
    // PPM information and details
    expect(text).to.include('Next Step: Wait for approval');
    expect(text).to.include('Weight (est.): 150');
    expect(text).to.include('Incentive (est.): $4,255.80 - 4,703.78');
  });
}
