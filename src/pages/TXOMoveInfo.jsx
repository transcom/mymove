import React, { lazy } from 'react';
import { Switch } from 'react-router-dom';

import { MoveTabNavWithRouter } from '../shared/Header/Office';

import PrivateRoute from 'shared/User/PrivateRoute';

const MoveDetails = lazy(() => import('./Office/MoveDetails/MoveDetails'));
const TOOMoveTaskOrder = lazy(() => import('./TOO/moveTaskOrder'));
const PaymentRequestShow = lazy(() => import('../scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('./moveHistory'));

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
