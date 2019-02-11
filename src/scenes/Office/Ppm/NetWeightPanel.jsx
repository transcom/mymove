import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';

import { updateOrders } from '../ducks';

const NetWeightDisplay = props => {
  const fieldProps = {
    schema: props.ppmSchema,
    values: props.ppm,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField title="Net Weight" fieldName="net_weight" required {...fieldProps} />
      </div>
    </React.Fragment>
  );
};

const NetWeightEdit = props => {
  const { ppmSchema } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <SwaggerField title="Net Weight" fieldName="net_weight" swagger={ppmSchema} required />
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_accounting';

let NetWeightPanel = editablePanelify(NetWeightDisplay, NetWeightEdit);
NetWeightPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(NetWeightPanel);

function mapStateToProps(state, ownProps) {
  const formValues = getFormValues(formName)(state);
  const ppm = selectPPMForMove(state, ownProps.moveId);
  return {
    // reduxForm
    initialValues: ppm,
    formValues,

    // Wrapper
    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    ppm,
    // editablePanelify
    // getUpdateArgs: function() {
    //   let values = getFormValues(formName)(state);
    //   values.new_duty_station_id = values.new_duty_station.id;
    //   return [orders.id, values];
    // },
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

export default connect(mapStateToProps, mapDispatchToProps)(NetWeightPanel);
