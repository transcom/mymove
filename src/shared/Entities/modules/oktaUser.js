import { get } from 'lodash';

import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export function getOktaUser() {
  const swaggerTag = 'okta_profile.showOktaInfo';
  return swaggerRequest(getClient, swaggerTag, {});
}

// load Okta user
export function selectOktaUser(state) {
  return get(state, `entities.oktaUser`);
}
