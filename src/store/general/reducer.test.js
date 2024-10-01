import generalStateReducer, { initialState } from './reducer';
import { setCanAddOrders, setMoveId } from './actions';

describe('generalStateReducer', () => {
  it('returns the initial state by default', () => {
    expect(generalStateReducer(undefined, undefined)).toEqual(initialState);
  });

  it('handles the setMoveId action', () => {
    expect(generalStateReducer(initialState, setMoveId('test'))).toEqual({
      ...initialState,
      moveId: 'test',
    });
  });

  it('handles the setCanAddOrders action', () => {
    expect(generalStateReducer(initialState, setCanAddOrders(true))).toEqual({
      ...initialState,
      canAddOrders: true,
    });
  });
});
