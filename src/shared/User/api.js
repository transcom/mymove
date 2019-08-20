import Swagger from 'swagger-client';

import { getClient, checkResponse, requestInterceptor } from 'shared/Swagger/api';

export async function GetLoggedInUser() {
  const client = await getClient();
  const response = await client.apis.users.showLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export async function LogoutUser() {
  const logoutEndpoint = '/auth/logout';
  const req = {
    url: logoutEndpoint,
    method: 'POST',
    credentials: 'same-origin', // Passes through CSRF cookies
    requestInterceptor,
  };
  try {
    // Successful logout should return a redirect url
    let resp = await Swagger.http(req);
    window.location.href = resp.text;
  } catch (err) {
    // Failure to logout should return user to homepage
    window.location.href = '/';
  }
}
