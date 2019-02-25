import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import { updateOrders, selectOrdersForMove } from 'shared/Entities/modules/orders';
import { selectMove } from 'shared/Entities/modules/moves';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';

const AccountingDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };
  const { sacIsRequired } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField title="Department indicator" fieldName="department_indicator" required {...fieldProps} />

        <PanelSwaggerField title="SAC" required={sacIsRequired} fieldName="sac" {...fieldProps} />
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField title="TAC" required fieldName="tac" {...fieldProps} />
      </div>
    </React.Fragment>
  );
};

const AccountingEdit = props => {
  const { ordersSchema, sacIsRequired } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField title="Department indicator" fieldName="department_indicator" swagger={ordersSchema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField title="TAC" fieldName="tac" swagger={ordersSchema} required />
      </div>
      <div className="editable-panel-column">
        <SwaggerField title="SAC" fieldName="sac" swagger={ordersSchema} required={sacIsRequired} />
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
  const { moveId } = ownProps;
  const move = selectMove(state, moveId);
  const orders = selectOrdersForMove(state, moveId);

  return {
    // reduxForm
    initialValues: orders,

    // Wrapper
    ordersSchema: get(state, 'swaggerInternal.spec.definitions.Orders', {}),
    orders: orders,
    sacIsRequired: !move.selected_move_type === 'PPM',

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
