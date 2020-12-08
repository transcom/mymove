import React, { Suspense, lazy } from 'react';
import { NavLink, Switch, useParams, Redirect, Route, useLocation, matchPath } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import 'styles/office.scss';
import TabNav from 'components/TabNav';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const MoveDetails = lazy(() => import('pages/Office/MoveDetails/MoveDetails'));
const MoveTaskOrder = lazy(() => import('pages/Office/MoveTaskOrder/MoveTaskOrder'));
const MoveOrders = lazy(() => import('pages/Office/MoveOrders/MoveOrders'));
const MoveAllowances = lazy(() => import('pages/Office/MoveAllowances/MoveAllowances'));
const PaymentRequestReview = lazy(() => import('pages/Office/PaymentRequestReview/PaymentRequestReview'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const MovePaymentRequests = lazy(() => import('pages/Office/MovePaymentRequests/MovePaymentRequests'));

const TXOMoveInfo = () => {
  const { moveOrderId } = useParams();
  const { pathname } = useLocation();

  const hideNav =
    matchPath(pathname, {
      path: '/moves/:moveOrderId/payment-requests/:id',
      exact: true,
    }) ||
    matchPath(pathname, {
      path: '/moves/:moveOrderId/orders',
      exact: true,
    });

  return (
    <>
      {!hideNav && (
        <header className="nav-header">
          <div className="grid-container-desktop-lg">
            <TabNav
              items={[
                <NavLink exact activeClassName="usa-current" to={`/moves/${moveOrderId}/details`} role="tab">
                  <span className="tab-title">Move details</span>
                  <Tag>2</Tag>
                </NavLink>,
                <NavLink
                  data-testid="MoveTaskOrder-Tab"
                  exact
                  activeClassName="usa-current"
                  to={`/moves/${moveOrderId}/mto`}
                  role="tab"
                >
                  <span className="tab-title">Move task order</span>
                </NavLink>,
                <NavLink exact activeClassName="usa-current" to={`/moves/${moveOrderId}/payment-requests`} role="tab">
                  <span className="tab-title">Payment requests</span>
                </NavLink>,
                <NavLink exact activeClassName="usa-current" to={`/moves/${moveOrderId}/history`} role="tab">
                  <span className="tab-title">History</span>
                </NavLink>,
              ]}
            />
          </div>
        </header>
      )}
      <Suspense fallback={<LoadingPlaceholder />}>
        <Switch>
          <Route path="/moves/:moveOrderId/details" exact>
            <MoveDetails />
          </Route>

          <Route path="/moves/:moveOrderId/orders" exact>
            <MoveOrders />
          </Route>

          <Route path="/moves/:moveOrderId/allowances" exact>
            <MoveAllowances />
          </Route>

          <Route path="/moves/:moveOrderId/mto" exact>
            <MoveTaskOrder />
          </Route>

          <Route path="/moves/:moveOrderId/payment-requests/:paymentRequestId" exact>
            <PaymentRequestReview />
          </Route>

          <Route path="/moves/:locator/payment-requests" exact>
            <MovePaymentRequests />
          </Route>

          <Route path="/moves/:moveOrderId/history" exact>
            <MoveHistory />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from="/moves/:moveOrderId" to="/moves/:moveOrderId/details" />
        </Switch>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
