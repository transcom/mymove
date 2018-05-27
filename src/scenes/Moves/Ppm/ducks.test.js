import { CREATE_OR_UPDATE_PPM, GET_PPM, ppmReducer } from './ducks';
import { GET_LOGGED_IN_USER } from 'shared/User/ducks';
import loggedInUserPayload, {
  emptyPayload,
} from 'shared/user/sampleLoggedInUserPayload';
describe('Ppm Reducer', () => {
  const samplePpm = { id: 'UUID', name: 'foo' };
  describe('GET_LOGGED_IN_USER', () => {
    it('Should handle GET_LOGGED_IN_USER.success', () => {
      const initialState = {
        pendingValue: '',
        hasSubmitError: false,
        hasSubmitSuccess: true,
      };

      const newState = ppmReducer(initialState, loggedInUserPayload);

      expect(newState).toEqual({
        currentPpm: {
          destination_postal_code: '76127',
          estimated_incentive: '$14954.09 - 16528.21',
          has_additional_postal_code: false,
          has_requested_advance: false,
          has_sit: false,
          id: 'cd67c9e4-ef59-45e5-94bc-767aaafe559e',
          pickup_postal_code: '80913',
          planned_move_date: '2018-06-28',
          size: 'L',
          status: 'DRAFT',
          weight_estimate: 9000,
        },
        hasSubmitError: false,
        hasSubmitSuccess: true,
        hasLoadError: false,
        hasLoadSuccess: true,
        incentive: '$14954.09 - 16528.21',
        pendingPpmSize: 'L',
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
        incentive: null,
        pendingPpmSize: null,
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
        pendingPpmSize: null,
        pendingPpmWeight: null,
        currentPpm: samplePpm,
        incentive: null,
        sitReimbursement: null,
        hasSubmitError: false,
        hasSubmitSuccess: true,
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
        incentive: null,
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
});
