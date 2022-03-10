import { TAC_VALIDATION_ACTIONS, initialState, reducer } from './tacValidation';

describe('reducers/tacValidation', () => {
  it('creates an initialState', () => {
    expect(initialState()).toEqual({
      HHG: { isValid: true, tac: '' },
      NTS: { isValid: true, tac: '' },
    });
  });

  it('accepts validation actions', () => {
    const state = initialState();
    expect(
      reducer(state, {
        type: TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE,
        loaType: 'HHG',
        isValid: false,
        tac: '1234',
      }),
    ).toEqual({
      NTS: { isValid: true, tac: '' },
      HHG: { isValid: false, tac: '1234' },
    });
  });
});
