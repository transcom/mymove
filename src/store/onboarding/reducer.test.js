import onboardingReducer, { initialState } from './reducer';
import { setConusStatus, setPPMEstimateError } from './actions';

describe('onboardingReducer', () => {
  it('returns the initial state by default', () => {
    expect(onboardingReducer(undefined, undefined)).toEqual(initialState);
  });

  it('handles the setConusStatus action', () => {
    expect(onboardingReducer(initialState, setConusStatus('CONUS'))).toEqual({
      ...initialState,
      conusStatus: 'CONUS',
    });
  });

  it('handles the setPPMEstimateError action', () => {
    expect(onboardingReducer(initialState, setPPMEstimateError({ message: 'This is a test error' }))).toEqual({
      ...initialState,
      ppmEstimateError: { message: 'This is a test error' },
    });
  });
});
