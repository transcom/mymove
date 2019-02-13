import { DefaultDetails } from './DefaultDetails';
import { Code105Details } from './Code105Details';

export function getDetailComponent(code, robustAccessorial, initialValues) {
  code = code ? code.toLowerCase() : '';
  if (initialValues && !initialValues.crate_dimensions) return DefaultDetails;
  if ((code && code.startsWith('105b')) || (code.startsWith('105e') && robustAccessorial)) return Code105Details;
  return DefaultDetails;
}
