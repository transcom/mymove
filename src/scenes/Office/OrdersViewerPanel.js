import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { compact, get } from 'lodash';
import { reduxForm, getFormValues, isValid } from 'redux-form';

import editablePanel from './editablePanel';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Alert from 'shared/Alert';
import { loadMoveDependencies, updateOrdersInfo } from './ducks.js';
import { formatDate, formatDateTime } from 'shared/formatters';

import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

import './office.css';

const OrdersViewerDisplay = props => {
  const orders = props.orders;
  const uploads = [];
  const ordersFieldsProps = {
    values: props.orders,
    schema: props.schema,
  };

  return (
    <React.Fragment>
      <div>
        <PanelField title="Move Locator">{props.move.locator}</PanelField>
        <PanelField title="DoD ID">{props.serviceMember.edipi}</PanelField>

        <h3>
          Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </h3>
        {uploads.length > 0 && (
          <p className="uploaded-at">
            Uploaded{' '}
            {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}
          </p>
        )}

        <PanelSwaggerField fieldName="orders_number" {...ordersFieldsProps} />

        <PanelField title="Date issued" value={formatDate(orders.issue_date)} />

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
          <PanelField className="Todo" title="Dependents" value="Authorized" />
        )}

        <PanelSwaggerField
          title="Dept. Indicator"
          fieldName="department_indicator"
          {...ordersFieldsProps}
        />
        <PanelSwaggerField title="TAC" fieldName="tac" {...ordersFieldsProps} />

        <PanelField className="Todo" title="Doc status" />
      </div>
    </React.Fragment>
  );
};

const OrdersViewerEdit = props => {
  const orders = props.orders;
  const uploads = [];
  const schema = props.schema;

  return (
    <React.Fragment>
      <div>
        <PanelField title="Move Locator">{props.move.locator}</PanelField>
        <PanelField title="DoD ID">{props.serviceMember.edipi}</PanelField>

        <h3>
          Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </h3>
        {uploads.length > 0 && (
          <p className="uploaded-at">
            Uploaded{' '}
            {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}
          </p>
        )}

        <PanelSwaggerField fieldName="orders_number" swagger={schema} />

        <PanelField
          title="Date issued"
          fieldName="issue_date"
          swagger={schema}
        />

        <PanelSwaggerField fieldName="orders_type" swagger={schema} />
        <PanelSwaggerField fieldName="orders_type_detail" swagger={schema} />

        <PanelField
          title="Report by"
          fieldName="report_by_date"
          swagger={schema}
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
            fieldName="has_dependents"
            swagger={schema}
          />
        )}

        <PanelSwaggerField
          title="Dept. Indicator"
          fieldName="department_indicator"
          swagger={schema}
        />
        <PanelSwaggerField title="TAC" fieldName="tac" swagger={schema} />
        <PanelField
          className="Todo"
          title="Doc status"
          fieldName="orders_status"
          swagger={schema}
        />
      </div>
    </React.Fragment>
  );
};

const formName = 'orders_document_viewer';

let OrdersViewerPanel = editablePanel(OrdersViewerDisplay, OrdersViewerEdit);
OrdersViewerPanel = reduxForm({ form: formName })(OrdersViewerPanel);

function mapStateToProps(state) {
  return {
    // reduxForm
    initialValues: {
      orders: get(state, 'office.officeOrders', {}),
    },

    schema: get(state, 'swagger.spec.definitions.Orders', {}),

    hasError: false,
    errorMessage: state.office.error,
    isUpdating: false,

    orders: get(state, 'office.officeOrders', {}),
    serviceMember: get(state, 'office.officeServiceMember', {}),
    move: get(state, 'office.officeMove', {}),

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      return [
        get(state, 'office.officeOrders.id'),
        values.orders,
        get(state, 'office.officeServiceMember.id'),
        values.serviceMember,
      ];
    },
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateOrdersInfo,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(OrdersViewerPanel);
