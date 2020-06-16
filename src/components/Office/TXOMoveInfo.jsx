import React, { lazy } from 'react';
import { Switch } from 'react-router-dom';
import PrivateRoute from 'shared/User/PrivateRoute';

import { MoveTabNavWithRouter } from 'shared/Header/Office';

const MoveDetails = lazy(() => import('pages/TOO/moveDetails'));
const TOOMoveTaskOrder = lazy(() => import('pages/TOO/moveTaskOrder'));
const PaymentRequestShow = lazy(() => import('../../scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('pages/TIO/moveHistory'));

const TXOMoveInfo = () => {
  return (
    <>
      <MoveTabNavWithRouter />
      <Switch>
        <PrivateRoute path="/moves/:moveOrderId/details" exact component={MoveDetails} />
        <PrivateRoute path="/moves/:moveTaskOrderId/mto" exact component={TOOMoveTaskOrder} />
        <PrivateRoute path="/moves/:id/payment-requests" exact component={PaymentRequestShow} />
        <PrivateRoute path="/moves/:moveOrderId/history" exact component={MoveHistory} />
      </Switch>
    </>
  );
};

export default TXOMoveInfo;
