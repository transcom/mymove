import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid } from 'redux-form';

import editablePanel from '../editablePanel';
import { no_op } from 'shared/utils';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField } from 'shared/EditablePanel';

// TODO: add shipment_type
// TODO: combine source_gbloc and destination_gbloc values to one value
// TODO: debug why code_of_service value is not in response (could be datagen needs updating)
const RoutingPanelDisplay = props => {
  return (
    <div className="editable-panel-column">
      <PanelSwaggerField
        title="Shipment market"
        fieldName="market"
        values={{
          market: props.initialValues.market,
          source_gbloc: props.initialValues.source_gbloc,
          code_of_service: props.initialValues.code_of_service,
        }}
        schema={props.shipmentSchema}
      />
      <PanelSwaggerField
        title="Channel"
        fieldName="source_gbloc"
        values={{
          source_gbloc: props.initialValues.source_gbloc,
        }}
        schema={props.shipmentSchema}
      />
      <PanelSwaggerField
        title="Code of service"
        fieldName="code_of_service"
        values={{
          code_of_service: props.initialValues.code_of_service,
        }}
        schema={props.shipmentSchema}
      />
    </div>
  );
};

const RoutingPanelEdit = props => {
  const { shipmentSchema } = props;
  return (
    <div className="editable-panel-column">
      <SwaggerField
        title="Market"
        fieldName="market"
        swagger={shipmentSchema}
        required
      />
      <SwaggerField
        title="Channel"
        fieldName="source_gbloc"
        swagger={shipmentSchema}
        required
      />
      <SwaggerField
        title="Code of service"
        fieldName="code_of_service"
        swagger={shipmentSchema}
        required
      />
    </div>
  );
};

const formName = 'office_shipment_routing';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let RoutingPanel = editablePanel(
  RoutingPanelDisplay,
  RoutingPanelEdit,
  editEnabled,
);
RoutingPanel = reduxForm({ form: formName })(RoutingPanel);

function mapStateToProps(state) {
  let shipment = get(state, 'office.officeMove.shipments.0', {});

  return {
    // reduxForm
    initialValues: shipment,

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

export default connect(mapStateToProps, mapDispatchToProps)(RoutingPanel);
