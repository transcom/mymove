import React, { Suspense, lazy } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router';
import { NavLink, Switch } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import PrivateRoute from 'shared/User/PrivateRoute';
import { roleTypes } from 'constants/userRoles';
import TabNav from 'components/TabNav';
import { MatchShape } from 'types/router';
import { withContext } from 'shared/AppContext';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const MoveDetails = lazy(() => import('pages/Office/MoveDetails/MoveDetails'));
const TOOMoveTaskOrder = lazy(() => import('pages/TOO/moveTaskOrder'));
const PaymentRequestShow = lazy(() => import('scenes/Office/TIO/paymentRequestShow'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));

const TXOMoveInfo = ({
  context: {
    flags: { too, tio },
  },
  match,
}) => {
  const { moveOrderId } = match.params;

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
          {too && (
            <PrivateRoute
              path="/moves/:moveOrderId/details"
              exact
              component={MoveDetails}
              requiredRoles={[roleTypes.TOO]}
              hideSwitcher
            />
          )}
          {too && (
            <PrivateRoute
              path="/moves/:moveTaskOrderId/mto"
              exact
              component={TOOMoveTaskOrder}
              requiredRoles={[roleTypes.TOO]}
              hideSwitcher
            />
          )}
          {tio && (
            <PrivateRoute
              path="/moves/:id/payment-requests"
              exact
              component={PaymentRequestShow}
              requiredRoles={[roleTypes.TIO]}
              hideSwitcher
            />
          )}
          {tio && (
            <PrivateRoute
              path="/moves/:moveOrderId/history"
              exact
              component={MoveHistory}
              requiredRoles={[roleTypes.TIO]}
              hideSwitcher
            />
          )}
        </Switch>
      </Suspense>
    </>
  );
};

TXOMoveInfo.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      too: PropTypes.bool,
      tio: PropTypes.bool,
    }),
  }),
  match: MatchShape.isRequired,
};

TXOMoveInfo.defaultProps = {
  context: {
    flags: {
      too: false,
      tio: false,
    },
  },
};

export default withContext(withRouter(TXOMoveInfo));
