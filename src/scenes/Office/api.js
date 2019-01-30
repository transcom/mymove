import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload } from 'shared/utils';

// MOVE QUEUE
export async function RetrieveMovesForOffice(queueType) {
  const client = await getClient();
  const response = await client.apis.queues.showQueue({
    queueType,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}

// MOVE
export async function LoadMove(moveId) {
  const client = await getClient();
  const response = await client.apis.moves.showMove({
    moveId,
  });
  checkResponse(response, 'failed to load move due to server error');
  return response.body;
}

// SHIPMENT
export async function PatchShipment(shipmentId, shipment) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.Shipment;
  const response = await client.apis.shipments.patchShipment({
    shipmentId,
    shipment: formatPayload(shipment, payloadDef),
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

// ORDERS
export async function LoadOrders(ordersId) {
  const client = await getClient();
  const response = await client.apis.orders.showOrders({
    ordersId,
  });
  checkResponse(response, 'failed to load orders due to server error');
  return response.body;
}

// SERVICE MEMBER
export async function LoadServiceMember(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.service_members.showServiceMember({
    serviceMemberId,
  });
  checkResponse(response, 'failed to load service member due to server error');
  return response.body;
}

export async function UpdateServiceMember(serviceMemberId, payload) {
  const client = await getClient();
  const response = await client.apis.service_members.patchServiceMember({
    serviceMemberId,
    patchServiceMemberPayload: payload,
  });
  checkResponse(response, 'failed to update service member due to server error');
  return response.body;
}

// BACKUP CONTACT
export async function LoadBackupContacts(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.backup_contacts.indexServiceMemberBackupContacts({
    serviceMemberId,
  });
  checkResponse(response, 'failed to load backup contacts due to server error');
  return response.body;
}

export async function UpdateBackupContact(backupContactId, payload) {
  const client = await getClient();
  const response = await client.apis.backup_contacts.updateServiceMemberBackupContact({
    backupContactId,
    updateServiceMemberBackupContactPayload: payload,
  });
  checkResponse(response, 'failed to load backup contacts due to server error');
  return response.body;
}

// PPM
export async function LoadPPMs(moveId) {
  const client = await getClient();
  const response = await client.apis.ppm.indexPersonallyProcuredMoves({
    moveId,
  });
  checkResponse(response, 'failed to load PPMs due to server error');
  return response.body;
}

// HHG invoice
export async function SendHHGInvoice(shipmentId) {
  const client = await getClient();
  const response = await client.apis.shipments.createAndSendHHGInvoice({
    shipmentId,
  });
  checkResponse(response, 'failed to send invoice to server error');
  return response.body;
}

// Reimbursement status
export async function ApproveReimbursement(reimbursementId) {
  const client = await getClient();
  const response = await client.apis.office.approveReimbursement({
    reimbursementId,
  });
  checkResponse(response, 'failed to approve reimbursement due to server error');
  return response.body;
}

// PPM attachments
export async function DownloadPPMAttachments(ppmId, docTypes) {
  const client = await getClient();
  const response = await client.apis.ppm.createPPMAttachments({
    personallyProcuredMoveId: ppmId,
    docTypes: docTypes,
  });
  checkResponse(response, 'failed to create PPM attachments PDF');
  return response.body;
}
