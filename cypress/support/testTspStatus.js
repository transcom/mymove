/* global cy */

export function tspUserVerifiesShipmentStatus(status) {
  cy.get('[data-cy="shipment-status"]').contains(`Status: ${status}`);
}
