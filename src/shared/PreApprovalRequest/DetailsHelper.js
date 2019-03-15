import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { Code35Form } from './Code35Form';
import { Code226Form } from './Code226Form';
import { get } from 'lodash';
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
  } else if (robustAccessorial && code.startsWith('226')) {
    return Code226Form;
  }
  return DefaultForm;
}

export function getDetailsComponent(code, robustAccessorial, isNewAccessorial) {
  return (code === '105B' || code === '105E') && isNewAccessorial ? Code105Details : DefaultDetails;
}
