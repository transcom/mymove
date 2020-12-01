import { REMOVE_SUCCESS_BANNER, SHOW_SUCCESS_BANNER, signedCertificationReducer } from './ducks';

describe('Feedback Reducer', () => {
  it('Should handle REMOVE_SUCCESS_BANNER', () => {
    const initialState = { moveSubmitSuccess: true };

    const newState = signedCertificationReducer(initialState, {
      type: REMOVE_SUCCESS_BANNER,
    });

    expect(newState).toEqual({
      moveSubmitSuccess: false,
    });
  });

  it('Should handle SHOW_SUCCESS_BANNER', () => {
    const initialState = { moveSubmitSuccess: false };

    const newState = signedCertificationReducer(initialState, {
      type: SHOW_SUCCESS_BANNER,
    });

    expect(newState).toEqual({
      moveSubmitSuccess: true,
    });
  });
});
