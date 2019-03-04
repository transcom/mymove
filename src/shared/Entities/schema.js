/* eslint no-use-before-define: 0 */
import { schema } from 'normalizr';

// User
export const user = new schema.Entity('users');

// Uploads
export const upload = new schema.Entity('uploads');
export const uploads = new schema.Array(upload);

// PPMs
export const reimbursement = new schema.Entity('reimbursements');
export const personallyProcuredMove = new schema.Entity('personallyProcuredMoves');
personallyProcuredMove.define({
  advance: reimbursement,
});

export const personallyProcuredMoves = new schema.Array(personallyProcuredMove);
export const indexPersonallyProcuredMove = personallyProcuredMoves;

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

// Orders
export const orders = new schema.Entity('orders');
orders.define({
  moves: moves,
});
export const ordersArray = new schema.Array(orders);

// ServiceMemberBackupContacts
export const backupContact = new schema.Entity('backupContacts');

export const backupContacts = new schema.Array(backupContact);

export const indexServiceMemberBackupContacts = new schema.Array(backupContact);

export const serviceMemberBackupContact = backupContact;

// Service Member
export const serviceMember = new schema.Entity('serviceMembers', {
  backup_contacts: backupContacts,
  user: user,
});

// Documents
export const documentModel = new schema.Entity('documents', {
  uploads: uploads,
  service_member: serviceMember,
});
orders.define({
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
  invoice: invoice,
});
export const shipmentLineItems = new schema.Array(shipmentLineItem);

// AvailableMoveDates
export const availableMoveDates = new schema.Entity('availableMoveDates', {}, { idAttribute: 'start_date' });

// MoveDatesSummary
export const moveDatesSummary = new schema.Entity('moveDatesSummaries');

// TransportationServiceProviders
export const transportationServiceProvider = new schema.Entity('transportationServiceProviders');

// StorageInTransits
export const storageInTransit = new schema.Entity('storageInTransits');

export const storageInTransits = new schema.Array(storageInTransit);
