import {
  setConusStatus,
  SET_CONUS_STATUS,
  initOnboarding,
  INIT_ONBOARDING,
  initOnboardingFailed,
  INIT_ONBOARDING_FAILED,
  initOnboardingComplete,
  INIT_ONBOARDING_COMPLETE,
  fetchCustomerData,
  FETCH_CUSTOMER_DATA,
  setPPMEstimateError,
  SET_PPM_ESTIMATE_ERROR,
} from './actions';

describe('Onboarding actions', () => {
  it('setConusStatus returns the expected action', () => {
    const expectedAction = {
      type: SET_CONUS_STATUS,
      moveType: 'CONUS',
    };

    expect(setConusStatus('CONUS')).toEqual(expectedAction);
  });

  it('initOnboarding returns the expected action', () => {
    const expectedAction = {
      type: INIT_ONBOARDING,
    };

    expect(initOnboarding()).toEqual(expectedAction);
  });

  it('initOnboardingFailed returns the expected action', () => {
    const expectedAction = {
      type: INIT_ONBOARDING_FAILED,
      error: 'Test Error',
    };

    expect(initOnboardingFailed('Test Error')).toEqual(expectedAction);
  });

  it('initOnboardingComplete returns the expected action', () => {
    const expectedAction = {
      type: INIT_ONBOARDING_COMPLETE,
    };

    expect(initOnboardingComplete()).toEqual(expectedAction);
  });

  it('fetchCustomerData returns the expected action', () => {
    const expectedAction = {
      type: FETCH_CUSTOMER_DATA,
    };

    expect(fetchCustomerData()).toEqual(expectedAction);
  });

  it('setPPMEstimateError returns the expected action', () => {
    const expectedAction = {
      type: SET_PPM_ESTIMATE_ERROR,
      error: { message: 'This is a test error' },
    };

    expect(setPPMEstimateError({ message: 'This is a test error' })).toEqual(expectedAction);
  });
});
