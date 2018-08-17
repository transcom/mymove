import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid } from 'redux-form';

import { no_op } from 'shared/utils';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import {
  PanelField,
  PanelSwaggerField,
  SwaggerValue,
  editablePanelify,
} from 'shared/EditablePanel';

// TODO: add shipment_type
// TODO: combine source_gbloc and destination_gbloc values to one value
const RoutingPanelDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  return (
    <div className="editable-panel-column">
      <PanelField title="Shipment type">HHG</PanelField>
      <PanelSwaggerField
        title="Shipment market"
        fieldName="market"
        {...fieldProps}
      />
      <PanelField title="Channel">
        <SwaggerValue fieldName="source_gbloc" {...fieldProps} /> -{' '}
        <SwaggerValue fieldName="destination_gbloc" {...fieldProps} />
      </PanelField>
      <PanelSwaggerField
        title="Code of service"
        fieldName="code_of_service"
        {...fieldProps}
      />
    </div>
  );
};

const RoutingPanelEdit = props => {
  const { shipmentSchema } = props;
  return (
    <div className="editable-panel-column">
      <PanelField title="Shipment type">HHG</PanelField>
      <PanelField title="Market">dHHG</PanelField>
      <SwaggerField
        title="Source GBLOC"
        fieldName="source_gbloc"
        swagger={shipmentSchema}
        required
      />
      <SwaggerField
        title="Destination GBLOC"
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

let RoutingPanel = editablePanelify(
  RoutingPanelDisplay,
  RoutingPanelEdit,
  editEnabled,
);
RoutingPanel = reduxForm({ form: formName })(RoutingPanel);

function mapStateToProps(state) {
  let shipment = get(state, 'office.officeMove.shipments.0', {});

  return {
    initialValues: shipment,
    // Wrapper
    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),
    hasError:
      state.office.shipmentHasLoadError || state.office.shipmentHasUpdateError,
    errorMessage: state.office.error,

    shipment: shipment,
    isUpdating: state.office.shipmentIsUpdating,

    // editablePanelify
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
