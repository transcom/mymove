import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function RetrieveMovesForOffice(queueType) {
  const client = await getClient();
  const response = await client.apis.queues.showQueue({
    queueType,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}

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

export async function LoadMove(moveId) {
  const client = await getClient();
  const response = await client.apis.office.showMove({
    moveId,
  });
  checkResponse(response, 'failed to load move due to server error');
  return response.body;
}
