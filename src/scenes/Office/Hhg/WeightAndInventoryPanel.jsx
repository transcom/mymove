import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid } from 'redux-form';

import editablePanel from '../editablePanel';
import { no_op } from 'shared/utils';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelField } from 'shared/EditablePanel';

const WeightAndInventoryDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">Weights</span>
        <PanelField title="Customer Estimate">
          {get(fieldProps, 'values.weight_estimate', '').toLocaleString()} lbs
        </PanelField>
      </div>
      <div className="editable-panel-column">
        <span className="column-subhead">Special Items</span>
        <PanelField title="Customer Entered" className="Todo-phase2">
          None
        </PanelField>
      </div>
    </React.Fragment>
  );
};

const WeightAndInventoryEdit = props => {
  const { shipmentSchema } = props;
  return (
    <div className="editable-panel-column">
      <SwaggerField
        title="Customer Estimate"
        fieldName="weight_estimate"
        swagger={shipmentSchema}
        required
      />
    </div>
  );
};

const formName = 'office_shipment_info_weight_and_inventory';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let WeightAndInventoryPanel = editablePanel(
  WeightAndInventoryDisplay,
  WeightAndInventoryEdit,
  editEnabled,
);
WeightAndInventoryPanel = reduxForm({ form: formName })(
  WeightAndInventoryPanel,
);

function mapStateToProps(state) {
  let shipment = get(state, 'office.officeMove.shipments.0', {});
  return {
    // Wrapper
    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),
    hasError:
      state.office.shipmentHasLoadError || state.office.shipmentHasUpdateError,
    errorMessage: state.office.error,

    shipment: shipment,
    isUpdating: state.office.shipmentIsUpdating,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      return [shipment.id, values];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(
  WeightAndInventoryPanel,
);
