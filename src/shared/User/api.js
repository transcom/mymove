import { getClient, checkResponse } from 'shared/Swagger/api';
import * as Cookies from 'js-cookie';

export async function GetLoggedInUser() {
  const client = await getClient();
  const response = await client.apis.users.showLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export function LogoutUser() {
  const token = Cookies.get('masked_gorilla_csrf');
  console.log('============');
  console.log(token);
  console.log(document.cookie);
  console.log('============');
  const logoutEndpoint = '/auth/logout';
  fetch(logoutEndpoint, {
    method: 'POST',
    headers: { 'X-CSRF-Token': token },
    redirect: 'follow',
  }).then(response => {
    window.location = '/';
  });
}
