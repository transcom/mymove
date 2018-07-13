import { GetPpmIncentive } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';
import { combineReducers } from 'redux';

const actionName = 'GET_PPM_INCENTIVE';
export const getIncentiveActionType = ReduxHelpers.generateAsyncActionTypes(
  actionName,
);

export const getPpmIncentive = ReduxHelpers.generateAsyncActionCreator(
  actionName,
  GetPpmIncentive,
);

const reducer = ReduxHelpers.generateAsyncReducer(actionName, v => ({
  calculation: { ...v },
}));

//this is to put reducer in it's own slice
export default combineReducers({ incentive: reducer });
