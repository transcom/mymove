import { get, pick } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';
import { selectPPMForMove, updatePPM } from 'shared/Entities/modules/ppms';
import { selectAllDocumentsForMove } from 'shared/Entities/modules/moveDocuments';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';

const NetWeightDisplay = ({ ppmSchema, ppm, netWeight }) => {
  const fieldProps = {
    schema: ppmSchema,
    values: ppm,
  };
  return (
    <div className="editable-panel-column">
      <PanelSwaggerField title="Net Weight" fieldName="net_weight" required {...fieldProps}>
        {netWeight}
      </PanelSwaggerField>
    </div>
  );
};

const NetWeightEdit = ({ ppmSchema }) => {
  return (
    <div className="editable-panel-column net-weight">
      <SwaggerField className="short-field" title="Net Weight" fieldName="net_weight" swagger={ppmSchema} required />lbs
    </div>
  );
};

const formName = 'ppm_net_weight';

let NetWeightPanel = editablePanelify(NetWeightDisplay, NetWeightEdit, false);
NetWeightPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(NetWeightPanel);

function mapStateToProps(state, ownProps) {
  const formValues = getFormValues(formName)(state);
  const ppm = selectPPMForMove(state, ownProps.moveId);
  const moveDocs = selectAllDocumentsForMove(state, ownProps.moveId);
  const netWeight = moveDocs.reduce((accum, { move_document_type, status, full_weight, empty_weight }) => {
    if (move_document_type === 'WEIGHT_TICKET_SET' && status === 'OK') {
      return accum + (full_weight - empty_weight);
    }
    return accum;
  }, 0);

  return {
    // reduxForm
    initialValues: pick(ppm, 'net_weight'),
    formValues,

    // Wrapper
    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    ppm,
    // editablePanelify
    getUpdateArgs: function() {
      const values = getFormValues(formName)(state);
      return [ownProps.moveId, ppm.id, values];
    },
    netWeight,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updatePPM,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(NetWeightPanel);
