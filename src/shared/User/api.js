import { getClient, checkResponse } from 'shared/api';

export async function GetLoggedInUser() {
  const client = await getClient();
  const response = await client.apis.users.showLoggedInUser({});
  checkResponse(response, 'failed to get move due to server error');
  return response.body;
}
