import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, FormSection } from 'redux-form';

import { PanelSwaggerField, PanelField, editablePanelify } from 'shared/EditablePanel';
import { formatCentsRange } from 'shared/formatters';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { selectPPMForMove, updatePPM } from 'shared/Entities/modules/ppms';

import { loadEntitlements } from 'scenes/Office/ducks';

const validateWeight = (value, formValues, props, fieldName) => {
  if (value && props.entitlement && value > props.entitlement.sum) {
    return `Cannot be more than full entitlement weight (${props.entitlement.sum} lbs)`;
  }
};

const EstimatesDisplay = props => {
  const ppm = props.PPMEstimate;
  const fieldProps = {
    schema: props.ppmSchema,
    values: ppm,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelField title="Incentive estimate">
          {formatCentsRange(ppm.incentive_estimate_min, ppm.incentive_estimate_max)}
        </PanelField>
        <PanelSwaggerField fieldName="weight_estimate" {...fieldProps} />
        <PanelSwaggerField title="Planned departure" fieldName="planned_move_date" {...fieldProps} />
        <PanelField title="Storage planned" fieldName="has_sit">
          {fieldProps.values.has_sit ? 'Yes' : 'No'}
        </PanelField>
        {fieldProps.values.has_sit && (
          <PanelSwaggerField title="Planned days in storage" fieldName="days_in_storage" {...fieldProps} />
        )}
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField title="Origin zip code" fieldName="pickup_postal_code" {...fieldProps} />
        <PanelSwaggerField title="Additional stop zip code" fieldName="additional_pickup_postal_code" {...fieldProps} />
        <PanelSwaggerField title="Destination zip code" fieldName="destination_postal_code" {...fieldProps} />
      </div>
    </React.Fragment>
  );
};

const EstimatesEdit = props => {
  const ppm = props.PPMEstimate;
  const schema = props.ppmSchema;

  return (
    <React.Fragment>
      <FormSection name="PPMEstimate">
        <div className="editable-panel-column">
          <PanelField title="Incentive estimate">
            {formatCentsRange(ppm.incentive_estimate_min, ppm.incentive_estimate_max)}
          </PanelField>
          <SwaggerField
            className="short-field"
            fieldName="weight_estimate"
            swagger={schema}
            validate={validateWeight}
            required
          />{' '}
          lbs
          <SwaggerField title="Planned departure date" fieldName="planned_move_date" swagger={schema} required />
          <div className="panel-subhead">Storage</div>
          <SwaggerField title="Storage planned?" fieldName="has_sit" swagger={schema} component={YesNoBoolean} />
          {get(props, 'formValues.PPMEstimate.has_sit', false) && (
            <SwaggerField title="Planned days in storage" fieldName="days_in_storage" swagger={schema} />
          )}
        </div>
        <div className="editable-panel-column">
          <SwaggerField title="Origin zip code" fieldName="pickup_postal_code" swagger={schema} required />
          <SwaggerField title="Additional stop zip code" fieldName="additional_pickup_postal_code" swagger={schema} />
          <SwaggerField title="Destination zip code" fieldName="destination_postal_code" swagger={schema} required />
        </div>
      </FormSection>
    </React.Fragment>
  );
};

const formName = 'ppm_estimate_and_details';

let PPMEstimatesPanel = editablePanelify(EstimatesDisplay, EstimatesEdit);
PPMEstimatesPanel = reduxForm({ form: formName })(PPMEstimatesPanel);

function mapStateToProps(state, ownProps) {
  const PPMEstimate = selectPPMForMove(state, ownProps.moveId);
  const formValues = getFormValues(formName)(state);

  return {
    // reduxForm
    formValues: formValues,
    initialValues: { PPMEstimate: PPMEstimate },

    // Wrapper
    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    hasError: false,
    errorMessage: get(state, 'office.error'),
    PPMEstimate: PPMEstimate,
    isUpdating: false,
    entitlement: loadEntitlements(state, ownProps.moveId),

    // editablePanelify
    getUpdateArgs: function() {
      if (
        formValues.PPMEstimate.additional_pickup_postal_code !== '' &&
        formValues.PPMEstimate.additional_pickup_postal_code !== undefined
      ) {
        formValues.PPMEstimate.has_additional_postal_code = true;
      } else {
        delete formValues.PPMEstimate.additional_pickup_postal_code;
        formValues.PPMEstimate.has_additional_postal_code = false;
      }
      return [ownProps.moveId, formValues.PPMEstimate.id, formValues.PPMEstimate];
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

export default connect(mapStateToProps, mapDispatchToProps)(PPMEstimatesPanel);
