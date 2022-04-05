import { checkResponse, getClient } from 'shared/Swagger/api';

export async function SearchDutyLocations(query) {
  const client = await getClient();
  const response = await client.apis.duty_locations.searchDutyLocations({
    search: query,
  });
  checkResponse(response, 'failed to query duty locations due to server error');
  return response.body;
}

export async function ShowAddress(addressId) {
  const client = await getClient();
  const response = await client.apis.addresses.showAddress({
    addressId,
  });
  checkResponse(response, 'failed to query address for duty location');
  return response.body;
}
