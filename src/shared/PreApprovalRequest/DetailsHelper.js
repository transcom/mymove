import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { Code35Form } from './Code35Form';
import { Code226Form } from './Code226Form';
import { get } from 'lodash';
import { Code35Details } from './Code35Details';
import { Code105Details } from './Code105Details';
import { DefaultDetails } from './DefaultDetails';
import { Code226Details } from './Code226Details';
import { Code125Form } from './Code125Form';
import { Code125Details } from './Code125Details';

export function getFormComponent(code, robustAccessorialFlag, initialValues) {
  code = code ? code.toLowerCase() : '';
  const isNew = !initialValues;
  if (code.startsWith('105b') || code.startsWith('105e')) {
    if (isNew || get(initialValues, 'crate_dimensions', false)) return Code105Form;
  } else if (code.startsWith('35')) {
    if (isNew || get(initialValues, 'estimate_amount_cents')) return Code35Form;
  } else if (code.startsWith('226')) {
    if (isNew || get(initialValues, 'actual_amount_cents')) return Code226Form;
  } else if (robustAccessorialFlag && code.startsWith('125')) {
    if (isNew || get(initialValues, 'address')) return Code125Form;
  }
  return DefaultForm;
}

export function getDetailsComponent(code, robustAccessorialFlag, isRobustAccessorial) {
  if (!isRobustAccessorial) return DefaultDetails;
  if (code === '105B' || code === '105E') return Code105Details;
  if (code === '35A') return Code35Details;
  if (code === '226A') return Code226Details;
  if (code.startsWith('125') && robustAccessorialFlag) return Code125Details;
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
  if (code.startsWith('125') && !item.address) {
    return false;
  }
  return true;
};
