import React, { Suspense, lazy } from 'react';
import { NavLink, Switch, useParams } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import 'styles/office.scss';
import PrivateRoute from 'containers/PrivateRoute';
import { roleTypes } from 'constants/userRoles';
import TabNav from 'components/TabNav';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const MoveDetails = lazy(() => import('pages/Office/MoveDetails/MoveDetails'));
const TOOMoveTaskOrder = lazy(() => import('pages/TOO/moveTaskOrder'));
const PaymentRequestShow = lazy(() => import('scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const MoveOrders = lazy(() => import('pages/Office/MoveOrders/MoveOrders'));

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
          <PrivateRoute path="/moves/:id/orders" exact component={MoveOrders} />
          <PrivateRoute
            path="/moves/:moveTaskOrderId/mto"
            exact
            component={TOOMoveTaskOrder}
            requiredRoles={[roleTypes.TOO]}
          />
          <PrivateRoute
            path="/moves/:id/payment-requests"
            exact
            component={PaymentRequestShow}
            requiredRoles={[roleTypes.TIO]}
          />
          <PrivateRoute
            path="/moves/:moveOrderId/history"
            exact
            component={MoveHistory}
            requiredRoles={[roleTypes.TIO]}
          />
        </Switch>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
