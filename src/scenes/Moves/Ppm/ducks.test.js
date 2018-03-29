import {
  CREATE_OR_UPDATE_PPM_SUCCESS,
  CREATE_OR_UPDATE_PPM_FAILURE,
  ppmReducer,
} from './ducks';

describe('Ppm Reducer', () => {
  it('Should handle CREATE_OR_UPDATE_PPM_SUCCESS', () => {
    const initialState = { pendingValue: '' };

    const newState = ppmReducer(initialState, {
      type: CREATE_OR_UPDATE_PPM_SUCCESS,
      item: 'Successful ppm!',
    });

    expect(newState).toEqual({
      pendingValue: '',
      pendingPpmSize: null,
      currentPpm: 'Successful ppm!',
      hasSubmitError: false,
      hasSubmitSuccess: true,
    });
  });

  it('Should handle CREATE_OR_UPDATE_PPM_FAILURE', () => {
    const initialState = { pendingValue: '' };

    const newState = ppmReducer(initialState, {
      type: CREATE_OR_UPDATE_PPM_FAILURE,
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      pendingValue: '',
      currentPpm: null,
      hasSubmitError: true,
      hasSubmitSuccess: false,
    });
  });
});
