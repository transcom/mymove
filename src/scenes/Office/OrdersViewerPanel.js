import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { reduxForm, getFormValues, isValid } from 'redux-form';

import editablePanel from './editablePanel';
import { updateOrdersInfo } from './ducks.js';
import { formatDate, formatDateTime } from 'shared/formatters';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

import './office.css';

const OrdersViewerDisplay = props => {
  const orders = props.orders;
  const uploads = [];
  const ordersFieldsProps = {
    values: props.orders,
    schema: props.ordersSchema,
  };

  return (
    <React.Fragment>
      <div>
        <PanelField title="Move Locator">{props.move.locator}</PanelField>
        <PanelField title="DoD ID">{props.serviceMember.edipi}</PanelField>
        <span className="panel-subhead">
          <FontAwesomeIcon
            aria-hidden
            className="icon approval-waiting"
            icon={faClock}
            title="Awaiting Review"
          />Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </span>

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
  const schema = props.ordersSchema;
  const serviceMemberSchema = props.serviceMemberSchema;

  return (
    <React.Fragment>
      <div>
        <PanelField title="Move Locator">{props.move.locator}</PanelField>
        <PanelField title="DoD ID">{props.serviceMember.edipi}</PanelField>
        <span className="panel-subhead">
          <FontAwesomeIcon
            aria-hidden
            className="icon approval-waiting"
            icon={faClock}
            title="Awaiting Review"
          />Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </span>
        {uploads.length > 0 && (
          <p className="uploaded-at">
            Uploaded{' '}
            {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}
          </p>
        )}

        <SwaggerField fieldName="orders_number" swagger={schema} />

        <PanelField
          title="Date issued"
          fieldName="issue_date"
          swagger={schema}
        />

        <SwaggerField fieldName="orders_type" swagger={schema} />
        <SwaggerField fieldName="orders_type_detail" swagger={schema} />

        <PanelField
          title="Report by"
          fieldName="report_by_date"
          swagger={schema}
        />

        <SwaggerField
          title="Current Duty Station"
          fieldname="current_station"
          swagger={serviceMemberSchema}
        />

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

        <SwaggerField
          title="Dept. Indicator"
          fieldName="department_indicator"
          swagger={schema}
        />
        <SwaggerField title="TAC" fieldName="tac" swagger={schema} />
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

    ordersSchema: get(state, 'swagger.spec.definitions.Orders', {}),
    serviceMemberSchema: get(
      state,
      'swagger.spec.definitions.ServiceMemberPayload',
      {},
    ),

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
