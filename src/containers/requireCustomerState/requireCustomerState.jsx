import React, { useEffect } from 'react';
import { useDispatch } from 'react-redux';

import { initOnboarding } from 'store/onboarding/actions';

const requireCustomerState = (Component) => {
  const RequireCustomerState = (props) => {
    const dispatch = useDispatch();

    useEffect(() => {
      // Only call initOnboarding on mount (once)
      dispatch(initOnboarding());
    }, [dispatch]);

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Component {...props} />;
  };

  return RequireCustomerState;
};

export default requireCustomerState;
