/* global cy, Cypress*/
export default function createOrders() {
  cy.fixture('orders.json').then(orders =>
    cy.window().then(window => {
      orders.service_member_id =
        window._state.serviceMember.currentServiceMember.id;
      return cy
        .request('internal/duty_stations?search=fort%20worth')
        .then(result => {
          orders.new_duty_station_id = Cypress._.get(result, 'body.[0].id');
          return cy.request('POST', `/internal/orders`);
        });
      //         .then(result => {
      //           const ordersId  =Cypress._.get(result, 'body.[0].id');
      //           cy.fixture("orders.jpg").as("orders")

      //   // convert the logo base64 string to a blob
      //   return Cypress.Blob.base64StringToBlob(this.logo, "image/png").then((blob) => {

      //     // pass the blob to the fileupload jQuery plugin
      //     // used in your application's code
      //     // which initiates a programmatic upload
      //     $input.fileupload("add", { files: blob })
      // })
      //         });
    }),
  );
}
