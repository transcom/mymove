import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function RetrieveDutyStations(payload) {
  const client = await getClient();
  const response = await client.apis.duty_stations.searchDutyStations({
    branch: 'AIRFORCE',
    search: 'afb',
  });
  checkResponse(response, 'failed to create move due to server error');
  return response.body;
}
