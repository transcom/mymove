import { DefaultDetails } from './DefaultDetails';
import { Code105Details } from './Code105Details';
import { has } from 'lodash';

export function getDetailComponent(code, robustAccessorial, initialValues) {
  code = code ? code.toLowerCase() : '';
  if (initialValues && !initialValues.crate_dimensions) return DefaultDetails;
  const hasDimensions = has(initialValues, 'crate_dimensions');
  if ((code.startsWith('105b') || code.startsWith('105e')) && hasDimensions && robustAccessorial) return Code105Details;
  return DefaultDetails;
}
