import { getClient, checkResponse } from 'shared/Swagger/api';

export async function GetCookieURL(cookieName) {
  const client = await getClient();
  const response = await client.apis.dps_auth.getCookieURL({
    cookie_name: cookieName,
  });
  checkResponse(response, 'Failed to set DPS Auth cookie');
  return response.body;
}
