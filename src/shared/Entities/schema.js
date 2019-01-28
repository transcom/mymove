/* eslint no-use-before-define: 0 */
import { schema } from 'normalizr';

// User
export const user = new schema.Entity('users');

// Uploads
export const upload = new schema.Entity('uploads');
export const uploads = new schema.Array(upload);

// PPMs
export const personallyProcuredMove = new schema.Entity('personallyProcuredMove');
export const personallyProcuredMoves = new schema.Array(personallyProcuredMove);

// Addresses
export const address = new schema.Entity('addresses');
export const addresses = new schema.Array(address);

// Shipments
export const shipment = new schema.Entity('shipments');

export const shipments = new schema.Array(shipment);

export const serviceAgent = new schema.Entity('serviceAgents');

export const serviceAgents = new schema.Array(serviceAgent);

// Moves
export const move = new schema.Entity('moves', {
  personally_procured_moves: personallyProcuredMoves,
  shipments: shipments,
});
export const moves = new schema.Array(move);
personallyProcuredMove.define({
  move: move,
});

// Orders
export const order = new schema.Entity('orders', {
  moves: moves,
});

export const orders = new schema.Array(order);

// Service Member
export const serviceMember = new schema.Entity('serviceMembers', {
  user: user,
  orders: orders,
});

// Documents
export const documentModel = new schema.Entity('documents', {
  uploads: uploads,
  service_member: serviceMember,
});
order.define({
  uploaded_orders: documentModel,
});

// MoveDocuments
export const moveDocument = new schema.Entity('moveDocuments', {
  document: documentModel,
});

export const moveDocuments = new schema.Array(moveDocument);
moveDocument.define({
  move: move,
});

export const moveDocumentPayload = moveDocument;

// Tariff400ngItems
export const tariff400ngItem = new schema.Entity('tariff400ngItems');
export const tariff400ngItems = new schema.Array(tariff400ngItem);

// Invoice
export const invoice = new schema.Entity('invoices');
export const invoices = new schema.Array(invoice);

// ShipmentLineItem
export const shipmentLineItem = new schema.Entity('shipmentLineItems', {
  tariff400ng_item: tariff400ngItem,
  invoice: invoice,
});
export const shipmentLineItems = new schema.Array(shipmentLineItem);

// AvailableMoveDates
export const availableMoveDates = new schema.Entity('availableMoveDates', {}, { idAttribute: 'start_date' });

// MoveDatesSummary
export const moveDatesSummary = new schema.Entity('moveDatesSummaries');

// TransportationServiceProviders
export const transportationServiceProvider = new schema.Entity('transportationServiceProviders');
