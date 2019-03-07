/* global cy */
export function fillAndSavePreApprovalRequest() {
  // Click on add Pre Approval Request
  cy
    .get('.add-request')
    .contains('Add a request')
    .click();

  // Verify tariff items that don't require approval are not loaded into drop down
  cy
    .get('.tariff400-select #react-select-2-input')
    .first()
    .type('Linehaul Transportation{downarrow}{enter}', { force: true, delay: 75 });
  cy.get('.tariff400__single-value').should('not.exist');

  //  Enter details in form and create pre approval request
  cy.selectTariff400ngItem('Article: Motorcycle');

  cy.get('select[name="location"]').select('ORIGIN');

  cy.typeInInput({ name: 'quantity_1', value: 2 });
  cy.typeInTextarea({ name: 'notes', value: `notes notes` });

  cy
    .get('button')
    .contains('Save & Close')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save & Close')
    .click();
}

export function editPreApprovalRequest() {
  cy
    .get('[data-test=edit-request]')
    .first()
    .click();

  cy.typeInTextarea({ name: 'notes', value: `edited` });
  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();
}

export function approvePreApprovalRequest() {
  cy.get('[data-test=approve-request]').should('exist');
  cy
    .get('[data-test=approve-request]')
    .first()
    .click();
}

export function deletePreApprovalRequest() {
  cy
    .get('[data-test=delete-request]')
    .first()
    .click();
  cy
    .get('button')
    .contains('No, do not delete')
    .click();

  cy
    .get('[data-test=delete-request]')
    .first()
    .click();

  cy
    .get('button')
    .contains('Yes, delete')
    .click();
}
