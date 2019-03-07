import { CREATE_OR_UPDATE_MOVE, GET_MOVE, SUBMIT_FOR_APPROVAL, moveReducer } from './ducks';
import loggedInUserPayload, { emptyPayload } from 'shared/User/sampleLoggedInUserPayload';

const expectedMove = {
  id: '593cc830-1a3e-44b3-ba5a-8809f02dfa7d',
  locator: 'WUMGLQ',
  orders_id: '51953e97-25a7-430c-ba6d-3bd980a38b71',
  selected_move_type: 'PPM',
  status: 'DRAFT',
};
const latestMove = {
  id: '593cc830-1a3e-44b3-ba5a-8809f02d000',
  locator: 'BLABLA',
  orders_id: '51953e97-25a7-430c-ba6d-3bd980a38b00',
  selected_move_type: 'PPM',
  status: 'CANCELED',
};
const movePayload = {
  created_at: '2018-05-25T21:36:10.235Z',
  id: '593cc830-1a3e-44b3-ba5a-8809f02dfa7d',
  locator: 'WUMGLQ',
  orders_id: '51953e97-25a7-430c-ba6d-3bd980a38b71',
  personally_procured_moves: [
    {
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
  ],
  selected_move_type: 'PPM',
  status: 'DRAFT',
};
describe('move Reducer', () => {
  describe('GET_LOGGED_IN_USER', () => {
    it('Should handle GET_LOGGED_IN_USER.success', () => {
      const initialState = {};

      const newState = moveReducer(initialState, loggedInUserPayload);

      expect(newState).toEqual({
        currentMove: { ...expectedMove },
        latestMove: { ...latestMove },
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    });
    it('Should handle GET_LOGGED_IN_USER.success with empty payload', () => {
      const initialState = {};

      const newState = moveReducer(initialState, emptyPayload);

      expect(newState).toEqual({
        currentMove: null,
        latestMove: null,
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    });
  });

  describe('CREATE_OR_UPDATE_MOVE', () => {
    it('Should handle CREATE_OR_UPDATE_MOVE_SUCCESS', () => {
      const initialState = {};

      const newState = moveReducer(initialState, {
        type: CREATE_OR_UPDATE_MOVE.success,
        payload: movePayload,
      });

      expect(newState).toEqual({
        currentMove: { ...expectedMove },
        latestMove: null,
        error: null,
        hasSubmitError: false,
        hasSubmitSuccess: true,
        pendingMoveType: null,
      });
    });

    it('Should handle CREATE_OR_UPDATE_MOVE_FAILURE', () => {
      const initialState = {};

      const newState = moveReducer(initialState, {
        type: CREATE_OR_UPDATE_MOVE.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentMove: {},
        latestMove: null,
        error: 'No bueno.',
        hasSubmitError: true,
        hasSubmitSuccess: false,
      });
    });
  });

  describe('GET_MOVE', () => {
    it('Should handle GET_MOVE_SUCCESS', () => {
      const initialState = {};
      const newState = moveReducer(initialState, {
        type: GET_MOVE.success,
        payload: movePayload,
      });

      expect(newState).toEqual({
        currentMove: { ...expectedMove },
        latestMove: null,
        error: null,
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    });

    it('Should handle GET_MOVE_FAILURE', () => {
      const initialState = {};

      const newState = moveReducer(initialState, {
        type: GET_MOVE.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentMove: {},
        latestMove: null,
        hasLoadError: true,
        hasLoadSuccess: false,
        error: 'No bueno.',
      });
    });
  });

  describe('SUBMIT_FOR_APPROVAL', () => {
    it('Should handle SUCCESS', () => {
      const initialState = {};
      const newState = moveReducer(initialState, {
        type: SUBMIT_FOR_APPROVAL.success,
        payload: { ...movePayload, status: 'APPROVED' },
      });

      expect(newState).toEqual({
        currentMove: { ...expectedMove, status: 'APPROVED' },
        submittedForApproval: true,
      });
    });

    it('Should handle FAILURE', () => {
      const initialState = {};

      const newState = moveReducer(initialState, {
        type: SUBMIT_FOR_APPROVAL.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        submittedForApproval: false,
        error: 'No bueno.',
      });
    });
  });
});
