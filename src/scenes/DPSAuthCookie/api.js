import { getClient, checkResponse } from 'shared/Swagger/api';

export async function GetCookieURL(values) {
  const client = await getClient();
  const response = await client.apis.dps_auth.getCookieURL({
    cookie_name: values.cookie_name,
    dps_redirect_url: values.dps_redirect_url,
  });
  checkResponse(response, 'Failed to set DPS Auth cookie');
  return response.body;
}
