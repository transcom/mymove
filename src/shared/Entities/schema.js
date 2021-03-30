/* eslint no-use-before-define: 0 */
import { schema } from 'normalizr';

// User
export const user = new schema.Entity('users');

// Uploads
export const upload = new schema.Entity('upload');
export const uploads = new schema.Array(upload);

// PPMs
export const reimbursement = new schema.Entity('reimbursements');

export const personallyProcuredMove = new schema.Entity('personallyProcuredMoves');
personallyProcuredMove.define({
  advance: reimbursement,
});

export const personallyProcuredMoves = new schema.Array(personallyProcuredMove);
export const indexPersonallyProcuredMove = personallyProcuredMoves;

// MTO Shipments
export const mtoShipment = new schema.Entity('mtoShipments');
export const mtoShipments = new schema.Array(mtoShipment);

// Shipments
export const shipment = new schema.Entity('shipments');
export const shipments = new schema.Array(shipment);

// Addresses
export const address = new schema.Entity('addresses');
export const addresses = new schema.Array(address);

export const serviceAgent = new schema.Entity('serviceAgents');

export const serviceAgents = new schema.Array(serviceAgent);

// Moves
export const move = new schema.Entity('moves', {
  personally_procured_moves: personallyProcuredMoves,
  mto_shipments: mtoShipments,
});
export const moves = new schema.Array(move);

// Orders

export const order = new schema.Entity('orders');
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

// DutyStations and TransportationOffices
export const transportationOffice = new schema.Entity('transportationOffices');
export const transportationOffices = new schema.Array(transportationOffice);
export const dutyStation = new schema.Entity('dutyStations', {
  transportation_office: transportationOffice,
});
export const dutyStations = new schema.Array(dutyStation);

// Service Member
export const serviceMember = new schema.Entity('serviceMembers', {
  backup_contacts: backupContacts,
  user,
  orders: ordersArray,
});

// Loggedin User
export const loggedInUser = new schema.Entity('user', {
  service_member: serviceMember,
});

// Documents
export const document = new schema.Entity('documents', {
  uploads: uploads,
  service_member: serviceMember,
});
export const documents = new schema.Array(document);

// MoveDocuments
export const moveDocument = new schema.Entity('moveDocuments', {
  document: document,
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

// Signed Certificate
export const signedCertification = new schema.Entity('signedCertifications');

export const signedCertifications = new schema.Array(signedCertification);

// PPM Weight Estimate Range
export const ppmEstimateRange = new schema.Entity('ppmEstimateRanges');

// AvailableMoveDates
export const availableMoveDates = new schema.Entity(
  'availableMoveDates',
  {},
  {
    idAttribute: 'start_date',
  },
);

// MoveDatesSummary
export const moveDatesSummary = new schema.Entity('moveDatesSummaries');

// TransportationServiceProviders
export const transportationServiceProvider = new schema.Entity('transportationServiceProviders');

// StorageInTransits
export const storageInTransit = new schema.Entity('storageInTransits');

export const storageInTransits = new schema.Array(storageInTransit);

export const ppmSitEstimate = new schema.Entity('ppmSitEstimate');

// AccessCodes
export const accessCode = new schema.Entity('accessCodes');

// MTO Service Items
export const mtoServiceItem = new schema.Entity('mtoServiceItems');
export const mtoServiceItems = new schema.Array(mtoServiceItem);

// Payment Service Items
export const paymentServiceItem = new schema.Entity('paymentServiceItems');
export const paymentServiceItems = new schema.Array(paymentServiceItem);

// Move Task Orders
export const moveTaskOrder = new schema.Entity('moveTaskOrders');
export const moveTaskOrders = new schema.Array(moveTaskOrder);

// Customer
export const customer = new schema.Entity('customer');

// Payment Requests
export const paymentRequest = new schema.Entity('paymentRequests', {
  serviceItems: paymentServiceItems,
});

export const paymentRequests = new schema.Array(paymentRequest);

// MTO Agents
export const mtoAgent = new schema.Entity('mtoAgents');
export const mtoAgents = new schema.Array(mtoAgent);

// Queues
export const queueMove = new schema.Entity('queueMoves');
export const queueMoves = new schema.Array(queueMove);
export const queueMovesResult = new schema.Entity('queueMovesResult');

export const queuePaymentRequest = new schema.Entity('queuePaymentRequests');
export const queuePaymentRequests = new schema.Array(queuePaymentRequest);
export const queuePaymentRequestsResult = new schema.Entity('queuePaymentRequestsResult');
