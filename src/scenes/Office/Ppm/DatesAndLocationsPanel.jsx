import { get } from 'lodash';
import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { reduxForm, getFormValues } from 'redux-form';

import { editablePanelify, PanelSwaggerField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectPPMForMove, updatePPM } from 'shared/Entities/modules/ppms';

const DatesAndLocationDisplay = props => {
  const fieldProps = {
    schema: props.ppmSchema,
    values: props.ppm,
  };
  return (
    <div className="editable-panel-column">
      <PanelSwaggerField fieldName="actual_move_date" title="Departure date" required {...fieldProps} />
      <PanelSwaggerField title="Pickup zip code" fieldName="pickup_postal_code" required {...fieldProps} />
      <PanelSwaggerField title="Destination zip code" fieldName="destination_postal_code" required {...fieldProps} />
    </div>
  );
};

const DatesAndLocationEdit = props => {
  const schema = props.ppmSchema;
  return (
    <div className="editable-panel-column">
      <SwaggerField className="short-field" fieldName="actual_move_date" title="Departure date" swagger={schema} />
      <SwaggerField
        className="short-field"
        title="Pickup zip code"
        fieldName="pickup_postal_code"
        swagger={schema}
        required
      />
      <SwaggerField
        className="short-field"
        title="Destination zip code"
        fieldName="destination_postal_code"
        swagger={schema}
        required
      />
    </div>
  );
};

const formName = 'ppm_dates_and_locations';

let DatesAndLocationPanel = editablePanelify(DatesAndLocationDisplay, DatesAndLocationEdit);
DatesAndLocationPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(DatesAndLocationPanel);

function mapStateToProps(state, props) {
  const formValues = getFormValues(formName)(state);
  const ppm = selectPPMForMove(state, props.moveId);

  return {
    // reduxForm
    formValues,
    initialValues: {
      actual_move_date: ppm.actual_move_date,
      pickup_postal_code: ppm.pickup_postal_code,
      destination_postal_code: ppm.destination_postal_code,
    },

    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    ppm,

    hasError: !!props.error,
    errorMessage: get(state, 'office.error'),
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      const values = getFormValues(formName)(state);
      return [props.moveId, ppm.id, values];
    },
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

export default connect(mapStateToProps, mapDispatchToProps)(DatesAndLocationPanel);
