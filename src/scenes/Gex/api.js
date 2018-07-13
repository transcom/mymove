import { getClient, checkResponse } from 'shared/api';

export async function SendGexRequest(payload) {
  console.log('sending');
  console.log(payload);
  const client = await getClient();
  const response = await client.apis.gex.sendGexRequest({
    sendGexRequestPayload: payload,
  });
  checkResponse(response, 'failed to send GEX request');
}
