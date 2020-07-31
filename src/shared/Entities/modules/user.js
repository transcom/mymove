import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const showLoggedInUserLabel = 'ServiceMember.showLoggedInUser';

export function showLoggedInUser() {
  const swaggerTag = 'users.showLoggedInUser';
  return swaggerRequest(getClient, swaggerTag, {}, { label: showLoggedInUserLabel });
}

export function selectLoggedInUser(state) {
  if (state.entities.user) {
    return Object.values(state.entities.user)[0];
  }
  return {};
}
