import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import { reduxForm, FormSection, getFormValues } from 'redux-form';

import { PanelSwaggerField, editablePanelify } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const weightsFields = [
  'weight_estimate',
  'progear_weight_estimate',
  'spouse_progear_weight_estimate',
];

const WeightsDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField
          fieldName="weight_estimate"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="progear_weight_estimate"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="spouse_progear_weight_estimate"
          nullWarning
          {...fieldProps}
        />
      </div>
    </React.Fragment>
  );
};

const WeightsEdit = props => {
  const schema = props.shipmentSchema;
  return (
    <React.Fragment>
      <FormSection name="weights">
        <div className="editable-panel-column">
          <SwaggerField fieldName="weight_estimate" swagger={schema} required />
          <SwaggerField
            fieldName="progear_weight_estimate"
            swagger={schema}
            required
          />
          <SwaggerField
            fieldName="spouse_progear_weight_estimate"
            swagger={schema}
            required
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

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(WeightsPanel);
