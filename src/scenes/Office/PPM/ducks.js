import { GetPpmIncentive } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

const actionName = 'GET_PPM_INCENTIVE';
export const getIncentiveActionType = ReduxHelpers.generateAsyncActionTypes(
  actionName,
);

export const getPpmIncentive = ReduxHelpers.generateAsyncActionCreator(
  actionName,
  GetPpmIncentive,
);

export default ReduxHelpers.generateAsyncReducer(actionName, v => ({
  calculation: { ...v },
}));
