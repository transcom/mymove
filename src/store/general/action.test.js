import {
  setMoveId,
  SET_MOVE_ID,
  SET_CAN_ADD_ORDERS,
  setCanAddOrders,
  SET_SHOW_LOADING_SPINNER,
  setShowLoadingSpinner,
} from './actions';

describe('GeneralState actions', () => {
  it('setMoveId returns the expected action', () => {
    const expectedAction = {
      type: SET_MOVE_ID,
      payload: 'test',
    };

    expect(setMoveId('test')).toEqual(expectedAction);
  });

  it('canAddOrders returns the expected action', () => {
    const expectedAction = {
      type: SET_CAN_ADD_ORDERS,
      payload: true,
    };

    expect(setCanAddOrders(true)).toEqual(expectedAction);
  });

  it('setShowLoadingSpinner returns the expected action', () => {
    const expectedAction = {
      type: SET_SHOW_LOADING_SPINNER,
      showSpinner: true,
      loadingSpinnerMessage: 'test message',
    };

    expect(setShowLoadingSpinner(true, 'test message')).toEqual(expectedAction);
  });
});
