export const INIT_ONBOARDING = 'INIT_ONBOARDING';
export const INIT_ONBOARDING_FAILED = 'INIT_ONBOARDING_FAILED';
export const INIT_ONBOARDING_COMPLETE = 'INIT_ONBOARDING_COMPLETE';

export const FETCH_CUSTOMER_DATA = 'FETCH_CUSTOMER_DATA';

export const SET_CONUS_STATUS = 'SET_CONUS_STATUS';

export const initOnboarding = () => ({
  type: INIT_ONBOARDING,
});

export const initOnboardingFailed = (error) => ({
  type: INIT_ONBOARDING_FAILED,
  error,
});

export const initOnboardingComplete = () => ({
  type: INIT_ONBOARDING_COMPLETE,
});

export const fetchCustomerData = () => ({
  type: FETCH_CUSTOMER_DATA,
});

export const setConusStatus = (moveType) => ({
  type: SET_CONUS_STATUS,
  moveType,
});
