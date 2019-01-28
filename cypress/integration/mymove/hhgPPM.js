/* global cy */

describe('service member adds a ppm to an hhg', function() {
  it('service member clicks on Add PPM (DITY) Move', function() {
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
    serviceMemberEditsProfileSection();
    serviceMemberVerifiesProfileWasEdited();
    serviceMemberEditsOrdersSection();
    serviceMemberVerifiesOrderWasEdited();
    serviceMemberEditsContactInfoSection();
    serviceMemberVerifiesContactInfoWasEdited();
    serviceMemberEditsBackupContactInfoSection();
    serviceMemberVerifiesBackupContactInfoWasEdited();
    serviceMemberEditsHHGMoveDates();
    serviceMemberEditsPPMDatesAndLocations();
    serviceMemberVerifiesPPMDatesAndLocationsEdited();
    serviceMemberEditsPPMWeight();
    serviceMemberVerifiesPPMWeightsEdited();
    serviceMemberGoesBackToHomepage();
  });
});

function serviceMemberGoesBackToHomepage() {
  cy
    .get('.back-to-home')
    .contains('BACK TO HOME')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.eq('/');
  });
}

function serviceMemberEditsHHGMoveDates() {
  cy.get('[data-cy="edit-move"]').click();
  // TODO: Currently does not support changing move dates for 2019. Add test to edit dates ewhen fixed

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/edit/);
  });
}

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
    expect(text).to.include('New Duty Station:  Fort Gordon');
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

  cy.get('.ppm-container').should($div => {
    const text = $div.text();

    // HHG Panel
    expect(text).to.include('Shipment - Government moves all of your stuff (HHG)');
    expect(text).to.include('Movers Packing: Fri, May 11 - Mon, May 14');
    expect(text).to.include('Loading Truck: Tue, May 15');
    expect(text).to.include('Move in Transit:Wed, May 16 - Sun, May 20');
    expect(text).to.include('Delivery:Mon, May 21');
    expect(text).to.include(
      'Weight Estimate:2,000 lbs + 225 lbs pro-gear + 312 lbs spouse pro-gear Great! You appear within your weight allowance.',
    );

    // PPM Panel
    expect(text).to.include('Shipment - You move your stuff (PPM)');
    expect(text).to.include('Move Date: 05/20/2018');
    expect(text).to.include('Pickup ZIP Code:  90210');
    expect(text).to.include('Delivery ZIP Code:  30813');
    expect(text).not.to.include('Storage: Not requested');
    expect(text).to.include('Estimated Weight:  1,50');
    expect(text).to.include('Estimated PPM Incentive:  $4,362.66 - 4,821.88');
  });
}

function serviceMemberVerifiesPPMWeightsEdited() {
  cy.get('.ppm-container').should($div => {
    const text = $div.text();

    expect(text).to.include('Estimated Weight:  1,700 lbs');
    expect(text).to.include('Estimated PPM Incentive:  $4,858.82 - 5,370.28');
  });
}

function serviceMemberEditsPPMWeight() {
  cy.get('[data-cy="edit-ppm-weight"]').click();

  typeInInput({ name: 'weight_estimate', value: '1700' });

  cy.get('strong').contains('$4,858.82 - 5,370.28');
  cy.get('.subtext').contains('Originally $4,362.66 - 4,821.88');

  cy
    .get('button')
    .contains('Save')
    .click();
}
function serviceMemberVerifiesPPMDatesAndLocationsEdited() {
  cy.get('.ppm-container').should($div => {
    const text = $div.text();
    expect(text).to.include('Move Date: 05/28/2018');
    expect(text).to.include('Pickup ZIP Code:  91206');
    expect(text).to.include('Delivery ZIP Code:  30813');
  });
}
function serviceMemberEditsPPMDatesAndLocations() {
  cy.get('[data-cy="edit-ppm-dates"]').click();

  typeInInput({ name: 'planned_move_date', value: '5/28/2018' });
  typeInInput({ name: 'pickup_postal_code', value: '91206' });
  typeInInput({ name: 'destination_postal_code', value: '30813' });

  cy
    .get('button')
    .contains('Save')
    .click();
}

function serviceMemberVerifiesBackupContactInfoWasEdited() {
  cy.get('.review-content').should($div => {
    const text = $div.text();
    expect(text).to.include('Backup Contact: Backup Name');
    expect(text).to.include('Email:  backup@example.com');
    expect(text).to.include('Phone:  323-111-1111');
  });
}

