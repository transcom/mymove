import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { reduxForm, Field, FormSection, getFormValues } from 'redux-form';
import { Link } from 'react-router-dom';

import { loadEntitlements } from './ducks';
import { updateServiceMember } from 'shared/Entities/modules/serviceMembers';
import { selectOrdersForMove, updateOrders } from 'shared/Entities/modules/orders';
import { selectServiceMemberForOrders } from 'shared/Entities/modules/serviceMembers';
import { formatDate } from 'shared/formatters';

import { PanelSwaggerField, PanelField, SwaggerValue, editablePanelify } from 'shared/EditablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

function renderEntitlements(entitlements, orders) {
  return (
    <React.Fragment>
      <span className="panel-subhead">Entitlements</span>
      <PanelField title="Household Goods">{get(entitlements, 'weight', '').toLocaleString()} lbs</PanelField>
      <PanelField title="Pro-gear">{get(entitlements, 'pro_gear', '').toLocaleString()} lbs</PanelField>
      {orders.spouse_has_pro_gear && (
        <PanelField title="Spouse pro-gear">{get(entitlements, 'pro_gear_spouse', '').toLocaleString()} lbs</PanelField>
      )}
      <PanelField title="Short-term storage">90 days</PanelField>
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
        {props.orders.orders_number ? (
          <PanelField title="Orders Number" className="orders_number">
            <Link to={`/moves/${props.move.id}/orders`} target="_blank">
              <SwaggerValue fieldName="orders_number" {...fieldProps} />
              &nbsp;
              <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
            </Link>
          </PanelField>
        ) : (
          <PanelField title="Orders Number" className="missing orders_number">
            missing
            <Link to={`/moves/${props.move.id}/orders`} target="_blank">
              <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
            </Link>
          </PanelField>
        )}
        <PanelField title="Date issued" value={formatDate(props.orders.issue_date)} />
        <PanelSwaggerField fieldName="orders_type" {...fieldProps} />
        <PanelSwaggerField fieldName="orders_type_detail" required {...fieldProps} />
        <PanelField title="Report by" value={formatDate(props.orders.report_by_date)} />
        <PanelField title="Current Duty Station">{get(props.serviceMember, 'current_station.name', '')}</PanelField>
        <PanelField title="New Duty Station">{get(props.orders, 'new_duty_station.name', '')}</PanelField>

        <PanelSwaggerField title="Orders Issuing Agency" fieldName="orders_issuing_agency" {...fieldProps} />

        <PanelSwaggerField title="Paragraph Number" fieldName="paragraph_number" {...fieldProps} />
      </div>
      <div className="editable-panel-column">
        {renderEntitlements(props.entitlements, props.orders)}
        {props.orders.has_dependents && <PanelField title="Dependents" value="Authorized" />}
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
          <SwaggerField fieldName="orders_number" swagger={schema} className="half-width" required />
          <SwaggerField fieldName="issue_date" swagger={schema} className="half-width" />
          <SwaggerField fieldName="orders_type" swagger={schema} required />
          <SwaggerField fieldName="orders_type_detail" swagger={schema} required />
          <SwaggerField fieldName="report_by_date" swagger={schema} />
        </FormSection>

        <FormSection name="serviceMember">
          <div className="usa-input duty-station">
            <Field name="current_station" component={DutyStationSearchBox} props={{ title: 'Current Duty Station' }} />
          </div>
        </FormSection>

        <FormSection name="orders">
          <div className="usa-input duty-station">
            <Field name="new_duty_station" component={DutyStationSearchBox} props={{ title: 'New Duty Station' }} />
          </div>
        </FormSection>

        <FormSection name="orders">
          <SwaggerField fieldName="orders_issuing_agency" swagger={schema} className="half-width" />

          <SwaggerField fieldName="paragraph_number" swagger={schema} className="half-width" />
        </FormSection>
      </div>
      <div className="editable-panel-column">
        {renderEntitlements(props.entitlements, props.orders)}

        <FormSection name="orders">
          <SwaggerField fieldName="has_dependents" swagger={schema} title="Dependents authorized" />
          {get(props, 'formValues.orders.has_dependents', false) && (
            <SwaggerField fieldName="spouse_has_pro_gear" swagger={schema} title="Spouse has pro gear" />
          )}
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_orders';

let OrdersPanel = editablePanelify(OrdersDisplay, OrdersEdit);
OrdersPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(OrdersPanel);

function mapStateToProps(state, ownProps) {
  let formValues = getFormValues(formName)(state);
  const orders = selectOrdersForMove(state, ownProps.moveId);
  const serviceMember = selectServiceMemberForOrders(state, orders.id);

  return {
    // reduxForm
    formValues: formValues,
    initialValues: { orders, serviceMember },
    ordersSchema: get(state, 'swaggerInternal.spec.definitions.Orders', {}),
    hasError: false,
    errorMessage: state.office.error,
    entitlements: loadEntitlements(state, ownProps.moveId),
    isUpdating: false,
    orders,
    serviceMember,
    move: get(state, 'office.officeMove', {}),
    // editablePanelify
    getUpdateArgs: () => [orders.id, formValues.orders, serviceMember.id, formValues.serviceMember],
  };
}

function mapDispatchToProps(dispatch) {
  const update = (ordersId, orders, serviceMemberId, serviceMember) => {
    serviceMember.current_station_id = serviceMember.current_station.id;
    dispatch(updateServiceMember(serviceMemberId, { serviceMember }));

    if (!orders.has_dependents) {
      orders.spouse_has_pro_gear = false;
    }

    orders.new_duty_station_id = orders.new_duty_station.id;
    dispatch(updateOrders(ordersId, orders));
  };

  return { update };
}

export default connect(mapStateToProps, mapDispatchToProps)(OrdersPanel);
