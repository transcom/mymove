import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { compact, get } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import OrdersViewerPanel from './OrdersViewerPanel';
import { loadMoveDependencies } from './ducks.js';

import './office.css';

// Page displays an image or PDF.
const Page = function(props) {
  let content;
  if (props.contentType === 'application/pdf') {
    content = (
      <div className="pdf-placeholder">
        {props.filename && <span className="filename">{props.filename}</span>}
        This PDF can be <a href={props.url}>viewed here</a>.
      </div>
    );
  } else {
    content = (
      <img src={props.url} width="100%" height="100%" alt="document upload" />
    );
  }
  return <div className="page">{content}</div>;
};

class OrdersInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
  }

  render() {
    const orders = this.props.orders;
    const serviceMember = this.props.serviceMember;
    const name = compact([
      serviceMember.last_name,
      serviceMember.first_name,
    ]).join(', ');

    let uploads;
    if (orders && orders.uploaded_orders) {
      uploads = orders.uploaded_orders.uploads.map(upload => (
        <Page
          key={upload.url}
          url={upload.url}
          filename={upload.filename}
          contentType={upload.content_type}
        />
      ));
    } else {
      uploads = [];
    }

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
            {uploads}
          </div>
          <div className="usa-width-one-third orders-page-fields">
            <OrdersViewerPanel
              title={name}
              className="document-viewer"
              moveId={this.props.match.params.moveId}
            />
          </div>
        </div>
      </div>
    );
  }
}

OrdersInfo.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  ordersSchema: get(state, 'swagger.spec.definitions.CreateUpdateOrders', {}),
  orders: state.office.officeOrders || {},
  serviceMember: state.office.officeServiceMember || {},
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMoveDependencies }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OrdersInfo);
