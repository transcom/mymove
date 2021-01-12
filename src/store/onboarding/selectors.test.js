import { selectConusStatus } from './selectors';

describe('selectConusStatus', () => {
  it('returns the conusStatus value', () => {
    const testState = {
      onboarding: {
        conusStatus: 'CONUS',
      },
    };

    expect(selectConusStatus(testState)).toEqual(testState.onboarding.conusStatus);
  });
});
