import { LoadAccountingAPI, UpdateAccountingAPI } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

function generatePanelReducer(loadActionType, updateActionType) {
  const loadAsyncActionType = ReduxHelpers.generateAsyncActionTypes(
    loadActionType,
  );

  const updateAsyncActionType = ReduxHelpers.generateAsyncActionTypes(
    updateActionType,
  );

  const initialState = {
    isLoading: false,
    isUpdating: false,
    hasLoadError: false,
    hasLoadSuccess: null,
    hasUpdateError: false,
    hasUpdateSuccess: null,
  };

  return function panelReducer(state = initialState, action) {
    switch (action.type) {
      case loadAsyncActionType.start:
        return Object.assign({}, state, {
          isLoading: true,
          hasLoadSuccess: false,
        });
      case loadAsyncActionType.success:
        return Object.assign({}, state, {
          isLoading: false,
          accounting: action.payload,
          hasLoadSuccess: true,
          hasLoadError: false,
        });
      case loadAsyncActionType.failure:
        return Object.assign({}, state, {
          isLoading: false,
          data: null,
          hasLoadSuccess: false,
          hasLoadError: true,
          error: action.error.message,
        });

      case updateAsyncActionType.start:
        return Object.assign({}, state, {
          isUpdating: true,
          hasUpdateSuccess: false,
        });
      case updateAsyncActionType.success:
        return Object.assign({}, state, {
          isUpdating: false,
          data: action.payload,
          hasUpdateSuccess: true,
          hasUpdateError: false,
        });
      case updateAsyncActionType.failure:
        return Object.assign({}, state, {
          isUpdating: false,
          hasUpdateSuccess: false,
          hasUpdateError: true,
          error: action.error.message,
        });

      default:
        return state;
    }
  };
}

const loadAccountingType = 'LOAD_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';

export const loadAccounting = ReduxHelpers.generateAsyncActionCreator(
  loadAccountingType,
  LoadAccountingAPI,
);

export const updateAccounting = ReduxHelpers.generateAsyncActionCreator(
  updateAccountingType,
  UpdateAccountingAPI,
);

export const officeAccountingReducer = generatePanelReducer(
  loadAccountingType,
  updateAccountingType,
);
