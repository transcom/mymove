import React, { Component } from 'react';
import { Redirect, Switch, Route } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import TspHeader from 'shared/Header/Tsp';
import { getCurrentUserInfo } from 'shared/Data/users';
import { loadPublicSchema } from 'shared/Swagger/ducks';
import { detectIE11, no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';
import ScratchPad from 'shared/ScratchPad';
import { isProduction } from 'shared/constants';
import DocumentViewer from './DocumentViewerContainer';
import NewDocument from './NewDocumentContainer';
import ShipmentInfo from './ShipmentInfo';
import QueueList from './QueueList';
import QueueTable from './QueueTable';

import './tsp.scss';

class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <QueueList />
        </div>
        <div className="queue-list-column">
          <QueueTable queueType={this.props.match.params.queueType} />
        </div>
      </div>
    );
  }
}

class TspWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: TSP';
    this.props.loadPublicSchema();
    this.props.getCurrentUserInfo();
  }

  render() {
    const Tag = detectIE11() ? 'div' : 'main';
    return (
      <ConnectedRouter history={history}>
        <div className="TSP site">
          <TspHeader />
          <Tag role="main" className="site__content">
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
                <PrivateRoute path="/shipments/:shipmentId/documents/new" component={NewDocument} />
                <PrivateRoute path="/shipments/:shipmentId/documents/:moveDocumentId" component={DocumentViewer} />
                <PrivateRoute path="/shipments/:shipmentId" component={ShipmentInfo} />
                {/* Be specific about available routes by listing them */}
                <PrivateRoute
                  path="/queues/:queueType(new|accepted|approved|in_transit|delivered|completed|all)"
                  component={Queues}
                />
                {!isProduction && <PrivateRoute path="/playground" component={ScratchPad} />}
                {/* TODO: cgilmer (2018/07/31) Need a NotFound component to route to */}
                <Redirect from="*" to="/queues/new" component={Queues} />
              </Switch>
            </div>
          </Tag>
        </div>
      </ConnectedRouter>
    );
  }
}

TspWrapper.defaultProps = {
  loadPublicSchema: no_op,
  getCurrentUserInfo: no_op,
};

const mapStateToProps = state => ({
  swaggerError: state.swaggerPublic.hasErrored,
});

const mapDispatchToProps = dispatch => bindActionCreators({ loadPublicSchema, getCurrentUserInfo }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TspWrapper);
