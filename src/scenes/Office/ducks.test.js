import { APPROVE_REIMBURSEMENT, CANCEL_MOVE, officeReducer } from './ducks';

function createAdvance(status) {
  return {
    id: '5d2b211c-0dca-4e6a-b0c0-a987d4fed5fb',
    method_of_receipt: 'MIL_PAY',
    requested_amount: 1000,
    status,
  };
}

function createMove(status, cancelReason) {
  return {
    id: '1224ba31c-0dca-4e6a-b0c0-a987d4fed32c',
    locator: 'CAJ214',
    orders_id: '4134b2a1c-0dca-4e6a-b0c0-a987d4fe55ab',
    selected_move_type: 'PPM',
    personally_procured_moves: [],
    status,
    cancelReason,
  };
}

describe('office Reducer', () => {
  describe('APPROVE_REIMBURSEMENT', () => {
    it('handles SUCCESS', () => {
      const initialState = {
        otherState: 1,
        officePPMs: [{ id: '1234', advance: createAdvance('DRAFT') }],
      };
      const newState = officeReducer(initialState, {
        type: APPROVE_REIMBURSEMENT.success,
        payload: createAdvance('APPROVED'),
      });

      expect(newState).toEqual({
        otherState: 1,
        officePPMs: [{ id: '1234', advance: createAdvance('APPROVED') }],
        reimbursementIsApproving: false,
      });
    });

    it('handles START', () => {
      const initialState = {
        otherState: 1,
        officePPMs: [{ id: '1234', advance: createAdvance('DRAFT') }],
      };
      const newState = officeReducer(initialState, {
        type: APPROVE_REIMBURSEMENT.start,
      });

      expect(newState).toEqual({
        otherState: 1,
        officePPMs: [{ id: '1234', advance: createAdvance('DRAFT') }],
        reimbursementIsApproving: true,
      });
    });

    it('handles FAILURE', () => {
      const initialState = {
        otherState: 1,
        officePPMs: [{ id: '1234', advance: createAdvance('DRAFT') }],
      };
      const newState = officeReducer(initialState, {
        type: APPROVE_REIMBURSEMENT.failure,
        error: {
          message: 'something went wrong',
        },
      });

      expect(newState).toEqual({
        otherState: 1,
        officePPMs: [{ id: '1234', advance: createAdvance('DRAFT') }],
        reimbursementIsApproving: false,
        error: 'something went wrong',
      });
    });
  });
  describe('CANCEL_MOVE', () => {
    it('handles SUCCESS', () => {
      const initialState = {
        otherState: 1,
        officeMove: createMove('SUBMITTED', ''),
      };
      const newState = officeReducer(initialState, {
        type: CANCEL_MOVE.success,
        payload: createMove('CANCELED', 'Got tired'),
      });

      expect(newState).toEqual({
        otherState: 1,
        officeMove: createMove('CANCELED', 'Got tired'),
        moveIsCanceling: false,
        flashMessage: true,
      });
    });

    it('handles START', () => {
      const initialState = {
        otherState: 1,
        officeMove: createMove('DRAFT', ''),
      };
      const newState = officeReducer(initialState, {
        type: CANCEL_MOVE.start,
      });

      expect(newState).toEqual({
        otherState: 1,
        officeMove: createMove('DRAFT', ''),
        moveIsCanceling: true,
      });
    });

    it('handles FAILURE', () => {
      const initialState = {
        otherState: 1,
        officeMove: createMove('DRAFT', ''),
      };
      const newState = officeReducer(initialState, {
        type: CANCEL_MOVE.failure,
        error: {
          message: 'something went wrong',
        },
      });

      expect(newState).toEqual({
        otherState: 1,
        officeMove: createMove('DRAFT', ''),
        moveIsCanceling: false,
        flashMessage: false,
        error: 'something went wrong',
      });
    });
  });
});
