/* eslint-disable import/prefer-default-export */
// utility functions related to API interactions

import Swagger from 'swagger-client';
import qs from 'query-string';

import { getClient, checkResponse, requestInterceptor } from 'shared/Swagger/api';

export const getQueriesStatus = (queries) => {
  // Queries should be an array of statuses returned by useQuery (https://react-query.tanstack.com/docs/api#usequery)
  return {
    isLoading: queries.some((q) => q.isLoading),
    isError: queries.some((q) => q.isError),
    isSuccess: queries.every((q) => q.isSuccess),
    errors: queries.reduce((errors, q) => (q.error ? [...errors, q.error] : errors), []),
  };
};

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
    const resp = await Swagger.http(req);
    const redirectUrl = timedout ? qs.stringifyUrl({ url: resp.text, fragmentIdentifier: 'timedout' }) : resp.text;
    window.location.href = redirectUrl;
  } catch (err) {
    // Failure to logout should return user to homepage
    window.location.href = '/';
  }
}
