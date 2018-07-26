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

// MovingExpenseDocuments
export const movingExpenseDocument = new schema.Entity(
  'movingExpenseDocuments',
  {
    document: documentModel,
  },
);
export const movingExpenseDocuments = new schema.Array(movingExpenseDocument);

// PPMs
export const personallyProcuredMove = new schema.Entity(
  'personallyProcuredMove',
);
export const personallyProcuredMoves = new schema.Array(personallyProcuredMove);

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
