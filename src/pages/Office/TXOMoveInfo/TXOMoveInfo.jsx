import React, { Suspense, lazy } from 'react';
import { NavLink, Switch, useParams, Redirect } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import 'styles/office.scss';
import PrivateRoute from 'containers/PrivateRoute';
import { roleTypes } from 'constants/userRoles';
import TabNav from 'components/TabNav';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const MoveDetails = lazy(() => import('pages/Office/MoveDetails/MoveDetails'));
const TOOMoveTaskOrder = lazy(() => import('pages/TOO/moveTaskOrder'));
const MoveOrders = lazy(() => import('pages/Office/MoveOrders/MoveOrders'));
const PaymentRequestIndex = lazy(() => import('scenes/Office/TIO/paymentRequestIndex'));
const PaymentRequestShow = lazy(() => import('scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));

const TXOMoveInfo = () => {
  const { moveOrderId } = useParams();

  return (
    <>
      <header className="nav-header">
        <div className="grid-container-desktop-lg">
          <TabNav
            items={[
              <NavLink exact activeClassName="usa-current" to={`/moves/${moveOrderId}/details`} role="tab">
                <span className="tab-title">Move details</span>
                <Tag>2</Tag>
              </NavLink>,
              <NavLink exact activeClassName="usa-current" to={`/moves/${moveOrderId}/mto`} role="tab">
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
      <Suspense fallback={<LoadingPlaceholder />}>
        <Switch>
          <PrivateRoute
            path="/moves/:moveOrderId/details"
            exact
            component={MoveDetails}
            requiredRoles={[roleTypes.TOO]}
          />
          <PrivateRoute path="/moves/:moveOrderId/orders" exact component={MoveOrders} />

          {/* TODO - the nav to this url is passing moveOrderId instead of moveTaskOrderId */}
          <PrivateRoute
            path="/moves/:moveTaskOrderId/mto"
            exact
            component={TOOMoveTaskOrder}
            requiredRoles={[roleTypes.TOO]}
          />

          <PrivateRoute
            path="/moves/:moveOrderId/payment-requests/:id"
            exact
            component={PaymentRequestShow}
            requiredRoles={[roleTypes.TIO]}
          />
          <PrivateRoute
            path="/moves/:moveOrderId/payment-requests"
            exact
            component={PaymentRequestIndex}
            requiredRoles={[roleTypes.TIO]}
          />

          <PrivateRoute
            path="/moves/:moveOrderId/history"
            exact
            component={MoveHistory}
            requiredRoles={[roleTypes.TIO]}
          />

          <Redirect from="/moves/:moveOrderId" to="/moves/:moveOrderId/details" />
        </Switch>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
