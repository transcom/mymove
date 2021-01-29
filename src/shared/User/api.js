import Swagger from 'swagger-client';
import qs from 'query-string';

import { getClient, checkResponse, requestInterceptor } from 'shared/Swagger/api';

export async function GetLoggedInUser() {
  const client = await getClient();
  const response = await client.apis.users.showLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export async function GetIsLoggedIn() {
  const client = await getClient();
  const response = await client.apis.users.isLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export async function LogoutUser(timedout) {
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
    let redirect_url = timedout ? qs.stringifyUrl({ url: resp.text, fragmentIdentifier: 'timedout' }) : resp.text;
    window.location.href = redirect_url;
  } catch (err) {
    // Failure to logout should return user to homepage
    window.location.href = '/';
  }
}
