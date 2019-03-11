import Swagger from 'swagger-client';

import { getClient, checkResponse, requestInterceptor } from 'shared/Swagger/api';

export async function GetLoggedInUser() {
  const client = await getClient();
  const response = await client.apis.users.showLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export function LogoutUser() {
  const logoutEndpoint = '/auth/logout';
  const req = {
    url: logoutEndpoint,
    method: 'POST',
    credentials: 'same-origin',
    requestInterceptor,
  };
  Swagger.http(req)
    .then(response => {
      window.location = '/';
    })
    .catch(err => {
      console.log(err);
    });
}
