import React, { lazy, Suspense } from 'react';
import propTypes from 'prop-types';
import { Switch } from 'react-router-dom';

import { MoveTabNavWithRouter } from '../shared/Header/Office';
import { RenderWithOrWithoutHeader } from '../scenes/Office/index';

import PrivateRoute from 'shared/User/PrivateRoute';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { roleTypes } from 'constants/userRoles';

const MoveDetails = lazy(() => import('./Office/MoveDetails/MoveDetails'));
const TOOMoveTaskOrder = lazy(() => import('./TOO/moveTaskOrder'));
const PaymentRequestShow = lazy(() => import('../scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('./moveHistory'));

const TXOMoveInfo = ({ too, tio, tag }) => {
  /* eslint-disable react/jsx-props-no-spreading */
  return (
    <>
      <MoveTabNavWithRouter />
      <Switch>
        {too && (
          <PrivateRoute
            path="/moves/:moveOrderId/details"
            exact
            component={(componentProps) => (
              <Suspense fallback={<LoadingPlaceholder />}>
                <RenderWithOrWithoutHeader component={MoveDetails} tag={tag} {...componentProps} />
              </Suspense>
            )}
            requiredRoles={[roleTypes.TOO]}
            hideSwitcher
          />
        )}
        {too && (
          <PrivateRoute
            path="/moves/:moveTaskOrderId/mto"
            exact
            component={(componentProps) => (
              <Suspense fallback={<LoadingPlaceholder />}>
                <RenderWithOrWithoutHeader component={TOOMoveTaskOrder} tag={tag} {...componentProps} />
              </Suspense>
            )}
            requiredRoles={[roleTypes.TOO]}
            hideSwitcher
          />
        )}
        {tio && (
          <PrivateRoute
            path="/moves/:id/payment-requests"
            exact
            component={(componentProps) => (
              <Suspense fallback={<LoadingPlaceholder />}>
                <RenderWithOrWithoutHeader component={PaymentRequestShow} tag={tag} {...componentProps} />
              </Suspense>
            )}
            requiredRoles={[roleTypes.TIO]}
            hideSwitcher
          />
        )}
        {tio && (
          <PrivateRoute
            path="/moves/:moveOrderId/history"
            exact
            component={(componentProps) => (
              <Suspense fallback={<LoadingPlaceholder />}>
                <RenderWithOrWithoutHeader component={MoveHistory} tag={tag} {...componentProps} />
              </Suspense>
            )}
            requiredRoles={[roleTypes.TIO]}
            hideSwitcher
          />
        )}
      </Switch>
    </>
  );
};

TXOMoveInfo.propTypes = {
  too: propTypes.bool,
  tio: propTypes.bool,
  tag: propTypes.string.isRequired,
};

TXOMoveInfo.defaultProps = {
  too: false,
  tio: false,
};

export default TXOMoveInfo;
