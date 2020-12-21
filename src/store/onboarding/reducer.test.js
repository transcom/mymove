import onboardingReducer, { initialState } from './reducer';
import { setConusStatus } from './actions';

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
});
