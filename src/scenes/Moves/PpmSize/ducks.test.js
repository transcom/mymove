import { CREATE_PPM_SUCCESS, CREATE_PPM_FAILURE, ppmReducer } from './ducks';

describe('Feedback Reducer', () => {
  it('Should handle CREATE_PPM_SUCCESS', () => {
    const initialState = { pendingValue: '' };

    const newState = ppmReducer(initialState, {
      type: CREATE_PPM_SUCCESS,
      item: 'Successful ppm!',
    });

    expect(newState).toEqual({
      pendingValue: '',
      currentPpm: 'Successful ppm!',
      hasSubmitError: false,
      hasSubmitSuccess: true,
    });
  });

  it('Should handle CREATE_PPM_FAILURE', () => {
    const initialState = { pendingValue: '' };

    const newState = ppmReducer(initialState, {
      type: CREATE_PPM_FAILURE,
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      pendingValue: '',
      currentPpm: {},
      hasSubmitError: true,
      hasSubmitSuccess: false,
    });
  });
});
