import { DefaultDetails } from './DefaultDetails';
import { Code105Details } from './Code105Details';

export function getDetailComponent(code, robustAccessorial) {
  code = code ? code.toLowerCase() : '';
  if (code && code.startsWith('105') && !code.startsWith('105d') && robustAccessorial) return Code105Details;
  return DefaultDetails;
}
