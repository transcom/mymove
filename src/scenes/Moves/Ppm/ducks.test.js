import { CREATE_OR_UPDATE_PPM, GET_PPM, ppmReducer, getMaxAdvance } from './ducks';
import loggedInUserPayload, { emptyPayload } from 'shared/User/sampleLoggedInUserPayload';
describe('Ppm Reducer', () => {
  const samplePpm = { id: 'UUID', name: 'foo' };
  describe('GET_LOGGED_IN_USER', () => {
    it('Should handle GET_LOGGED_IN_USER_SUCCESS', () => {
      const initialState = {
        pendingValue: '',
        hasSubmitError: false,
        hasSubmitSuccess: true,
      };

      const newState = ppmReducer(initialState, loggedInUserPayload);

      expect(newState).toEqual({
        currentPpm: {
          destination_postal_code: '76127',
          incentive_estimate_min: 1495409,
          incentive_estimate_max: 1652821,
          has_additional_postal_code: false,
          has_requested_advance: false,
          has_sit: false,
          id: 'cd67c9e4-ef59-45e5-94bc-767aaafe559e',
          pickup_postal_code: '80913',
          original_move_date: '2018-06-28',
          status: 'DRAFT',
          weight_estimate: 9000,
        },
        hasSubmitError: false,
        hasSubmitSuccess: true,
        hasLoadError: false,
        hasLoadSuccess: true,
        incentive_estimate_min: 1495409,
        incentive_estimate_max: 1652821,
        pendingPpmWeight: 9000,
        pendingValue: '',
        sitReimbursement: null,
      });
    });
    it('Should handle emptyPayload', () => {
      const initialState = {
        hasSubmitError: false,
        hasSubmitSuccess: true,
      };

      const newState = ppmReducer(initialState, emptyPayload);

      expect(newState).toEqual({
        currentPpm: null,
        hasSubmitError: false,
        hasSubmitSuccess: true,
        hasLoadError: false,
        hasLoadSuccess: true,
        incentive_estimate_min: null,
        incentive_estimate_max: null,
        pendingPpmWeight: null,
        sitReimbursement: null,
      });
    });
  });
  describe('CREATE_OR_UPDATE_PPM', () => {
    it('Should handle CREATE_OR_UPDATE_PPM_SUCCESS', () => {
      const initialState = { pendingValue: '' };

      const newState = ppmReducer(initialState, {
        type: CREATE_OR_UPDATE_PPM.success,
        payload: samplePpm,
      });

      expect(newState).toEqual({
        pendingValue: '',
        pendingPpmWeight: null,
        currentPpm: samplePpm,
        sitReimbursement: null,
        id: 'UUID',
        name: 'foo',
        hasSubmitError: false,
        hasSubmitSuccess: true,
        hasSubmitInProgress: false,
      });
    });

    it('Should handle CREATE_OR_UPDATE_PPM_FAILURE', () => {
      const initialState = { pendingValue: '', currentPpm: { id: 'bad' } };

      const newState = ppmReducer(initialState, {
        type: CREATE_OR_UPDATE_PPM.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        pendingValue: '',
        currentPpm: { id: 'bad' },
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
        hasSubmitInProgress: false,
      });
    });
  });
  describe('GET_PPM', () => {
    it('Should handle GET_PPM_SUCCESS', () => {
      const initialState = { pendingValue: '' };
      const newState = ppmReducer(initialState, {
        type: GET_PPM.success,
        payload: [samplePpm],
      });

      expect(newState).toEqual({
        pendingValue: '',
        currentPpm: samplePpm,
        incentive_estimate_min: null,
        incentive_estimate_max: null,
        pendingPpmWeight: null,
        sitReimbursement: null,
        hasLoadError: false,
        hasLoadSuccess: true,
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
        hasLoadError: true,
        hasLoadSuccess: false,
        error: 'No bueno.',
      });
    });
  });

  describe('getMaxAdvance', () => {
    describe('when there is a max estimated incentive', () => {
      const state = { ppm: { incentive_estimate_max: 10000 } };
      it('should return 60% of max estimated incentive', () => {
        expect(getMaxAdvance(state)).toEqual(6000);
      });
    });
  });
  describe('when there is no max estimated incentive', () => {
    const state = {};
    it('should return 60% of max estimated incentive', () => {
      expect(getMaxAdvance(state)).toEqual(20000000);
    });
  });
});
