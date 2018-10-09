import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import { reduxForm, FormSection, getFormValues } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import {
  PanelSwaggerField,
  editablePanelify,
  PanelField,
} from 'shared/EditablePanel';

const weightsFields = [
  'net_weight',
  'gross_weight',
  'tare_weight',
  'pm_survey_weight_estimate',
  'pm_survey_progear_weight_estimate',
  'pm_survey_spouse_progear_weight_estimate',
];

const WeightsDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  const values = props.shipment;
  return (
    <Fragment>
      <div className="editable-panel-column">
        <div className="column-subhead">Weight</div>
        <PanelSwaggerField
          fieldName="weight_estimate"
          required
          title="Customer estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="pm_survey_weight_estimate"
          required
          title="TSP estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Actual"
          fieldName="net_weight"
          required
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-column">
        <div className="column-subhead">Pro-gear (Service member + spouse)</div>
        <PanelField title="Customer estimate">
          {values.progear_weight_estimate ? (
            <span>{values.progear_weight_estimate} lbs</span>
          ) : (
            '0 lbs'
          )}
          {values.spouse_progear_weight_estimate ? (
            <span> + {values.spouse_progear_weight_estimate} lbs</span>
          ) : (
            ' + 0 lbs'
          )}
        </PanelField>
        <PanelField title="TSP estimate">
          {values.pm_survey_progear_weight_estimate ? (
            <span>{values.pm_survey_progear_weight_estimate} lbs</span>
          ) : (
            '0 lbs'
          )}
          {values.pm_survey_spouse_progear_weight_estimate ? (
            <span>
              {' '}
              + {values.pm_survey_spouse_progear_weight_estimate} lbs
            </span>
          ) : (
            ' + 0 lbs'
          )}
        </PanelField>
      </div>
    </Fragment>
  );
};

const WeightsEdit = props => {
  const schema = props.shipmentSchema;
  const fieldProps = {
    schema,
    values: props.shipment,
  };
  const values = props.shipment;
  return (
    <Fragment>
      <FormSection name="weights">
        <div className="editable-panel-column">
          <div className="column-head">Weights</div>
          <PanelSwaggerField
            fieldName="weight_estimate"
            required
            title="Customer estimate"
            {...fieldProps}
          />
          <PanelSwaggerField
            fieldName="pm_survey_weight_estimate"
            required
            title="TSP estimate"
            {...fieldProps}
          />
          <div className="column-subhead">Actual Weights</div>
          <SwaggerField
            className="short-field"
            fieldName="gross_weight"
            swagger={schema}
            required
          />{' '}
          lbs
          <SwaggerField
            className="short-field"
            fieldName="tare_weight"
            swagger={schema}
            required
          />{' '}
          lbs
          <SwaggerField
            title="Net (Gross - Tare)"
            className="short-field"
            fieldName="net_weight"
            swagger={schema}
            required
          />{' '}
          lbs
        </div>
        <div className="editable-panel-column">
          <div className="column-head">Pro-gear (Service member + spouse)</div>
          <PanelField title="Customer estimate">
            {values.progear_weight_estimate ? (
              <span>{values.progear_weight_estimate} lbs</span>
            ) : (
              '0 lbs'
            )}
            {values.spouse_progear_weight_estimate ? (
              <span> + {values.spouse_progear_weight_estimate} lbs</span>
            ) : (
              ' + 0 lbs'
            )}
          </PanelField>
          <div className="column-subhead">TSP Estimate</div>
          <SwaggerField
            className="short-field"
            fieldName="pm_survey_progear_weight_estimate"
            title="Service member"
            swagger={schema}
            required
          />{' '}
          lbs
          <SwaggerField
            className="short-field"
            fieldName="pm_survey_spouse_progear_weight_estimate"
            title="Spouse"
            swagger={schema}
            required
          />{' '}
          lbs
        </div>
      </FormSection>
    </Fragment>
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
  schema: PropTypes.object,
};

function mapStateToProps(state, props) {
  let formValues = getFormValues(formName)(state);

  return {
    // reduxForm
    formValues,
    initialValues: {
      weights: pick(props.shipment, weightsFields),
    },

    shipmentSchema: get(state, 'swaggerPublic.spec.definitions.Shipment', {}),

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
