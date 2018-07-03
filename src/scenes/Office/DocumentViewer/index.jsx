import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { compact, get } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { PanelField } from 'shared/EditablePanel';
import { loadMoveDependencies } from '../ducks.js';
import { RoutedTabs, NavTab } from 'react-router-tabs';
import PrivateRoute from 'shared/User/PrivateRoute';
import { Switch, Redirect } from 'react-router-dom';

import './index.css';
class DocumentViewer extends Component {
  componentDidMount() {
    //this is probably overkill, but works for now
    this.props.loadMoveDependencies(this.props.match.params.moveId);
  }
  componentWillUpdate() {
    document.title = 'Document Viewer';
  }
  render() {
    const { serviceMember, move } = this.props;

    const name = compact([
      serviceMember.last_name,
      serviceMember.first_name,
    ]).join(', ');

    const listUrl = `${this.props.match.path}/list`;
    const detailUrl = `${this.props.match.path}/details`;
    const newUrl = `${this.props.match.url}/new`;

    if (
      !this.props.loadDependenciesHasSuccess &&
      !this.props.loadDependenciesHasError
    )
      return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="usa-grid">
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              Something went wrong contacting the server.
            </Alert>
          </div>
        </div>
      );
    return (
      <div className="usa-grid doc-viewer">
        <div className="usa-width-two-thirds">
          <div style={{ minWidth: '400px' }}> Document Upload Coming soon</div>
        </div>
        <div className="usa-width-one-third">
          <h3>{name}</h3>
          <PanelField title="Move Locator">{move.locator}</PanelField>
          <PanelField title="DoD ID">{serviceMember.edipi}</PanelField>
          <RoutedTabs
            startPathWith={this.props.match.url}
            className="doc-viewer-tabs"
          >
            <NavTab to="/list">
              <span className="title">Document(s)</span>
            </NavTab>

            <NavTab to="/details">
              <span className="title">Details</span>
            </NavTab>
          </RoutedTabs>
          <div className="tab-content">
            <Switch>
              <PrivateRoute
                exact
                path={this.props.match.url}
                render={() => <Redirect replace to={listUrl} />}
              />
              <PrivateRoute
                path={newUrl}
                render={() => <div> new list coming soon</div>}
              />
              <PrivateRoute
                path={listUrl}
                render={() => <div>list coming soon</div>}
              />
              <PrivateRoute
                path={detailUrl}
                render={() => <div> details coming soon</div>}
              />
            </Switch>
          </div>
        </div>
      </div>
    );
  }
}

DocumentViewer.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  ordersSchema: get(state, 'swagger.spec.definitions.CreateUpdateOrders', {}),
  orders: state.office.officeOrders || {},
  move: get(state, 'office.officeMove', {}),
  serviceMember: state.office.officeServiceMember || {},
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
  error: state.office.error,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMoveDependencies }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(DocumentViewer);
