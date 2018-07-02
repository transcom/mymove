import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { compact, get } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { PanelField } from 'shared/EditablePanel';
import { loadMoveDependencies } from './ducks.js';

import './office.css';

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
      <div>
        <div className="usa-grid">
          <div className="usa-width-two-thirds orders-page-column">
            <div style={{ minWidth: '400px' }}>
              {' '}
              Document Upload Coming soon
            </div>
          </div>
          <div className="usa-width-one-third orders-page-fields">
            <h3>{name}</h3>
            <PanelField title="Move Locator">{move.locator}</PanelField>
            <PanelField title="DoD ID">{serviceMember.edipi}</PanelField>
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
