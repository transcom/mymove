import { get, isEmpty } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import { updateOrders, selectOrdersForMove } from 'shared/Entities/modules/orders';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';

const AccountingDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };
  const isRequired = props.isHHG ? true : false;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField title="Department indicator" fieldName="department_indicator" required {...fieldProps} />

        <PanelSwaggerField title="SAC" required={isRequired} fieldName="sac" {...fieldProps} />
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField title="TAC" required fieldName="tac" {...fieldProps} />
      </div>
    </React.Fragment>
  );
};

const AccountingEdit = props => {
  const { ordersSchema } = props;
  const isRequired = props.isHHG ? true : false;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField title="Department indicator" fieldName="department_indicator" swagger={ordersSchema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField title="TAC" fieldName="tac" swagger={ordersSchema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField title="SAC" fieldName="sac" swagger={ordersSchema} required={isRequired} />
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_accounting';

let AccountingPanel = editablePanelify(AccountingDisplay, AccountingEdit);
AccountingPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(AccountingPanel);

function mapStateToProps(state, ownProps) {
  const orders = selectOrdersForMove(state, ownProps.moveId);

  return {
    // reduxForm
    initialValues: orders,

    // Wrapper
    ordersSchema: get(state, 'swaggerInternal.spec.definitions.Orders', {}),
    hasError: state.office.ordersHaveLoadError || state.office.ordersHaveUpdateError,
    errorMessage: state.office.error,

    orders: orders,
    isHHG: !isEmpty(get(state, 'office.officeMove.shipments.0', {})),
    isUpdating: state.office.ordersAreUpdating,

    // editablePanelify
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      values.new_duty_station_id = values.new_duty_station.id;
      return [orders.id, values];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateOrders,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(AccountingPanel);
