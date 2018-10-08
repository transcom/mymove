import React, { Component } from 'react';
import { Redirect, Switch } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import DocumentUploader from 'shared/DocumentViewer/DocumentUploader';
import TspHeader from 'shared/Header/Tsp';
import { loadLoggedInUser } from 'shared/User/ducks';
import { loadPublicSchema } from 'shared/Swagger/ducks';
import { no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';
import ScratchPad from 'shared/ScratchPad';
import { isProduction } from 'shared/constants';
import DocumentViewer from './DocumentViewerContainer';
import ShipmentInfo from './ShipmentInfo';
import QueueList from './QueueList';
import QueueTable from './QueueTable';

import './tsp.css';

class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <QueueList />
        </div>
        <div className="queue-list-column">
          <div className="queue-table-scrollable">
            <QueueTable queueType={this.props.match.params.queueType} />
          </div>
        </div>
      </div>
    );
  }
}

class TestNewDocument extends Component {
  render() {
    return (
      <div>
        <DocumentUploader />
      </div>
    );
  }
}

class TspWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: TSP';
    this.props.loadPublicSchema();
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <div className="TSP site">
          <TspHeader />
          <main className="site__content">
            <div>
              <LogoutOnInactivity />
              <Switch>
                <Redirect from="/" to="/queues/new" exact />
                <PrivateRoute
                  path="/shipments/:shipmentId/documents/new"
                  component={TestNewDocument}
                />
                <PrivateRoute
                  path="/shipments/:shipmentId/documents/:moveDocumentId"
                  component={DocumentViewer}
                />
                <PrivateRoute
                  path="/shipments/:shipmentId"
                  component={ShipmentInfo}
                />
                {/* Be specific about available routes by listing them */}
                <PrivateRoute
                  path="/queues/:queueType(new|approved|in_transit|delivered|all)"
                  component={Queues}
                />
                {!isProduction && (
                  <PrivateRoute path="/playground" component={ScratchPad} />
                )}
                {/* TODO: cgilmer (2018/07/31) Need a NotFound component to route to */}
              </Switch>
            </div>
          </main>
        </div>
      </ConnectedRouter>
    );
  }
}

TspWrapper.defaultProps = {
  loadPublicSchema: no_op,
  loadLoggedInUser: no_op,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadPublicSchema, loadLoggedInUser }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TspWrapper);
