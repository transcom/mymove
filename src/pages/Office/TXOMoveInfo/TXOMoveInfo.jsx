import React, { Suspense, lazy } from 'react';
import { NavLink, Switch, useParams, Redirect, Route, useLocation, matchPath } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import 'styles/office.scss';
import TabNav from 'components/TabNav';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import CustomerHeader from 'components/CustomerHeader';
import { useTXOMoveInfoQueries } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const MoveDetails = lazy(() => import('pages/Office/MoveDetails/MoveDetails'));
const MoveDocumentWrapper = lazy(() => import('pages/Office/MoveDocumentWrapper/MoveDocumentWrapper'));
const MoveTaskOrder = lazy(() => import('pages/Office/MoveTaskOrder/MoveTaskOrder'));
const PaymentRequestReview = lazy(() => import('pages/Office/PaymentRequestReview/PaymentRequestReview'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const MovePaymentRequests = lazy(() => import('pages/Office/MovePaymentRequests/MovePaymentRequests'));

const TXOMoveInfo = () => {
  const [unapprovedShipmentCount, setUnapprovedShipmentCount] = React.useState(0);
  const [unapprovedServiceItemCount, setUnapprovedServiceItemCount] = React.useState(0);
  const [pendingPaymentRequestCount, setPendingPaymentRequestCount] = React.useState(0);

  const { moveCode } = useParams();
  const { pathname } = useLocation();
  const { order, customerData, isLoading, isError } = useTXOMoveInfoQueries(moveCode);
  const hideNav =
    matchPath(pathname, {
      path: '/moves/:moveCode/payment-requests/:id',
      exact: true,
    }) ||
    matchPath(pathname, {
      path: '/moves/:moveCode/orders',
      exact: true,
    }) ||
    matchPath(pathname, {
      path: '/moves/:moveCode/allowances',
      exact: true,
    });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <>
      <CustomerHeader order={order} customer={customerData} moveCode={moveCode} />
      {!hideNav && (
        <header className="nav-header">
          <div className="grid-container-desktop-lg">
            <TabNav
              items={[
                <NavLink
                  exact
                  activeClassName="usa-current"
                  to={`/moves/${moveCode}/details`}
                  role="tab"
                  data-testid="MoveDetails-Tab"
                >
                  <span className="tab-title">Move details</span>
                  {unapprovedShipmentCount > 0 && <Tag>{unapprovedShipmentCount}</Tag>}
                </NavLink>,
                <NavLink
                  data-testid="MoveTaskOrder-Tab"
                  exact
                  activeClassName="usa-current"
                  to={`/moves/${moveCode}/mto`}
                  role="tab"
                >
                  <span className="tab-title">Move task order</span>
                  {unapprovedServiceItemCount > 0 && <Tag>{unapprovedServiceItemCount}</Tag>}
                </NavLink>,
                <NavLink exact activeClassName="usa-current" to={`/moves/${moveCode}/payment-requests`} role="tab">
                  <span className="tab-title">Payment requests</span>
                  {pendingPaymentRequestCount > 0 && <Tag>{pendingPaymentRequestCount}</Tag>}
                </NavLink>,
                <NavLink exact activeClassName="usa-current" to={`/moves/${moveCode}/history`} role="tab">
                  <span className="tab-title">History</span>
                </NavLink>,
              ]}
            />
          </div>
        </header>
      )}
      <Suspense fallback={<LoadingPlaceholder />}>
        <Switch>
          <Route path="/moves/:moveCode/details" exact>
            <MoveDetails
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            />
          </Route>

          <Route path={['/moves/:moveCode/allowances', '/moves/:moveCode/orders']} exact>
            <MoveDocumentWrapper />
          </Route>

          <Route path="/moves/:moveCode/mto" exact>
            <MoveTaskOrder
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            />
          </Route>

          <Route path="/moves/:moveCode/payment-requests/:paymentRequestId" exact>
            <PaymentRequestReview />
          </Route>

          <Route path="/moves/:moveCode/payment-requests" exact>
            <MovePaymentRequests
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
              setPendingPaymentRequestCount={setPendingPaymentRequestCount}
            />
          </Route>

          <Route path="/moves/:moveCode/history" exact>
            <MoveHistory />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from="/moves/:moveCode" to="/moves/:moveCode/details" />
        </Switch>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
