import { setConusStatus, SET_CONUS_STATUS } from './actions';

describe('Onboarding actions', () => {
  it('setConusStatus returns the expected action', () => {
    const expectedAction = {
      type: SET_CONUS_STATUS,
      moveType: 'CONUS',
    };

    expect(setConusStatus('CONUS')).toEqual(expectedAction);
  });
});
