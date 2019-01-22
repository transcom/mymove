/* global cy */

export function tspUserVerifiesShipmentStatus(status) {
  cy.get('.shipment-status').contains(`Status: ${status}`);
}
