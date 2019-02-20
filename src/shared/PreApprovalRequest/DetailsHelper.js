import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { has } from 'lodash';
import { Code105Details } from './Code105Details';
import { DefaultDetails } from './DefaultDetails';

export function getFormComponent(code, robustAccessorial, initialValues) {
  code = code ? code.toLowerCase() : '';
  const hasDimensions = !initialValues || has(initialValues, 'crate_dimensions');
  if ((code.startsWith('105b') || code.startsWith('105e')) && hasDimensions && robustAccessorial) return Code105Form;
  return DefaultForm;
}

export function getDetailsComponent(code, robustAccessorial, isNewAccessorial) {
  return (code === '105B' || code === '105E') && robustAccessorial && isNewAccessorial
    ? Code105Details
    : DefaultDetails;
}
