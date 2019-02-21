import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const getCurrentUserInfoLabel = 'Users.getCurrentUser';

export function getCurrentUserInfo(label = getCurrentUserInfoLabel) {
  const swaggerTag = 'users.showLoggedInUser';
  return swaggerRequest(getClient, swaggerTag, {}, { label });
}

export function selectCurrentUser(state) {
  const currentUserInfo = Object.values(state.entities.users).first || {};
  const userAppInfo = state.user;
  return { ...currentUserInfo, ...userAppInfo };
}
