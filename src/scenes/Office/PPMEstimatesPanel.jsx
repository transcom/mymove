import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from './editablePanel';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import { createOrUpdatePpm } from 'scenes/Moves/Ppm/ducks';

const EstimatesDisplay = props => {
  const fieldProps = {
    schema: props.PPMEstimateSchema,
    values: props.PPMEstimate,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="estimated_incentive" {...fieldProps} />
        <PanelSwaggerField fieldName="weight_estimate" {...fieldProps} />
        <PanelSwaggerField
          title="Planned departure"
          fieldName="planned_move_date"
          {...fieldProps}
        />
        <PanelField title="Storage planned" fieldName="days_in_storage">
          {fieldProps.values.has_sit ? 'Yes' : 'No'}
        </PanelField>
        <PanelSwaggerField
          title="Storage days"
          fieldName="days_in_storage"
          {...fieldProps}
        />
        <PanelField
          title="Max. storage cost"
          value="Max. storage cost"
          className="Todo"
        />
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField
          title="Origin zip code"
          fieldName="pickup_postal_code"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Additional stop zip code"
          fieldName="additional_pickup_postal_code"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Destination zip code"
          fieldName="destination_postal_code"
          {...fieldProps}
        />
        <PanelField
          title="Distance estimate"
          value="863 miles"
          className="Todo"
        />
      </div>
    </React.Fragment>
  );
};

const EstimatesEdit = props => {
  const schema = props.PPMEstimateSchema;
  return (
    <React.Fragment>
      <FormSection name="PPMEstimate">
        <div className="editable-panel-column">
          <SwaggerField fieldName="estimated_incentive" swagger={schema} />
          <SwaggerField fieldName="weight_estimate" swagger={schema} />
          <SwaggerField
            title="Planned departure"
            fieldName="planned_move_date"
            swagger={schema}
          />
          <SwaggerField
            title="Storage planned"
            fieldName="days_in_storage"
            swagger={schema}
          />
          <SwaggerField
            title="Storage days"
            fieldName="days_in_storage"
            swagger={schema}
          />
          <SwaggerField
            title="Max. storage cost"
            swagger={schema}
            className="Todo"
          />
        </div>
        <div className="editable-panel-column">
          <SwaggerField
            title="Origin zip code"
            fieldName="pickup_postal_code"
            swagger={schema}
          />
          <SwaggerField
            title="Additional stop zip code"
            fieldName="additional_pickup_postal_code"
            swagger={schema}
          />
          <SwaggerField
            title="Destination zip code"
            fieldName="destination_postal_code"
            swagger={schema}
          />
          {/*<SwaggerField
          title="Distance estimate"
          fieldName="destination_postal_code"
          value="863 miles"
          className="Todo"
        />*/}
        </div>
      </FormSection>
    </React.Fragment>
  );
};

const formName = 'ppm_estimate_and_details';

let PPMEstimatesPanel = editablePanel(EstimatesDisplay, EstimatesEdit);
PPMEstimatesPanel = reduxForm({ form: formName })(PPMEstimatesPanel);

function mapStateToProps(state) {
  let PPMEstimate = get(state, 'office.officePPMs[0]', {});
  let officeMove = get(state, 'office.officeMove', {});

  return {
    // reduxForm
    formData: state.form[formName],
    initialValues: { PPMEstimate: PPMEstimate },

    // Wrapper
    PPMEstimateSchema: get(
      state,
      'swagger.spec.definitions.PersonallyProcuredMovePayload',
    ),
    hasError: false,
    errorMessage: get(state, 'office.error'),
    PPMEstimate: PPMEstimate,
    isUpdating: false,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      return [officeMove.id, values.PPMEstimate.id, values.PPMEstimate];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: createOrUpdatePpm,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PPMEstimatesPanel);
