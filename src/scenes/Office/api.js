import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function FunctionNameGoesHere(payload) {
  // Pull from finished API definition
  const client = await getClient();
  const response = await client.apis.duty_stations.searchDutyStations({
    status: 0, // Parameters to come; how many to add?
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}
