import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { compact, get } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { loadMoveDependencies } from './ducks.js';
import { formatDate, formatDateTime } from 'shared/formatters';

import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

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
    const ordersFieldsProps = {
      values: this.props.orders,
      schema: this.props.ordersSchema,
    };

    const move = this.props.move;
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
            <h2 className="usa-heading">{name}</h2>

            <PanelField title="Move Locator">{move.locator}</PanelField>
            <PanelField title="DoD ID">{serviceMember.edipi}</PanelField>

            <h3>
              Orders {orders.orders_number} ({formatDate(orders.issue_date)})
            </h3>
            {uploads.length > 0 && (
              <p className="uploaded-at">
                Uploaded{' '}
                {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}
              </p>
            )}

            <PanelSwaggerField
              fieldName="orders_number"
              {...ordersFieldsProps}
            />

            <PanelField
              title="Date issued"
              value={formatDate(orders.issue_date)}
            />

            <PanelSwaggerField fieldName="orders_type" {...ordersFieldsProps} />
            <PanelSwaggerField
              fieldName="orders_type_detail"
              {...ordersFieldsProps}
            />

            <PanelField
              title="Report by"
              value={formatDate(orders.report_by_date)}
            />

            <PanelField title="Current Duty Station">
              {orders.current_duty_station && orders.current_duty_station.name}
            </PanelField>
            <PanelField title="New Duty Station">
              {orders.new_duty_station && orders.new_duty_station.name}
            </PanelField>

            {orders.has_dependents && (
              <PanelField
                className="Todo"
                title="Dependents"
                value="Authorized"
              />
            )}

            <PanelSwaggerField
              title="Dept. Indicator"
              fieldName="department_indicator"
              {...ordersFieldsProps}
            />
            <PanelSwaggerField
              title="TAC"
              fieldName="tac"
              {...ordersFieldsProps}
            />

            <PanelField className="Todo" title="Doc status" />
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
  move: state.office.officeMove || {},
  orders: state.office.officeOrders || {},
  serviceMember: state.office.officeServiceMember || {},
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMoveDependencies }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OrdersInfo);
