import { getClient, checkResponse } from 'shared/Swagger/api';

export async function SetDPSAuthCookie(cookieName) {
  const client = await getClient();
  const response = await client.apis.dps_auth.setDPSAuthCookie({
    cookie_name: cookieName,
  });
  checkResponse(response, 'Failed to set DPS Auth cookie');
  return response.body;
}
