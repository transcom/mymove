import { LOA_VALIDATION_ACTIONS, initialState, reducer } from './loaValidation';

describe('reducers/loaValidation', () => {
  it('creates an initialState', () => {
    expect(initialState()).toEqual({
      HHG: {
        isValid: false,
        longLineOfAccounting: '',
        loa: null,
      },
      NTS: {
        isValid: false,
        longLineOfAccounting: '',
        loa: null,
      },
    });
  });

  it('accepts validation actions', () => {
    const state = initialState();
    expect(
      reducer(state, {
        type: LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        payload: {
          isValid: true,
          longLineOfAccounting: '1234',
          loa: '1234',
          loaType: 'HHG',
        },
      }),
    ).toEqual({
      NTS: { isValid: false, loa: null, longLineOfAccounting: '' },
      HHG: { isValid: true, loa: '1234', longLineOfAccounting: '1234' },
    });
  });
});
