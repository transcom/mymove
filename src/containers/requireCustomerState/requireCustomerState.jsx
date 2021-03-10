import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { push } from 'connected-react-router';

import { selectServiceMemberProfileState } from 'store/entities/selectors';
import { findNextServiceMemberStep } from 'utils/customer';

const requireCustomerState = (Component, requiredState) => {
  const RequireCustomerState = (props) => {
    const dispatch = useDispatch();
    const currentProfileState = useSelector(selectServiceMemberProfileState);

    useEffect(() => {
      // Only verify state on mount (once)
      if (requiredState !== currentProfileState) {
        const redirectTo = findNextServiceMemberStep(currentProfileState);
        dispatch(push(redirectTo));
      }
    }, [currentProfileState, dispatch]);

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Component {...props} />;
  };

  return RequireCustomerState;
};

export default requireCustomerState;
