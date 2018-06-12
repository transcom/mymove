import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import {
  reduxForm,
  getFormValues,
  isValid,
  FormSection,
  Field,
} from 'redux-form';

import editablePanel from './editablePanel';
import { updateOrdersInfo } from './ducks.js';
import { formatDate, formatDateTime } from 'shared/formatters';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

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
          {get(props.serviceMember, 'current_station.name', '')}
        </PanelField>
        <PanelField title="New Duty Station">
          {get(orders, 'new_duty_station.name', '')}
        </PanelField>
        {orders.has_dependents && (
          <PanelField title="Dependents" value="Authorized" />
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

        <FormSection name="orders">
          <SwaggerField fieldName="orders_number" swagger={schema} />
          <SwaggerField
            fieldName="issue_date"
            swagger={schema}
            className="half-width"
          />
          <SwaggerField fieldName="orders_type" swagger={schema} />
          <SwaggerField fieldName="orders_type_detail" swagger={schema} />
          <PanelField
            title="Report by"
            fieldName="report_by_date"
            swagger={schema}
          />
        </FormSection>
        <FormSection name="serviceMember">
          <div className="usa-input duty-station">
            <Field
              name="current_station"
              component={DutyStationSearchBox}
              props={{ title: 'Current Duty Station' }}
            />
          </div>
        </FormSection>
        <FormSection name="orders">
          <div className="usa-input duty-station">
            <Field
              name="new_duty_station"
              component={DutyStationSearchBox}
              props={{ title: 'New Duty Station' }}
            />
          </div>
          <SwaggerField
            fieldName="has_dependents"
            swagger={schema}
            title="Dependents authorized"
          />
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
        </FormSection>
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
      serviceMember: get(state, 'office.officeServiceMember', {}),
    },

    ordersSchema: get(state, 'swagger.spec.definitions.Orders', {}),

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
