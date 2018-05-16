import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert'; // eslint-disable-line
import { loadMoveDependencies } from './ducks.js';

// import FontAwesomeIcon from '@fortawesome/react-fontawesome';
// import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';

import './office.css';

const Page = function(props) {
  let content;
  if (props.contentType === 'application/pdf') {
    content = (
      <div>
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
    const ordersFieldsProps = {
      values: this.props.officeOrders,
      schema: this.props.officeSchema,
    };

    const officeMove = this.props.officeMove || {};
    const officeOrders = this.props.officeOrders || {};
    const officeServiceMember = this.props.officeServiceMember || {};

    let uploads;
    if (officeOrders && officeOrders.uploaded_orders) {
      uploads = officeOrders.uploaded_orders.uploads.map(upload => (
        <Page
          key={upload.url}
          url={upload.url}
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
          <div className="usa-width-one-third nav-controls">
            <PanelSwaggerField
              fieldName="orders_number"
              {...ordersFieldsProps}
            />
            <PanelSwaggerField
              title="Date issued"
              fieldName="issue_date"
              {...fieldProps}
            />
            <PanelSwaggerField fieldName="orders_type" {...ordersFieldsProps} />
            <PanelSwaggerField
              fieldName="orders_type_detail"
              {...ordersFieldsProps}
            />
            <PanelSwaggerField
              title="Report by"
              fieldName="report_by_date"
              {...ordersFieldsProps}
            />
            <PanelField title="Current Duty Station">
              {officeOrders.current_duty_station &&
                officeOrders.current_duty_station.name}
            </PanelField>
            <PanelField title="New Duty Station">
              {officeOrders.new_duty_station &&
                officeOrders.new_duty_station.name}
            </PanelField>
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
  officeMove: state.office.officeMove,
  officeOrders: state.office.officeOrders,
  officeServiceMember: state.office.officeServiceMember,
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMoveDependencies }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OrdersInfo);
