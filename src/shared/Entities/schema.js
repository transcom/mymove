/* eslint no-use-before-define: 0 */
import { schema } from 'normalizr';

// User
export const user = new schema.Entity('users');

// Service Member
export const serviceMember = new schema.Entity('serviceMembers', {
  user: user,
});

// Uploads
export const upload = new schema.Entity('uploads');
export const uploads = new schema.Array(upload);

// Documents
export const documentModel = new schema.Entity('documents', {
  uploads: uploads,
  service_member: serviceMember,
});

// MoveDocuments
export const moveDocument = new schema.Entity('moveDocuments', {
  document: documentModel,
});
export const moveDocuments = new schema.Array(moveDocument);

// PPMs
export const personallyProcuredMove = new schema.Entity(
  'personallyProcuredMove',
);
export const personallyProcuredMoves = new schema.Array(personallyProcuredMove);

// Addresses
export const address = new schema.Entity('address');

// Shipments
export const shipment = new schema.Entity('shipment');
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
});
export const moves = new schema.Array(move);
personallyProcuredMove.define({
  move: move,
});
moveDocument.define({
  move: move,
});

// Orders
export const order = new schema.Entity('orders', {
  uploaded_orders: documentModel,
  moves: moves,
});
