import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { push } from 'connected-react-router';

import { selectServiceMemberFromLoggedInUser, selectServiceMemberProfileState } from 'store/entities/selectors';
import { findNextServiceMemberStep } from 'utils/customer';

const requireCustomerState = (Component, requiredState) => {
  const RequireCustomerState = (props) => {
    const dispatch = useDispatch();
    const serviceMember = useSelector(selectServiceMemberFromLoggedInUser);
    const currentProfileState = useSelector(selectServiceMemberProfileState);

    useEffect(() => {
      // Only verify state on mount (once)
      if (requiredState !== currentProfileState) {
        const redirectTo = findNextServiceMemberStep(serviceMember.id, currentProfileState);
        dispatch(push(redirectTo));
      }
    }, [currentProfileState, dispatch, serviceMember]);

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Component {...props} />;
  };

  return RequireCustomerState;
};

export default requireCustomerState;
