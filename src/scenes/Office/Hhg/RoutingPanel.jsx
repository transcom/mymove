import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import { no_op } from 'shared/utils';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelField, PanelSwaggerField, SwaggerValue, editablePanelify } from 'shared/EditablePanel';
import { selectShipmentForMove } from 'shared/Entities/modules/shipments';

// TODO: add shipment_type
// TODO: combine source_gbloc and destination_gbloc values to one value
const RoutingPanelDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  const tdlFieldProps = {
    schema: props.tdlSchema,
    values: props.shipment.traffic_distribution_list,
  };
  return (
    <div className="editable-panel-column">
      <PanelField title="Shipment type">HHG</PanelField>
      <PanelSwaggerField title="Shipment market" fieldName="market" {...fieldProps} />
      <PanelField title="Channel">
        <SwaggerValue fieldName="source_gbloc" {...fieldProps} /> -{' '}
        <SwaggerValue fieldName="destination_gbloc" {...fieldProps} />
      </PanelField>
      <PanelSwaggerField title="Code of service" fieldName="code_of_service" {...tdlFieldProps} />
    </div>
  );
};

const RoutingPanelEdit = props => {
  const { shipmentSchema, tdlSchema } = props;
  return (
    <div className="editable-panel-column">
      <PanelField title="Shipment type">HHG</PanelField>
      <PanelField title="Market">dHHG</PanelField>
      <SwaggerField title="Source GBLOC" fieldName="source_gbloc" swagger={shipmentSchema} required />
      <SwaggerField title="Destination GBLOC" fieldName="source_gbloc" swagger={shipmentSchema} required />
      <SwaggerField title="Code of service" fieldName="code_of_service" swagger={tdlSchema} required />
    </div>
  );
};

const formName = 'office_shipment_routing';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let RoutingPanel = editablePanelify(RoutingPanelDisplay, RoutingPanelEdit, editEnabled);
RoutingPanel = reduxForm({ form: formName })(RoutingPanel);

function mapStateToProps(state, ownProps) {
  const { moveId } = ownProps;
  const shipment = selectShipmentForMove(state, moveId);

  return {
    initialValues: shipment,
    // Wrapper
    tdlSchema: get(state, 'swaggerInternal.spec.definitions.TrafficDistributionList'),
    shipmentSchema: get(state, 'swaggerInternal.spec.definitions.Shipment', {}),
    shipment: shipment,

    // editablePanelify
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
