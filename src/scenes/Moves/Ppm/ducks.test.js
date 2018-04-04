import { CREATE_OR_UPDATE_PPM, GET_PPM, ppmReducer } from './ducks';

describe('Ppm Reducer', () => {
  describe('CREATE_OR_UPDATE_PPM', () => {
    it('Should handle CREATE_OR_UPDATE_PPM_SUCCESS', () => {
      const initialState = { pendingValue: '' };

      const newState = ppmReducer(initialState, {
        type: CREATE_OR_UPDATE_PPM.success,
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
        type: CREATE_OR_UPDATE_PPM.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        pendingValue: '',
        currentPpm: null,
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
      });
    });
  });
  describe('GET_PPM', () => {
    it('Should handle GET_PPM_SUCCESS', () => {
      const initialState = { pendingValue: '' };

      const newState = ppmReducer(initialState, {
        type: GET_PPM.success,
        item: 'Successful ppm!',
      });

      expect(newState).toEqual({
        pendingValue: '',
        currentPpm: 'Successful ppm!',
        hasSubmitError: false,
        hasSubmitSuccess: true,
      });
    });

    it('Should handle GET_PPM_FAILURE', () => {
      const initialState = { pendingValue: '' };

      const newState = ppmReducer(initialState, {
        type: GET_PPM.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        pendingValue: '',
        currentPpm: null,
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
      });
    });
  });
});
