/* eslint no-use-before-define: 0 */
import { schema } from 'normalizr';

// User
export const user = new schema.Entity('users');

// Uploads
export const upload = new schema.Entity('uploads');
export const uploads = new schema.Array(upload);

// PPMs
export const personallyProcuredMove = new schema.Entity(
  'personallyProcuredMove',
);
export const personallyProcuredMoves = new schema.Array(personallyProcuredMove);

// Addresses
export const address = new schema.Entity('addresses');
export const addresses = new schema.Array(address);

// Shipments
export const shipment = new schema.Entity('shipments');
shipment.define({
  pickup_address: address,
  secondary_pickup_address: address,
  delivery_address: address,
  partial_sit_delivery_address: address,
});

export const shipments = new schema.Array(shipment);

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

// AvailableMoveDates
export const availableMoveDates = new schema.Entity(
  'availableMoveDates',
  {},
  { idAttribute: 'start_date' },
);

// MoveDatesSummary
export const moveDatesSummary = new schema.Entity('moveDatesSummaries');
