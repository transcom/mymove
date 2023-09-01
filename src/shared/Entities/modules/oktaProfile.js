import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export function getOktaProfile() {
  const swaggerTag = 'okta_profile.showOktaInfo';
  return swaggerRequest(getClient, swaggerTag);
}