function serviceMemberEditsBackupContactInfoSection() {
  cy
    .get('.review-content .edit-section-link')
    .eq(3)
    .click();

  typeInInput({ name: 'name', value: 'Backup Name' });
  typeInInput({ name: 'email', value: 'backup@example.com' });
  typeInInput({ name: 'telephone', value: '323-111-1111' });

  cy
    .get('button')
    .contains('Save')
    .click();
}

function serviceMemberVerifiesContactInfoWasEdited() {
  cy.get('.review-content').should($div => {
    const text = $div.text();
    expect(text).to.include('Best Contact Phone: 213-111-1111');
    expect(text).to.include('Alt. Phone: 222-222-2222');
    expect(text).to.include('Personal Email: hhgforppm@awarded.com');
    expect(text).to.include('Preferred Contact Method: Phone, Text, Email');
    expect(text).to.include('Current Mailing Address: 321 Any Street');
    expect(text).to.include('Los Angeles, CO 91206');

    expect(text).to.include('Backup Mailing Address: 333 Any Street');
    expect(text).to.include('P.O Box 54321');
    expect(text).to.include('Los Angeles, CT 91206');
  });
}

function serviceMemberEditsContactInfoSection() {
  cy
    .get('.review-content .edit-section-link')
    .eq(2)
    .click();

  typeInInput({ name: 'serviceMember.telephone', value: '213-111-1111' });
  typeInInput({ name: 'serviceMember.secondary_telephone', value: '222-222-2222' });
  typeInInput({ name: 'serviceMember.personal_email', value: 'hhgforppm@awarded.com' });
  cy.get('[type="checkbox"]').check({ force: true });
  typeInInput({ name: 'resAddress.street_address_1', value: '321 Any Street' });
  typeInInput({ name: 'resAddress.street_address_2', value: 'P.O Box 54321' });
  typeInInput({ name: 'resAddress.city', value: 'Los Angeles' });
  cy.get('select[name="resAddress.state"]').select('CO');
  typeInInput({ name: 'resAddress.postal_code', value: '91206' });

  typeInInput({ name: 'backupAddress.street_address_1', value: '333 Any Street' });
  typeInInput({ name: 'backupAddress.street_address_2', value: 'P.O Box 54321' });
  typeInInput({ name: 'backupAddress.city', value: 'Los Angeles' });
  cy.get('select[name="backupAddress.state"]').select('CT');
  typeInInput({ name: 'backupAddress.postal_code', value: '91206' });

  cy
    .get('button')
    .contains('Save')
    .click();
}

function serviceMemberVerifiesOrderWasEdited() {
  cy.get('.review-content').should($div => {
    const text = $div.text();
    expect(text).to.include('Orders Type: Permanent Change Of Station');
    expect(text).to.include('Orders Date:  05/26/2018');
    expect(text).to.include('Report-by Date: 09/01/2018');
    expect(text).to.include('New Duty Station:  NAS Fort Worth');
    expect(text).to.include('Dependents?:  No');
    expect(text).to.include('Orders Uploaded: 1');
  });
}

function serviceMemberEditsOrdersSection() {
  cy
    .get('.review-content .edit-section-link')
    .eq(1)
    .click();

  cy.get('select[name="orders_type"]').select('Permanent Change Of Station');
  typeInInput({ name: 'issue_date', value: '5/26/2018' });
  typeInInput({ name: 'report_by_date', value: '9/1/2018' });
  cy.get('input[type="radio"]').check('no', { force: true }); // checks yes for both radios on form
  cy.selectDutyStation('NAS Fort Worth', 'new_duty_station');
  cy
    .get('button')
    .contains('Save')
    .click();
}

function typeInInput({ name, value }) {
  cy
    .get(`input[name="${name}"]`)
    .clear()
    .type(value)
    .blur();
}

function serviceMemberVerifiesProfileWasEdited() {
  cy.get('.review-content').should($div => {
    const text = $div.text();
    console.log(text);
    expect(text).to.include('Name: Harry James Potter Sr');
    expect(text).to.include('Branch:Air Force');
    expect(text).to.include('Rank/Pay Grade: E-9');
    expect(text).to.include('DoD ID#: 9876543210');
    expect(text).to.include('Current Duty Station: NAS Fort Worth');
  });
}

