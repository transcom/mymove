import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, Field, FormSection, getFormValues } from 'redux-form';
import editablePanel from './editablePanel';

import { updateOrdersInfo } from './ducks';
import { loadEntitlements } from 'scenes/Office/ducks';

import {
  PanelSwaggerField,
  PanelField,
  SwaggerValue,
} from 'shared/EditablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

function renderEntitlements(entitlements) {
  return (
    <React.Fragment>
      <span className="panel-subhead">Entitlements</span>
      <PanelField title="Household Goods">
        {get(entitlements, 'total', '').toLocaleString()} lbs
      </PanelField>
      <PanelField title="Pro-gear">
        {get(entitlements, 'pro_gear', '').toLocaleString()} lbs
      </PanelField>
      <PanelField title="Spouse pro-gear">
        {get(entitlements, 'pro_gear_spouse', '').toLocaleString()} lbs
      </PanelField>
      <PanelField className="Todo" title="Short-term storage">
        90 days
      </PanelField>
    </React.Fragment>
  );
}

const OrdersDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelField title="Orders Number">
          <a
            href={get(props.orders, 'uploaded_orders.uploads[0].url')}
            target="_blank"
          >
            <SwaggerValue fieldName="orders_number" {...fieldProps} />&nbsp;
            <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
          </a>
        </PanelField>
        <PanelSwaggerField
          title="Date issued"
          fieldName="issue_date"
          {...fieldProps}
        />
        <PanelSwaggerField fieldName="orders_type" {...fieldProps} />
        <PanelSwaggerField fieldName="orders_type_detail" {...fieldProps} />
        <PanelSwaggerField
          title="Report by"
          fieldName="report_by_date"
          {...fieldProps}
        />
        <PanelField title="Current Duty Station">
          {get(props.serviceMember, 'current_station.name', '')}
        </PanelField>
        <PanelField title="New Duty Station">
          {get(props.orders, 'new_duty_station.name', '')}
        </PanelField>
      </div>
      <div className="editable-panel-column">
        {renderEntitlements(props.entitlements)}
        {props.orders.has_dependents && (
          <PanelField title="Dependents" value="Authorized" />
        )}
      </div>
    </React.Fragment>
  );
};

const OrdersEdit = props => {
  const schema = props.ordersSchema;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <FormSection name="orders">
          <SwaggerField
            fieldName="orders_number"
            swagger={schema}
            className="half-width"
          />
          <SwaggerField
            fieldName="issue_date"
            swagger={schema}
            className="half-width"
          />
          <SwaggerField fieldName="orders_type" swagger={schema} />
          <SwaggerField fieldName="orders_type_detail" swagger={schema} />
          <SwaggerField fieldName="report_by_date" swagger={schema} />
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
        </FormSection>
      </div>
      <div className="editable-panel-column">
        {renderEntitlements(props.entitlements)}

        <FormSection name="orders">
          <SwaggerField
            fieldName="has_dependents"
            swagger={schema}
            title="Dependents authorized"
          />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_orders';

let OrdersPanel = editablePanel(OrdersDisplay, OrdersEdit);
OrdersPanel = reduxForm({ form: formName })(OrdersPanel);

function mapStateToProps(state) {
  return {
    // reduxForm
    initialValues: {
      orders: get(state, 'office.officeOrders', {}),
      serviceMember: get(state, 'office.officeServiceMember', {}),
    },

    ordersSchema: get(state, 'swagger.spec.definitions.Orders', {}),
    serviceMemberSchema: get(
      state,
      'swagger.spec.definitions.ServiceMemberPayload',
      {},
    ),

    hasError: false,
    errorMessage: state.office.error,
    orders: get(state, 'office.officeOrders', {}),
    serviceMember: get(state, 'office.officeServiceMember', {}),

    entitlements: loadEntitlements(state),
    isUpdating: false,

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

export default connect(mapStateToProps, mapDispatchToProps)(OrdersPanel);
