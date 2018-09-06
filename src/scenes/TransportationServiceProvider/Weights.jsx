import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import { reduxForm, FormSection, getFormValues } from 'redux-form';

import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const weightsFields = ['actual_weight'];

const WeightsDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <div className="column-head">Weights</div>
        <div className="column-subhead">Total weight</div>
        <PanelSwaggerField
          fieldName="weight_estimate"
          nullWarning
          title="Customer estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="pm_survey_weight_estimate"
          nullWarning
          title="TSP estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="actual_weight"
          nullWarning
          {...fieldProps}
        />
        <div className="column-subhead">Pro-gear</div>
        <PanelSwaggerField
          fieldName="progear_weight_estimate"
          nullWarning
          title="Customer estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="pm_survey_progear_weight_estimate"
          nullWarning
          title="TSP estimate"
          {...fieldProps}
        />
        <div className="column-subhead">Spouse pro-gear</div>
        <PanelSwaggerField
          fieldName="spouse_progear_weight_estimate"
          nullWarning
          title="Customer estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="pm_survey_spouse_progear_weight_estimate"
          nullWarning
          title="TSP estimate"
          {...fieldProps}
        />
      </div>
    </React.Fragment>
  );
};

const WeightsEdit = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  const schema = props.shipmentSchema;
  return (
    <React.Fragment>
      <FormSection name="weights">
        <div className="editable-panel-column">
          <div className="column-head">Weights</div>
          <div className="column-subhead">Total weight</div>
          <PanelSwaggerField
            fieldName="weight_estimate"
            nullWarning
            title="Customer estimate"
            {...fieldProps}
          />
          <PanelSwaggerField
            fieldName="pm_survey_weight_estimate"
            nullWarning
            title="TSP estimate"
            {...fieldProps}
          />
          <SwaggerField fieldName="actual_weight" swagger={schema} required />
          <div className="column-subhead">Pro-gear</div>
          <PanelSwaggerField
            fieldName="progear_weight_estimate"
            nullWarning
            title="Customer estimate"
            {...fieldProps}
          />
          <PanelSwaggerField
            fieldName="pm_survey_progear_weight_estimate"
            nullWarning
            title="TSP estimate"
            {...fieldProps}
          />
          <div className="column-subhead">Spouse pro-gear</div>
          <PanelSwaggerField
            fieldName="spouse_progear_weight_estimate"
            nullWarning
            title="Customer estimate"
            {...fieldProps}
          />
          <PanelSwaggerField
            fieldName="pm_survey_spouse_progear_weight_estimate"
            nullWarning
            title="TSP estimate"
            {...fieldProps}
          />
        </div>
      </FormSection>
    </React.Fragment>
  );
};

const formName = 'shipment_weights';

let WeightsPanel = editablePanelify(WeightsDisplay, WeightsEdit);
WeightsPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(WeightsPanel);

WeightsPanel.propTypes = {
  shipment: PropTypes.object,
};

function mapStateToProps(state, props) {
  let formValues = getFormValues(formName)(state);

  return {
    // reduxForm
    formValues,
    initialValues: {
      weights: pick(props.shipment, weightsFields),
    },

    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),

    hasError: !!props.error,
    errorMessage: props.error,
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      return [get(props, 'shipment.id'), formValues.weights];
    },
  };
}

export default connect(mapStateToProps)(WeightsPanel);
