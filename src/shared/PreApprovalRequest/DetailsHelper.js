import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { has } from 'lodash';

export function getFormComponent(code, robustAccessorial, initialValues) {
  code = code ? code.toLowerCase() : '';
  const hasDimensions = !initialValues || has(initialValues, 'crate_dimensions');
  if ((code.startsWith('105b') || code.startsWith('105e')) && hasDimensions && robustAccessorial) return Code105Form;
  return DefaultForm;
}
