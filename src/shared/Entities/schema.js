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

export const currentMove = new schema.Array(move);
export const previousMoves = new schema.Array(move);

export const multiMoves = new schema.Entity('multiMoves', {
  currentMove: currentMove,
  previousMoves: previousMoves,
});

// Orders

export const order = new schema.Entity('orders');
export const orders = new schema.Entity('orders');
orders.define({
  moves,
});

export const ordersArray = new schema.Array(orders);

// ServiceMemberBackupContacts
export const backupContact = new schema.Entity('backupContacts');

export const backupContacts = new schema.Array(backupContact);

export const indexServiceMemberBackupContacts = new schema.Array(backupContact);

export const serviceMemberBackupContact = backupContact;

// DutyLocations and TransportationOffices
export const transportationOffice = new schema.Entity('transportationOffices');
export const transportationOffices = new schema.Array(transportationOffice);
export const dutyLocation = new schema.Entity('dutyLocations', {
  transportation_office: transportationOffice,
});
export const dutyLocations = new schema.Array(dutyLocation);

// Service Member
export const serviceMember = new schema.Entity('serviceMembers', {
  backup_contacts: backupContacts,
  user,
  orders: ordersArray,
});

// Okta Profile
export const oktaUser = new schema.Entity('oktaUser');

// Loggedin User
export const loggedInUser = new schema.Entity('user', {
  service_member: serviceMember,
});

// Documents
export const document = new schema.Entity('documents', {
  uploads,
  service_member: serviceMember,
});
export const documents = new schema.Array(document);

// MoveDocuments
export const moveDocument = new schema.Entity('moveDocuments', {
  document,
});

export const moveDocuments = new schema.Array(moveDocument);
moveDocument.define({
  move,
});

export const moveDocumentPayload = moveDocument;

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

// StorageInTransits
export const storageInTransit = new schema.Entity('storageInTransits');

export const storageInTransits = new schema.Array(storageInTransit);

export const ppmSitEstimate = new schema.Entity('ppmSitEstimate');

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
export const createdCustomer = new schema.Entity('createdCustomer');

// Payment Requests
export const paymentRequest = new schema.Entity('paymentRequests', {
  serviceItems: paymentServiceItems,
});

export const paymentRequests = new schema.Array(paymentRequest);

export const shipmentPaymentSITBalance = new schema.Entity(
  'shipmentsPaymentSITBalance',
  {},
  { idAttribute: 'shipmentID' },
);
export const shipmentsPaymentSITBalance = new schema.Array(shipmentPaymentSITBalance);

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

export const customerSupportRemark = new schema.Entity('customerSupportRemark');
export const customerSupportRemarks = new schema.Array(customerSupportRemark);

export const evaluationReport = new schema.Entity('evaluationReport');
export const evaluationReports = new schema.Array(evaluationReport);

export const searchMove = new schema.Entity('searchMoves');
export const searchMoves = new schema.Array(searchMove);
export const searchMovesResult = new schema.Entity('searchMovesResult');
