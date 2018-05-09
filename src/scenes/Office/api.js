import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  debugger;
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

export async function GetAccountingAPI(moveId) {
  const client = await getClient();
  const response = await client.apis.office.showAccounting({
    moveId: moveId,
  });
  checkResponse(
    response,
    'failed to get accounting for move due to server error',
  );
  return response.body;
}

export async function UpdateAccountingAPI(moveId, payload) {
  const client = await getClient();
  const response = await client.apis.office.updateAccounting({
    moveId,
    patchAccounting: payload,
  });
  checkResponse(
    response,
    'failed to update accounting for move due to server error',
  );
  return response.body;
}
