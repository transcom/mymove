import Swagger from 'swagger-client';
import * as Cookies from 'js-cookie';

import { getClient, checkResponse } from 'shared/Swagger/api';

export async function GetLoggedInUser() {
  const client = await getClient();
  const response = await client.apis.users.showLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export function LogoutUser() {
  const token = Cookies.get('masked_gorilla_csrf');
  const logoutEndpoint = '/auth/logout';
  const req = {
    url: logoutEndpoint,
    method: 'POST',
    headers: { 'X-CSRF-Token': token },
  };
  Swagger.http(req)
    .then(response => {
      window.location = '/';
    })
    .catch(err => {
      console.log(err);
    });
}
