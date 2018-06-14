import { APPROVE_REIMBURSEMENT, officeReducer } from './ducks';

function createAdvance(status) {
  return {
    id: '5d2b211c-0dca-4e6a-b0c0-a987d4fed5fb',
    method_of_receipt: 'MIL_PAY',
    requested_amount: 1000,
    status,
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
});
