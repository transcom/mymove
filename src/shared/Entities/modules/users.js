import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const getCurrentUserInfoLabel = 'Users.getCurrentUser';

export function getCurrentUserInfo(label = getCurrentUserInfoLabel) {
  const swaggerTag = 'users.showLoggedInUser';
  return swaggerRequest(getClient, swaggerTag, {}, { label });
}

export function selectCurrentUser(state) {
  const currentUserInfo = Object.values(state.entities.users)[0] || {};
  const userAppInfo = state.user;
  const serviceMember = currentUserInfo.service_member;
  if (serviceMember) {
    return { email: serviceMember.personal_email, first_name: serviceMember.first_name, ...userAppInfo };
  } else {
    return { email: '', ...currentUserInfo, ...userAppInfo };
  }
}
