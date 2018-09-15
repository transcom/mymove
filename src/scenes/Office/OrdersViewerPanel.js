import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { reduxForm, getFormValues, FormSection, Field } from 'redux-form';

import { updateOrdersInfo } from './ducks.js';
import { formatDate, formatDateTime } from 'shared/formatters';
import {
  PanelSwaggerField,
  PanelField,
  editablePanelify,
} from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import { renderStatusIcon } from 'shared/utils';

import './office.css';

const OrdersViewerDisplay = props => {
  const orders = props.orders;
  const currentDutyStation = get(
    props.serviceMember,
    'current_station.name',
    '',
  );
  const uploads = get(orders, 'uploaded_orders.uploads', []);
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
          {renderStatusIcon(orders.status)}
          Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </span>
        {uploads.length > 0 && (
          <p className="uploaded-at">
            Uploaded{' '}
            {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}
          </p>
        )}

        <PanelSwaggerField
          fieldName="orders_number"
          required
          {...ordersFieldsProps}
        />

        <PanelField
          title="Date issued"
          required
          value={formatDate(orders.issue_date)}
        />

        <PanelSwaggerField
          fieldName="orders_type"
          required
          {...ordersFieldsProps}
        />

        <PanelSwaggerField
          fieldName="orders_type_detail"
          required
          {...ordersFieldsProps}
        />

        <PanelField
          title="Report by"
          required
          value={formatDate(orders.report_by_date)}
        />

        <PanelField
          title="Current Duty Station"
          required
          value={currentDutyStation}
        />

        <PanelField
          title="New Duty Station"
          required
          value={get(orders, 'new_duty_station.name', '')}
        />

        {orders.has_dependents && (
          <PanelField title="Dependents" value="Authorized" />
        )}

        <PanelSwaggerField
          title="Dept. Indicator"
          fieldName="department_indicator"
          required
          {...ordersFieldsProps}
        />

        <PanelSwaggerField
          title="Orders Issuing Agency"
          fieldName="orders_issuing_agency"
          {...ordersFieldsProps}
        />
        <PanelSwaggerField
          title="Paragraph Number"
          fieldName="paragraph_number"
          {...ordersFieldsProps}
        />

        <PanelSwaggerField
          title="TAC"
          fieldName="tac"
          required
          {...ordersFieldsProps}
        />

        <PanelSwaggerField title="SAC" fieldName="sac" {...ordersFieldsProps} />
      </div>
    </React.Fragment>
  );
};

const OrdersViewerEdit = props => {
  const orders = props.orders;
  const uploads = get(orders, 'uploaded_orders.uploads', []);
  const schema = props.ordersSchema;

  return (
    <React.Fragment>
      <div>
        <PanelField title="Move Locator">{props.move.locator}</PanelField>
        <PanelField title="DoD ID">{props.serviceMember.edipi}</PanelField>
        <span className="panel-subhead">
          {renderStatusIcon(orders.status)}
          Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </span>
        {uploads.length > 0 && (
          <p className="uploaded-at">
            Uploaded{' '}
            {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}
          </p>
        )}

        <FormSection name="orders">
          <SwaggerField fieldName="orders_number" swagger={schema} required />
          <SwaggerField fieldName="issue_date" swagger={schema} />
          <SwaggerField fieldName="orders_type" swagger={schema} required />
          <SwaggerField
            fieldName="orders_type_detail"
            swagger={schema}
            required
          />
          <SwaggerField
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
          <SwaggerField
            title="Orders Issuing Agency"
            fieldName="orders_issuing_agency"
            swagger={schema}
          />
          <SwaggerField
            title="Paragraph Number"
            fieldName="paragraph_number"
            swagger={schema}
          />
          <SwaggerField title="TAC" fieldName="tac" swagger={schema} />
          <SwaggerField title="SAC" fieldName="sac" swagger={schema} />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'orders_document_viewer';

let OrdersViewerPanel = editablePanelify(OrdersViewerDisplay, OrdersViewerEdit);
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

    // editablePanelify
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
