import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

const validateAccessCodeLabel = 'ValidateAccessCode.validateAccessCode';
const claimAccessCodeLabel = 'AccessCode.claimAccessCode';

export function validateAccessCode(code, label = validateAccessCodeLabel) {
  const swaggerTag = 'accesscode.validateAccessCode';
  return swaggerRequest(getPublicClient, swaggerTag, { code }, { label });
}

export function claimAccessCode(code, serviceMemberId, label = claimAccessCodeLabel) {
  //const swaggerTag = 'accesscode.claimAccessCode';
  console.log('Claiming access code');
  return;
}
