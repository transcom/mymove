import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid } from 'redux-form';
import editablePanel from './editablePanel';

import { updateOrders } from './ducks';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

const AccountingDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        {props.orders.department_indicator ? (
          <PanelSwaggerField
            title="Department indicator"
            fieldName="department_indicator"
            {...fieldProps}
          />
        ) : (
          <PanelField title="Department indicator" className="missing">
            missing
          </PanelField>
        )}
      </div>
      <div className="editable-panel-column">
        {props.orders.tac ? (
          <PanelSwaggerField title="TAC" fieldName="tac" {...fieldProps} />
        ) : (
          <PanelField title="TAC" className="missing">
            missing
          </PanelField>
        )}
      </div>
    </React.Fragment>
  );
};

const AccountingEdit = props => {
  const { ordersSchema } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField
          title="Department indicator"
          fieldName="department_indicator"
          swagger={ordersSchema}
          required
        />
      </div>
      <div className="editable-panel-column">
        <SwaggerField
          title="TAC"
          fieldName="tac"
          swagger={ordersSchema}
          required
        />
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_accounting';

let AccountingPanel = editablePanel(AccountingDisplay, AccountingEdit);
AccountingPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(AccountingPanel);

function mapStateToProps(state) {
  let orders = get(state, 'office.officeOrders', {});

  return {
    // reduxForm
    initialValues: state.office.officeOrders,

    // Wrapper
    ordersSchema: get(state, 'swagger.spec.definitions.Orders', {}),
    hasError:
      state.office.ordersHaveLoadError || state.office.ordersHaveUpdateError,
    errorMessage: state.office.error,

    orders: orders,
    isUpdating: state.office.ordersAreUpdating,

    // editablePanel
    formIsValid: isValid(formName)(state),
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
