/* eslint-disable import/prefer-default-export */
// utility functions related to API interactions

import Swagger from 'swagger-client';

import { checkResponse, getClient, requestInterceptor } from 'shared/Swagger/api';

export const getQueriesStatus = (queries) => {
  // Queries should be an array of statuses returned by useQuery (https://react-query.tanstack.com/docs/api#usequery)
  return {
    // isIntialLoading is the react-query key designated for loading states (https://tanstack.com/query/v4/docs/guides/migrating-to-react-query-4#disabled-queries)
    isLoading: queries.some((q) => q.isInitialLoading),
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
