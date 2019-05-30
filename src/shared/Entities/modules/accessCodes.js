import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

const validateAccessCodeLabel = 'AccessCodes.getAccessCode';
const claimAccessCodeLabel = 'AccessCodes.updateAccessCodeForServiceMember';

export function validateAccessCode(code, label = validateAccessCodeLabel) {
  const swaggerTag = 'accesscode.validateAccessCode';
  return swaggerRequest(getPublicClient, swaggerTag, { code }, { label });
}

export function claimAccessCode(accessCode, label = claimAccessCodeLabel) {
  const swaggerTag = 'accesscode.claimAccessCode';
  return swaggerRequest(getPublicClient, swaggerTag, { accessCodePayload: accessCode }, { label });
}
