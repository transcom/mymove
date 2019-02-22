import {
  CREATE_OR_UPDATE_PPM,
  GET_PPM,
  GET_SIT_ESTIMATE,
  GET_PPM_ESTIMATE,
  ppmReducer,
  getMaxAdvance,
  getEstimatedRemainingWeight,
  getActualRemainingWeight,
} from './ducks';
import loggedInUserPayload, { emptyPayload } from 'shared/User/sampleLoggedInUserPayload';
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
          incentive_estimate_min: 1495409,
          incentive_estimate_max: 1652821,
          has_additional_postal_code: false,
          has_requested_advance: false,
          has_sit: false,
          id: 'cd67c9e4-ef59-45e5-94bc-767aaafe559e',
          pickup_postal_code: '80913',
          original_move_date: '2018-06-28',
          size: 'L',
          status: 'DRAFT',
          weight_estimate: 9000,
        },
        hasSubmitError: false,
        hasSubmitSuccess: true,
        hasLoadError: false,
        hasLoadSuccess: true,
        incentive_estimate_min: 1495409,
        incentive_estimate_max: 1652821,
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
        incentive_estimate_min: null,
        incentive_estimate_max: null,
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
        incentive_estimate_min: null,
        incentive_estimate_max: null,
        sitReimbursement: null,
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

  describe('CLEAR_SIT_ESTIMATE', () => {
    it('Should handle SUCCESS', () => {
      const initialState = {};
      const newState = ppmReducer(initialState, {
        type: 'CLEAR_SIT_ESTIMATE',
      });

      expect(newState).toEqual({
        sitReimbursement: null,
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        rateEngineError: null,
      });
    });
  });
  describe('GET_SIT_ESTIMATE', () => {
    it('Should handle SUCCESS', () => {
      const initialState = {};
      const newState = ppmReducer(initialState, {
        type: GET_SIT_ESTIMATE.success,
        payload: { estimate: 21505 },
      });

      expect(newState).toEqual({
        sitReimbursement: '$215.05',
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        rateEngineError: null,
      });
    });

    it('Should handle FAILURE', () => {
      const initialState = { pendingValue: '' };

      const newState = ppmReducer(initialState, {
        type: GET_SIT_ESTIMATE.failure,
        error: 'No bueno.',
      });
      // using special error here so it is not caught by WizardPage handling
      expect(newState).toEqual({
        hasEstimateError: true,
        hasEstimateInProgress: false,
        hasEstimateSuccess: false,
        pendingValue: '',
        rateEngineError: 'No bueno.',
        sitReimbursement: null,
      });
    });
  });

  describe('GET_PPM_ESTIMATE', () => {
    it('Should handle SUCCESS', () => {
      const initialState = {};
      const newState = ppmReducer(initialState, {
        type: GET_PPM_ESTIMATE.success,
        payload: { range_min: 21505, range_max: 44403 },
      });

      expect(newState).toEqual({
        incentive_estimate_min: 21505,
        incentive_estimate_max: 44403,
        hasEstimateSuccess: true,
        hasEstimateError: false,
        hasEstimateInProgress: false,
        rateEngineError: null,
      });
    });

    it('Should handle FAILURE', () => {
      const initialState = { pendingValue: '' };

      const newState = ppmReducer(initialState, {
        type: GET_PPM_ESTIMATE.failure,
        error: 'No bueno.',
      });
      // using special error here so it is not caught by WizardPage handling
      expect(newState).toEqual({
        hasEstimateError: true,
        hasEstimateInProgress: false,
        hasEstimateSuccess: false,
        pendingValue: '',
        rateEngineError: 'No bueno.',
        incentive_estimate_min: null,
        incentive_estimate_max: null,
        error: null,
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

  describe('getEstimatedRemainingWeight', () => {
    describe('when there is an estimated remaining weight from a pre move survey', () => {
      it('should return the proper estimated remaining weight based on the pre move survey', () => {
        const state = {
          moves: {
            currentMove: {
              selected_move_type: 'HHG_PPM',
            },
          },
          orders: {
            currentOrders: {
              id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
              has_dependents: true,
              spouse_has_pro_gear: true,
            },
          },
          serviceMember: {
            currentServiceMember: {
              rank: 'E_1',
            },
          },
          ui: {
            currentShipmentID: '0194fb44-3762-4d1c-b58a-f6daf984813b',
          },
          entities: {
            shipments: {
              '0194fb44-3762-4d1c-b58a-f6daf984813b': {
                gross_weight: 5000,
                id: '0194fb44-3762-4d1c-b58a-f6daf984813b',
                market: 'dHHG',
                move_id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                pm_survey_weight_estimate: 5000,
                progear_weight_estimate: 225,
                spouse_progear_weight_estimate: 312,
                status: 'AWARDED',
                tare_weight: 1500,
                weight_estimate: 2000,
              },
            },
          },
        };
        expect(getEstimatedRemainingWeight(state)).toEqual(5500);
      });
    });

    describe('when there is an estimated remaining weight from a service member entered weight', () => {
      it('should return the proper estimated remaining weight based on the entered values', () => {
        const state = {
          moves: {
            currentMove: {
              selected_move_type: 'HHG_PPM',
            },
          },
          orders: {
            currentOrders: {
              id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
              has_dependents: true,
              spouse_has_pro_gear: true,
            },
          },
          serviceMember: {
            currentServiceMember: {
              rank: 'E_1',
            },
          },
          ui: {
            currentShipmentID: '0194fb44-3762-4d1c-b58a-f6daf984813b',
          },
          entities: {
            shipments: {
              '0194fb44-3762-4d1c-b58a-f6daf984813b': {
                gross_weight: 5000,
                id: '0194fb44-3762-4d1c-b58a-f6daf984813b',
                market: 'dHHG',
                move_id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                progear_weight_estimate: 225,
                spouse_progear_weight_estimate: 312,
                status: 'AWARDED',
                tare_weight: 1500,
                weight_estimate: 2000,
              },
            },
          },
        };
        expect(getEstimatedRemainingWeight(state)).toEqual(8500);
      });
    });

    describe('getActualRemainingWeight', () => {
      describe('when there is an actual weight', () => {
        it('should return the proper actual weight', () => {
          const state = {
            moves: {
              currentMove: {
                selected_move_type: 'HHG_PPM',
              },
            },
            orders: {
              currentOrders: {
                id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
                has_dependents: true,
                spouse_has_pro_gear: true,
              },
            },
            serviceMember: {
              currentServiceMember: {
                rank: 'E_1',
              },
            },
            ui: {
              currentShipmentID: '0194fb44-3762-4d1c-b58a-f6daf984813b',
            },
            entities: {
              shipments: {
                '0194fb44-3762-4d1c-b58a-f6daf984813b': {
                  gross_weight: 5000,
                  id: '0194fb44-3762-4d1c-b58a-f6daf984813b',
                  market: 'dHHG',
                  move_id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                  pm_survey_weight_estimate: 5000,
                  progear_weight_estimate: 225,
                  spouse_progear_weight_estimate: 312,
                  status: 'AWARDED',
                  tare_weight: 1500,
                  weight_estimate: 2000,
                },
              },
            },
          };
          expect(getActualRemainingWeight(state)).toEqual(7000);
        });
      });
    });
  });
});
