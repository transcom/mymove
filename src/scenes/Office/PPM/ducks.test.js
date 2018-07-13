import reducer, { getIncentiveActionType } from './ducks';

describe('office ppm reducer', () => {
  describe('GET_PPM_INCENTIVE', () => {
    it('handles SUCCESS', () => {
      const newState = reducer(undefined, {
        type: getIncentiveActionType.success,
        payload: { gcc: 123400, incentive_percentage: 12400 },
      });

      expect(newState).toEqual({
        isLoading: false,
        hasErrored: false,
        hasSucceeded: true,
        calculation: { gcc: 123400, incentive_percentage: 12400 },
      });
    });
    it('handles START', () => {
      const newState = reducer(undefined, {
        type: getIncentiveActionType.start,
      });
      expect(newState).toEqual({
        isLoading: true,
        hasErrored: false,
        hasSucceeded: false,
      });
    });
    it('handles FAILURE', () => {
      const newState = reducer(undefined, {
        type: getIncentiveActionType.failure,
      });
      expect(newState).toEqual({
        isLoading: false,
        hasErrored: true,
        hasSucceeded: false,
      });
    });
  });
});
