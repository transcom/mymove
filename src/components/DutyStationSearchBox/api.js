import { getClient, checkResponse } from 'shared/Swagger/api';

export async function SearchDutyStations(query) {
  const client = await getClient();
  const response = await client.apis.duty_stations.searchDutyStations({
    search: query,
  });
  checkResponse(response, 'failed to query duty stations due to server error');
  return response.body;
}

export async function ShowAddress(addressId) {
  const client = await getClient();
  const response = await client.apis.addresses.showAddress({
    addressId,
  });
  checkResponse(response, 'failed to query address for duty station');
  return response.body;
}
