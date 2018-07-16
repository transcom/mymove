import { GetPpmIncentive } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

const GET_PPM_INCENTIVE = 'GET_PPM_INCENTIVE';
const CLEAR_PPM_INCENTIVE = 'CLEAR_PPM_INCENTIVE';
export const getIncentiveActionType = ReduxHelpers.generateAsyncActionTypes(
  GET_PPM_INCENTIVE,
);

export const getPpmIncentive = ReduxHelpers.generateAsyncActionCreator(
  GET_PPM_INCENTIVE,
  GetPpmIncentive,
);

export const clearPpmIncentive = () => ({ type: CLEAR_PPM_INCENTIVE });

function clearReducer(state, action) {
  if (action.type === CLEAR_PPM_INCENTIVE)
    return { ...state, calculation: null };
  return state;
}
const incentiveReducer = ReduxHelpers.generateAsyncReducer(
  GET_PPM_INCENTIVE,
  v => ({
    calculation: { ...v },
  }),
);

export default (state, action) => {
  const result = clearReducer(state, action);
  return incentiveReducer(result, action);
};
