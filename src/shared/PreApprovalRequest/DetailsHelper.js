import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { Code35Form } from './Code35Form';
import { get } from 'lodash';
import { Code35Details } from './Code35Details';
import { Code105Details } from './Code105Details';
import { DefaultDetails } from './DefaultDetails';

export function getFormComponent(code, robustAccessorial, initialValues) {
  code = code ? code.toLowerCase() : '';
  const hasCrateDimensions = get(initialValues, 'crate_dimensions', false);
  const isNew = !initialValues;
  if (code.startsWith('105b') || code.startsWith('105e')) {
    if (isNew || hasCrateDimensions) return Code105Form;
  } else if (robustAccessorial && code.startsWith('35')) {
    return Code35Form;
  }
  return DefaultForm;
}

export function getDetailsComponent(code, isNewAccessorial) {
  if (!isNewAccessorial) return DefaultDetails;
  if (code === '105B' || code === '105E') return Code105Details;
  if (code === '35A') return Code35Details;
  return DefaultDetails;
}
