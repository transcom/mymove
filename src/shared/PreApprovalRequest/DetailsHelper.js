import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { Code35Form } from './Code35Form';
import { Code226Form } from './Code226Form';
import { get } from 'lodash';
import { Code35Details } from './Code35Details';
import { Code105Details } from './Code105Details';
import { DefaultDetails } from './DefaultDetails';
import { Code226Details } from './Code226Details';

export function getFormComponent(code, robustAccessorial, initialValues) {
  code = code ? code.toLowerCase() : '';
  const isNew = !initialValues;
  if (code.startsWith('105b') || code.startsWith('105e')) {
    if (isNew || get(initialValues, 'crate_dimensions', false)) return Code105Form;
  } else if (robustAccessorial && code.startsWith('35')) {
    if (isNew || get(initialValues, 'estimate_amount_cents')) return Code35Form;
  } else if (robustAccessorial && code.startsWith('226')) {
    if (isNew || get(initialValues, 'actual_amount_cents')) return Code226Form;
  }
  return DefaultForm;
}

export function getDetailsComponent(code, robustAccessorialFlag, isRobustAccessorial) {
  if (!isRobustAccessorial) return DefaultDetails;
  if (code === '105B' || code === '105E') return Code105Details;
  if (code === '35A' && robustAccessorialFlag) return Code35Details;
  if (code === '226A' && robustAccessorialFlag) return Code226Details;
  return DefaultDetails;
}

export const isRobustAccessorial = item => {
  if (!item) return false;

  const code = item.tariff400ng_item.code;
  if ((code === '105B' || code === '105E') && !item.crate_dimensions) {
    return false;
  }
  if (code === '35A' && !item.estimate_amount_cents) {
    return false;
  }
  if (code === '226A' && !item.actual_amount_cents) {
    return false;
  }
  return true;
};
