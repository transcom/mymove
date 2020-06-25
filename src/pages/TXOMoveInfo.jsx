import React, { lazy, Suspense } from 'react';
import propTypes from 'prop-types';
import { withRouter, matchPath } from 'react-router';
import { NavLink, Switch } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import { RenderWithOrWithoutHeader } from '../scenes/Office/index';

import PrivateRoute from 'shared/User/PrivateRoute';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { roleTypes } from 'constants/userRoles';
import TabNav from 'components/TabNav';
import { MatchShape, LocationShape } from 'types/router';

const MoveDetails = lazy(() => import('./Office/MoveDetails/MoveDetails'));
const TOOMoveTaskOrder = lazy(() => import('./TOO/moveTaskOrder'));
const PaymentRequestShow = lazy(() => import('../scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('./moveHistory'));
const MoveOrders = lazy(() => import('pages/Office/MoveOrders/MoveOrders'));

const TXOMoveInfo = ({ too, tio, tag, match, location }) => {
  const { moveOrderId } = match.params;

  let matchDetails = false;
  let matchMTO = false;
  let matchPaymentRequests = false;
  let matchHistory = false;

  // Used to set aria-expanded attribute to selected tab
  if (matchPath(location.pathname, { path: '/moves/:id/details', exact: true })) {
    matchDetails = true;
  } else if (matchPath(location.pathname, { path: '/moves/:id/payment-requests', exact: true })) {
    matchPaymentRequests = true;
  } else if (matchPath(location.pathname, { path: '/moves/:id/mto', exact: true })) {
    matchMTO = true;
  } else if (matchPath(location.pathname, { path: '/moves/:id/history', exact: true })) {
    matchHistory = true;
  }

  /* eslint-disable react/jsx-props-no-spreading */
  return (
    <>
      <header className="nav-header">
        <div className="grid-container-desktop-lg">
          <TabNav
            items={[
              <NavLink
                exact
                activeClassName="usa-current"
                className=""
                to={`/moves/${moveOrderId}/details`}
                role="tab"
                aria-expanded={matchDetails}
              >
                <span className="tab-title">Move details</span>
                <Tag>2</Tag>
              </NavLink>,
              <NavLink
                exact
                activeClassName="usa-current"
                className=""
                to={`/moves/${moveOrderId}/mto`}
                role="tab"
                aria-expanded={matchMTO}
              >
                <span className="tab-title">Move task order</span>
              </NavLink>,
              <NavLink
                exact
                activeClassName="usa-current"
                className=""
                to={`/moves/${moveOrderId}/payment-requests`}
                role="tab"
                aria-expanded={matchPaymentRequests}
              >
                <span className="tab-title">Payment requests</span>
              </NavLink>,
              <NavLink
                exact
                activeClassName="usa-current"
                className=""
                to={`/moves/${moveOrderId}/history`}
                role="tab"
                aria-expanded={matchHistory}
              >
                <span className="tab-title">History</span>
              </NavLink>,
            ]}
          />
        </div>
      </header>
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
        <PrivateRoute path="/moves/:id/neworders" component={MoveOrders} /> {/* TODO fix this URL */}
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
  location: LocationShape.isRequired,
  match: MatchShape.isRequired,
};

TXOMoveInfo.defaultProps = {
  too: false,
  tio: false,
};

export default withRouter(TXOMoveInfo);
