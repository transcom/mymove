import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid } from 'redux-form';

import editablePanel from '../editablePanel';
import { no_op } from 'shared/utils';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField } from 'shared/EditablePanel';

const LocationsDisplay = props => {
  return (
    <div className="editable-panel-column">
      <PanelSwaggerField
        title="Primary"
        fieldName="pickup_address"
        values={props.initialValues.shipment.pickup_address}
        schema={props.shipmentSchema}
      />
    </div>
  );
};

const LocationsEdit = props => {
  const { shipmentSchema } = props;
  return (
    <div className="editable-panel-column">
      <SwaggerField
        title="Primary"
        fieldName="pickup_address"
        swagger={shipmentSchema}
        required
      />
    </div>
  );
};

const formName = 'office_shipment_info_locations';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let LocationsPanel = editablePanel(
  LocationsDisplay,
  LocationsEdit,
  editEnabled,
);
LocationsPanel = reduxForm({ form: formName })(LocationsPanel);

function mapStateToProps(state) {
  let shipment = get(state, 'office.officeMove.shipments.0', {});
  return {
    // reduxForm
    initialValues: {
      shipment: shipment,
    },

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

export default connect(mapStateToProps, mapDispatchToProps)(LocationsPanel);
