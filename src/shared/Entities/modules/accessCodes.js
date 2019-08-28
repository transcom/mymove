import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const validateAccessCodeLabel = 'AccessCodes.getAccessCode';
const fetchAccessCodeLabel = 'AccessCodes.getAccessCodeForServiceMember';
const claimAccessCodeLabel = 'AccessCodes.updateAccessCodeForServiceMember';

export function fetchAccessCode(label = fetchAccessCodeLabel) {
  const swaggerTag = 'accesscode.fetchAccessCode';
  return swaggerRequest(getClient, swaggerTag);
}

export function validateAccessCode(code, label = validateAccessCodeLabel) {
  const swaggerTag = 'accesscode.validateAccessCode';
  return swaggerRequest(getClient, swaggerTag, { code }, { label });
}

export function claimAccessCode(accessCode, label = claimAccessCodeLabel) {
  const swaggerTag = 'accesscode.claimAccessCode';
  return swaggerRequest(getClient, swaggerTag, { accessCode }, { label });
}
