import { getClient, checkResponse } from 'shared/api';

// MOVE QUEUE
export async function RetrieveMovesForOffice(queueType) {
  const client = await getClient();
  const response = await client.apis.queues.showQueue({
    queueType,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}

// ACCOUNTING
export async function LoadAccountingAPI(moveId) {
  const client = await getClient();
  const response = await client.apis.office.showAccounting({
    moveId: moveId,
  });
  checkResponse(
    response,
    'failed to load accounting for move due to server error',
  );
  return response.body;
}

export async function UpdateAccountingAPI(moveId, payload) {
  const client = await getClient();
  const response = await client.apis.office.patchAccounting({
    moveId,
    patchAccounting: payload,
  });
  checkResponse(
    response,
    'failed to update accounting for move due to server error',
  );
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

// BACKUP CONTACT
export async function LoadBackupContacts(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.backup_contacts.indexServiceMemberBackupContacts(
    {
      serviceMemberId,
    },
  );
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

// Move status
export async function ApproveBasics(moveId) {
  const client = await getClient();
  const response = await client.apis.office.approveMove({
    moveId,
  });
  checkResponse(response, 'failed to approve move gdue to server error');
  return response.body;
}
