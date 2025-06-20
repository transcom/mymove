/* eslint-disable import/prefer-default-export */
// utility functions related to API interactions

import Swagger from 'swagger-client';

import { checkResponse, getAdminClient, getClient, requestInterceptor } from 'shared/Swagger/api';

export const getQueriesStatus = (queries) => {
  // Queries should be an array of statuses returned by useQuery (https://react-query.tanstack.com/docs/api#usequery)
  return {
    // isIntialLoading is the react-query key designated for loading states (https://tanstack.com/query/v4/docs/guides/migrating-to-react-query-4#disabled-queries)
    isLoading: queries.some((q) => q.isInitialLoading),
    isFetching: queries.some((q) => q.isFetching),
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

export async function GetOktaUser() {
  const client = await getClient();
  const response = await client.apis.okta_profile.showOktaInfo({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export async function GetIsLoggedIn() {
  const client = await getClient();
  const response = await client.apis.users.isLoggedInUser({});
  checkResponse(response, 'failed to get user due to server error');
  return response.body;
}

export async function GetAdminUser() {
  const client = await getAdminClient();
  const response = await client.apis.User.getLoggedInAdminUser({});
  checkResponse(response, 'failed to get admin user due to server error');
  return response.body;
}

// logs a user out of MilMove and Okta
// redirects them back to their respective MilMove sign in page
export function LogoutUser() {
  const logoutEndpoint = '/auth/logout';
  const req = {
    url: logoutEndpoint,
    method: 'POST',
    credentials: 'same-origin', // Passes through CSRF cookies
    requestInterceptor,
  };
  return Swagger.http(req);
}

// logs a user out of MilMove & Okta
// redirects them back to the Okta sign in page
export function LogoutUserWithOktaRedirect() {
  const logoutEndpoint = '/auth/logoutOktaRedirect';
  const req = {
    url: logoutEndpoint,
    method: 'POST',
    credentials: 'same-origin',
    requestInterceptor,
  };
  return Swagger.http(req);
}

// updates a users server-side session
// with their new active role
export function UpdateActiveRoleServerSession(roleType) {
  const updateActiveRoleEndpoint = '/auth/activeRole';
  const req = {
    url: updateActiveRoleEndpoint,
    method: 'PATCH',
    credentials: 'same-origin',
    requestInterceptor,
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ roleType }),
  };
  return Swagger.http(req);
}
