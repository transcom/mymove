import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { Field, FormSection, getFormValues, reduxForm } from 'redux-form';

import { selectMove } from 'shared/Entities/modules/moves';
import { selectServiceMemberForOrders, updateServiceMember } from 'shared/Entities/modules/serviceMembers';
import { selectOrdersForMove, updateOrders } from 'shared/Entities/modules/orders';
import { formatDate, formatDateTime } from 'shared/formatters';
import { editablePanelify, PanelField, PanelSwaggerField, RowBasedHeader } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyLocationSearchBox from 'components/DutyLocationSearchBox/DutyLocationSearchBox';
import { renderStatusIcon } from 'shared/utils';

import './office.scss';

const OrdersViewerDisplay = (props) => {
  const orders = props.orders;
  const currentDutyLocation = get(props.serviceMember, 'current_location.name', '');
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
        <h3>
          {renderStatusIcon(orders.status)}
          Orders {orders.orders_number} ({formatDate(orders.issue_date)})
        </h3>
        {uploads.length > 0 && (
          <p className="uploaded-at">Uploaded {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}</p>
        )}

        <PanelSwaggerField fieldName="orders_number" required {...ordersFieldsProps} />

        <PanelField title="Date issued" required value={formatDate(orders.issue_date)} />

        <PanelSwaggerField fieldName="orders_type" required {...ordersFieldsProps} />

        <PanelSwaggerField fieldName="orders_type_detail" required {...ordersFieldsProps} />

        <PanelField title="Report by" required value={formatDate(orders.report_by_date)} />

        <PanelField title="Current Duty Location" required value={currentDutyLocation} />

        <PanelField title="New Duty Location" required value={get(orders, 'new_duty_location.name', '')} />

        {orders.has_dependents && <PanelField title="Dependents" value="Authorized" />}

        <PanelSwaggerField title="Dept. Indicator" fieldName="department_indicator" required {...ordersFieldsProps} />

        <PanelSwaggerField title="TAC / MDC" fieldName="tac" required {...ordersFieldsProps} />

        <PanelSwaggerField title="SAC / SDN" fieldName="sac" {...ordersFieldsProps} required />
      </div>
    </React.Fragment>
  );
};

const OrdersViewerEdit = (props) => {
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
          <p className="uploaded-at">Uploaded {formatDateTime(orders.uploaded_orders.uploads[0].created_at)}</p>
        )}

        <FormSection name="orders">
          <SwaggerField fieldName="orders_number" swagger={schema} required />
          <SwaggerField fieldName="issue_date" swagger={schema} />
          <SwaggerField fieldName="orders_type" swagger={schema} required />
          <SwaggerField fieldName="orders_type_detail" swagger={schema} required />
          <SwaggerField title="Report by" fieldName="report_by_date" swagger={schema} />
        </FormSection>
        <FormSection name="serviceMember">
          <div className="duty-location">
            <Field
              name="current_location"
              component={DutyLocationSearchBox}
              props={{ title: 'Current Duty Location' }}
            />
          </div>
        </FormSection>
        <FormSection name="orders">
          <div className="duty-location">
            <Field name="new_duty_location" component={DutyLocationSearchBox} props={{ title: 'New Duty Location' }} />
          </div>
          <SwaggerField fieldName="has_dependents" swagger={schema} title="Dependents authorized" />
          <SwaggerField title="Dept. Indicator" fieldName="department_indicator" swagger={schema} required />
          <SwaggerField title="TAC / MDC" fieldName="tac" swagger={schema} required />
          <SwaggerField title="SAC / SDN" fieldName="sac" swagger={schema} required />
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'orders_document_viewer';

const editEnabled = true;
let OrdersViewerPanel = editablePanelify(OrdersViewerDisplay, OrdersViewerEdit, editEnabled, RowBasedHeader);
OrdersViewerPanel = reduxForm({ form: formName })(OrdersViewerPanel);

function mapStateToProps(state, ownProps) {
  const { moveId } = ownProps;
  const orders = selectOrdersForMove(state, moveId);
  const serviceMember = selectServiceMemberForOrders(state, orders.id);

  return {
    // reduxForm
    initialValues: {
      orders,
      serviceMember,
    },

    ordersSchema: get(state, 'swaggerInternal.spec.definitions.Orders', {}),

    hasError: false,
    isUpdating: false,

    orders,
    serviceMember,
    move: selectMove(state, moveId),

    // editablePanelify
    getUpdateArgs: function () {
      let values = getFormValues(formName)(state);
      return [orders.id, values.orders, serviceMember.id, values.serviceMember];
    },
  };
}

function mapDispatchToProps(dispatch) {
  const update = (ordersId, orders, serviceMemberId, serviceMember) => {
    serviceMember.current_location_id = serviceMember.current_location.id;
    dispatch(updateServiceMember(serviceMemberId, serviceMember));

    if (!orders.has_dependents) {
      orders.spouse_has_pro_gear = false;
    }

    orders.new_duty_location_id = orders.new_duty_location.id;
    dispatch(updateOrders(ordersId, orders));
  };

  return { update };
}

export default connect(mapStateToProps, mapDispatchToProps)(OrdersViewerPanel);
