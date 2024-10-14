import { isBooleanFlagEnabled } from 'utils/featureFlags';

// List of states
export const statesList = [
  { value: 'AL', key: 'AL' },
  { value: 'AK', key: 'AK' },
  { value: 'AR', key: 'AR' },
  { value: 'AZ', key: 'AZ' },
  { value: 'CA', key: 'CA' },
  { value: 'CO', key: 'CO' },
  { value: 'CT', key: 'CT' },
  { value: 'DC', key: 'DC' },
  { value: 'DE', key: 'DE' },
  { value: 'FL', key: 'FL' },
  { value: 'GA', key: 'GA' },
  { value: 'HI', key: 'HI' },
  { value: 'IA', key: 'IA' },
  { value: 'ID', key: 'ID' },
  { value: 'IL', key: 'IL' },
  { value: 'IN', key: 'IN' },
  { value: 'KS', key: 'KS' },
  { value: 'KY', key: 'KY' },
  { value: 'LA', key: 'LA' },
  { value: 'MA', key: 'MA' },
  { value: 'MD', key: 'MD' },
  { value: 'ME', key: 'ME' },
  { value: 'MI', key: 'MI' },
  { value: 'MN', key: 'MN' },
  { value: 'MO', key: 'MO' },
  { value: 'MS', key: 'MS' },
  { value: 'MT', key: 'MT' },
  { value: 'NC', key: 'NC' },
  { value: 'ND', key: 'ND' },
  { value: 'NE', key: 'NE' },
  { value: 'NH', key: 'NH' },
  { value: 'NJ', key: 'NJ' },
  { value: 'NM', key: 'NM' },
  { value: 'NV', key: 'NV' },
  { value: 'NY', key: 'NY' },
  { value: 'OH', key: 'OH' },
  { value: 'OK', key: 'OK' },
  { value: 'OR', key: 'OR' },
  { value: 'PA', key: 'PA' },
  { value: 'RI', key: 'RI' },
  { value: 'SC', key: 'SC' },
  { value: 'SD', key: 'SD' },
  { value: 'TN', key: 'TN' },
  { value: 'TX', key: 'TX' },
  { value: 'UT', key: 'UT' },
  { value: 'VA', key: 'VA' },
  { value: 'VT', key: 'VT' },
  { value: 'WA', key: 'WA' },
  { value: 'WI', key: 'WI' },
  { value: 'WV', key: 'WV' },
  { value: 'WY', key: 'WY' },
];

export const unSupportedStates = [{ value: 'HI', key: 'HI' }];
export const unSupportedStatesDisabledAlaska = [
  { value: 'HI', key: 'HI' },
  { value: 'AK', key: 'AK' },
];

export const getUnSupportedStates = async () => {
  const enableAKFlag = await isBooleanFlagEnabled('enable_alaska');

  if (!enableAKFlag) {
    return unSupportedStatesDisabledAlaska;
  }

  return unSupportedStates;
};