function serviceMemberEditsProfileSection() {
  cy
    .get('.review-content .edit-section-link')
    .first()
    .click();

  typeInInput({ name: 'first_name', value: 'Harry' });
  typeInInput({ name: 'middle_name', value: 'James' });
  typeInInput({ name: 'last_name', value: 'Potter' });
  typeInInput({ name: 'suffix', value: 'Sr' });
  cy.get('select[name="affiliation"]').select('Air Force');
  cy.get('select[name="rank"]').select('E-9');
  cy
    .get('input[name="edipi"]')
    .clear()
    .type('9876543210');
  cy.selectDutyStation('NAS Fort Worth', 'current_station');
  cy
    .get('button')
    .contains('Save')
    .click();
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
    .contains('Add PPM (DITY) Move')
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
    expect(loc.pathname).to.match(/^\/$/);
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

  cy.get('.wizard-header').should('contain', 'Move Setup');
  cy.get('.wizard-header .progress-timeline .current').should('contain', 'Move Setup');
  cy
    .get('.wizard-header .progress-timeline .step')
    .last()
    .should('contain', 'Review');

  cy
    .get('input[name="planned_move_date"]')
    .should('have.value', '5/20/2018')
    .clear()
    .first()
    .type('9/2/2018{enter}')
    .blur();

  cy.get('input[name="pickup_postal_code"]').should('have.value', '90210');

  cy.get('input[name="destination_postal_code"]').should('have.value', '30813');

  cy.nextPage();
}

function serviceMemberSelectsWeightRange() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-size/);
  });

  cy.get('.wizard-header').should('contain', 'Move Setup');
  cy.get('.wizard-header .progress-timeline .current').should('contain', 'Move Setup');
  cy
    .get('.wizard-header .progress-timeline .step')
    .last()
    .should('contain', 'Review');

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

  cy.get('.wizard-header').should('contain', 'Move Setup');
  cy.get('.wizard-header .progress-timeline .current').should('contain', 'Move Setup');
  cy
    .get('.wizard-header .progress-timeline .step')
    .last()
    .should('contain', 'Review');

  // We usually poke the weight range slider to simulate user interaction,
  // but this can often move the slider handle by a pixel and throw off the numbers.
  // I'm commenting out this line in lieu of trying to build a slider interaction that can
  // verify that a desired weight is reached
  // cy.get('.rangeslider__handle').click();

  cy.get('.incentive').contains('$');

  cy.nextPage();
}

function serviceMemberCanReviewMoveSummary() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });

  cy.get('.wizard-header .usa-width-one-third').should('not.contain', 'Move Setup');
  cy.get('.wizard-header .usa-width-one-third').should('contain', 'Review');
  cy
    .get('.wizard-header .progress-timeline .step')
    .first()
    .should('contain', 'Move Setup');
  cy.get('.wizard-header .progress-timeline .current').should('contain', 'Review');
  cy.get('h3').should('not.contain', 'Profile and Orders');
  cy.get('h2').should('contain', 'Review Move Details');

  cy.get('body').should($div => expect($div.text()).not.to.include('Government moves all of your stuff (HHG)'));
  cy.get('.ppm-container').should($div => {
    const text = $div.text();
    expect(text).to.include('Shipment - You move your stuff (PPM)');
    expect(text).to.include('Move Date: 05/20/2018');
    expect(text).to.include('Pickup ZIP Code:  90210');
    expect(text).to.include('Delivery ZIP Code:  30813');
    expect(text).not.to.include('Storage: Not requested');
    expect(text).to.include('Estimated Weight:  1,50');
    expect(text).to.include('Estimated PPM Incentive:  $4,362.66 - 4,821.88');
  });

  cy.nextPage();
}
function serviceMemberCanSignAgreement() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-ppm-agreement/);
  });

  cy.get('.wizard-header').should('contain', 'Review');
  cy
    .get('.wizard-header .progress-timeline .step')
    .first()
    .should('contain', 'Move Setup');
  cy.get('.wizard-header .progress-timeline .current').should('contain', 'Review');

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

  cy.get('.usa-alert-success').contains("You've added a PPM shipment");
  cy
    .get('.usa-alert-success')
    .contains('Next, your shipment is awaiting approval and this can take up to 3 business days');

  cy.get('body').should($div => {
    expect($div.text()).to.include('Government Movers and Packers');
    expect($div.text()).to.include('Move your own stuff');
    expect($div.text()).to.not.include('Add PPM (DITY) Move');
  });

  cy.get('.usa-width-three-fourths').should($div => {
    const text = $div.text();
    // PPM information and details
    expect(text).to.include('Next Step: Wait for approval');
    expect(text).to.include('Weight (est.): 150');
    expect(text).to.include('Incentive (est.): $4,362.66 - 4,821.88');
  });
}
