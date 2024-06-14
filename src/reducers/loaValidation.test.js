import { LOA_VALIDATION_ACTIONS, initialState, reducer } from './loaValidation';

describe('reducers/loaValidation', () => {
  it('creates an initialState', () => {
    expect(initialState()).toEqual({
      isValid: false,
      loa: null,
    });
  });

  it('accepts validation actions', () => {
    const state = initialState();
    expect(
      reducer(state, {
        type: LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        payload: {
          isValid: true,
          loa: {
            id: '1234',
            loaSysId: '5678',
          },
        },
      }),
    ).toEqual({
      isValid: true,
      loa: {
        id: '1234',
        loaSysId: '5678',
      },
    });
  });

  it('handles invalid LOA validation', () => {
    const state = initialState();
    expect(
      reducer(state, {
        type: LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        payload: {
          isValid: false,
          loa: null,
        },
      }),
    ).toEqual({
      isValid: false,
      loa: null,
    });
  });
});
