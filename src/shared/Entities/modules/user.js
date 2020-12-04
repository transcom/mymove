import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export const showLoggedInUserLabel = 'ServiceMember.showLoggedInUser';

export function showLoggedInUser() {
  const swaggerTag = 'users.showLoggedInUser';
  return swaggerRequest(getClient, swaggerTag, {}, { label: showLoggedInUserLabel });
}
