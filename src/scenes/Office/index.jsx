import React, { Component } from 'react';
import { Redirect, Switch, Route } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import QueueHeader from 'shared/Header/Office';
import QueueList from './QueueList';
import QueueTable from './QueueTable';
import MoveInfo from './MoveInfo';
import OrdersInfo from './OrdersInfo';
import DocumentViewer from './DocumentViewer';
import { getCurrentUserInfo } from 'shared/Data/users';
import { loadInternalSchema, loadPublicSchema } from 'shared/Swagger/ducks';
import { no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';
import ScratchPad from 'shared/ScratchPad';
import { isProduction } from 'shared/constants';
import { RetrieveMovesForOffice } from './api';

import './office.scss';
import { detectIE11 } from '../../shared/utils';

class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <QueueList />
        </div>
        <div className="queue-list-column">
          <QueueTable queueType={this.props.match.params.queueType} retrieveMoves={RetrieveMovesForOffice} />
        </div>
      </div>
    );
  }
}

function DivTag() {
  return (
    <ConnectedRouter history={history}>
      <div className="Office site">
        <QueueHeader />
        <div role="main" className="site__content">
          <MainContent />
        </div>
      </div>
    </ConnectedRouter>
  );
}

function MainTag() {
  return (
    <ConnectedRouter history={history}>
      <div className="Office site">
        <QueueHeader />
        <main className="site__content">
          <MainContent />
        </main>
      </div>
    </ConnectedRouter>
  );
}

function MainContent() {
  return (
    <div>
      <LogoutOnInactivity />
      <Switch>
        <Route
          exact
          path="/"
          component={({ location }) => (
            <Redirect
              from="/"
              to={{
                ...location,
                pathname: '/queues/new',
              }}
            />
          )}
        />
        <PrivateRoute path="/queues/:queueType/moves/:moveId" component={MoveInfo} />
        <PrivateRoute path="/queues/:queueType" component={Queues} />
        <PrivateRoute path="/moves/:moveId/orders" component={OrdersInfo} />
        <PrivateRoute path="/moves/:moveId/documents/:moveDocumentId?" component={DocumentViewer} />
        {!isProduction && <PrivateRoute path="/playground" component={ScratchPad} />}
      </Switch>
    </div>
  );
}

class OfficeWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Office';
    this.props.loadInternalSchema();
    this.props.loadPublicSchema();
    this.props.getCurrentUserInfo();
  }

  render() {
    let isIE = detectIE11();
    // If we detect IE we'll make use of a component that uses <div role="main"> instead of <main>
    // so that we're compatible with IE11. Otherwise we'll do the normal <main>.
    if (isIE) {
      return <DivTag />;
    } else {
      return <MainTag />;
    }
  }
}

OfficeWrapper.defaultProps = {
  loadInternalSchema: no_op,
  loadPublicSchema: no_op,
};

const mapStateToProps = state => ({
  swaggerError: state.swaggerInternal.hasErrored,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadInternalSchema, loadPublicSchema, getCurrentUserInfo }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OfficeWrapper);
