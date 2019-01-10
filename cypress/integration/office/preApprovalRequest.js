import {
  fillAndSavePreApprovalRequest,
  fillInvalidPreApprovalRequest,
  editPreApprovalRequest,
  approvePreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/testCreatePreApprovalRequest';

/* global cy */
describe('office user interacts with pre approval request panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  // it('office user creates pre approval request', function() {
  //   officeUserCannotAddInvalidPreApprovalRequest();
  // });
  // it('office user edits pre approval request', function() {
  //   officeUserEditsPreApprovalRequest();
  // });
  // it('office user approves pre approval request', function() {
  //   officeUserApprovesPreApprovalRequest();
  // });
  // it('office user deletes pre approval request', function() {
  //   officeUserDeletesPreApprovalRequest();
  // });

  it('office user iterates through every pre approval request', () => {
    officeUserIterateThroughAllPARS();
  });
});

function officeUserIterateThroughAllPARS() {
  // Open new moves queue
  cy.visit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('DOOB')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  let PARS = [
    '4A',
    '4B',
    '28A',
    '28B',
    '28C',
    '35A',
    '105B',
    '105D',
    '105E',
    '120A',
    '120B',
    '120C',
    '120D',
    '120E',
    '120F',
    '125A',
    '125B',
    '125C',
    '125D',
    '130A',
    '130B',
    '130C',
    '130D',
    '130E',
    '130F',
    '130G',
    '130H',
    '130I',
    '130J',
    '175A',
    '225A',
    '225B',
    '226A',
  ];
  PARS.forEach(PAR => {
    //add a pre approval request first
    fillAndSavePreApprovalRequest(PAR);

    // wait a second for ui to catch up
    cy.wait(1000);

    // Verify data has been saved in the UI
    cy.get('td').contains(PAR);

    //approve PAR
    cy.get('[data-test=approve-request]').click({ multiple: true });
  });

  //invoice shipment
  cy
    .get('button')
    .contains('Approve Payment')
    .click()
    .then(() => {
      cy
        .get('button')
        .contains('Approve')
        .click()
        .then(() => {
          expect(cy.get('.invoice-panel').contains('Success!'));
        });
    });
}

function officeUserCannotAddInvalidPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  fillInvalidPreApprovalRequest();
}
function officeUserEditsPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  editPreApprovalRequest();
  // Verify data has been saved in the UI
  cy.get('td').contains('notes notes edited');
}

function officeUserApprovesPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  approvePreApprovalRequest();
  cy.get('.pre-approval-panel td').contains('Approved');
}

function officeUserDeletesPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  deletePreApprovalRequest();
  cy
    .get('.pre-approval-panel td')
    .first()
    .should('not.contain', 'Bulky Article: Motorcycle/Rec vehicle');
}
