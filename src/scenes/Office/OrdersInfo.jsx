import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import DocumentContent from 'shared/DocumentViewer/DocumentContent';
import OrdersViewerPanel from './OrdersViewerPanel';
import { loadMoveDependencies } from './ducks.js';
import { stringifyName } from 'shared/utils/serviceMember';

import './office.css';

class OrdersInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
  }

  render() {
    const orders = this.props.orders;
    const serviceMember = this.props.serviceMember;
    const name = stringifyName(serviceMember);

    let uploads;
    if (orders && orders.uploaded_orders) {
      uploads = orders.uploaded_orders.uploads.map(upload => (
        <DocumentContent
          key={upload.url}
          url={upload.url}
          filename={upload.filename}
          contentType={upload.content_type}
        />
      ));
    } else {
      uploads = [];
    }

    if (!this.props.loadDependenciesHasSuccess && !this.props.loadDependenciesHasError) return <LoadingPlaceholder />;
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
          <div className="usa-width-two-thirds document-contents">{uploads}</div>
          <div className="usa-width-one-third orders-page-fields">
            <OrdersViewerPanel title={name} className="document-viewer" moveId={this.props.match.params.moveId} />
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
  swaggerError: state.swaggerInternal.hasErrored,
  ordersSchema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
  orders: state.office.officeOrders || {},
  serviceMember: state.office.officeServiceMember || {},
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
});

const mapDispatchToProps = dispatch => bindActionCreators({ loadMoveDependencies }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OrdersInfo);
