// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import isMobile from 'is-mobile';
import { Link } from 'react-router-dom';
import { concat, get, reject, every, includes } from 'lodash';
import {
  reduxForm,
  Field,
  FormSection,
  getFormValues,
  isValid,
} from 'redux-form';

import {
  PanelSwaggerField,
  PanelField,
  SwaggerValue,
  editablePanelify,
} from 'shared/EditablePanel';
import { formatDate } from 'shared/formatters';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

import './index.css';

const SurveyDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-3-column">
        <PanelSwaggerField
          title="Pack Date"
          fieldName="pm_survey_pack_date"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Pickup Date"
          fieldName="pm_survey_pickup_date"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Latest Pickup Date"
          fieldName="pm_survey_latest_pickup_date"
          nullWarning
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-3-column">
        <PanelSwaggerField
          title="Earliest Pickup Date"
          fieldName="pm_survey_earliest_delivery_date"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Latest Delivery Date"
          fieldName="pm_survey_latest_delivery_date"
          nullWarning
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-3-column">
        <PanelSwaggerField
          title="Weight Estimate"
          fieldName="pm_survey_weight_estimate"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Progear Weight Estimate"
          fieldName="pm_survey_progear_weight_estimate"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Spouse Progear Weight Estimate"
          fieldName="pm_survey_spouse_progear_weight_estimate"
          nullWarning
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-column">
        <PanelSwaggerField
          title="Notes"
          fieldName="pm_survey_notes"
          nullWarning
          {...fieldProps}
        />
      </div>
    </React.Fragment>
  );
};

const SurveyEdit = props => {
  const schema = props.shipmentSchema;
  return (
    <React.Fragment>
      <FormSection name="survey">
        <div className="editable-panel-3-column">
          <SwaggerField
            fieldName="pm_survey_pack_date"
            swagger={schema}
            className="half-width"
            required
          />
          <SwaggerField
            fieldName="pm_survey_pickup_date"
            swagger={schema}
            className="half-width"
            required
          />
          <SwaggerField
            fieldName="pm_survey_latest_pickup_date"
            swagger={schema}
            className="half-width"
            required
          />
        </div>

        <div className="editable-panel-3-column">
          <SwaggerField
            fieldName="pm_survey_earliest_delivery_date"
            swagger={schema}
            className="half-width"
            required
          />
          <SwaggerField
            fieldName="pm_survey_latest_delivery_date"
            swagger={schema}
            className="half-width"
            required
          />
        </div>

        <div className="editable-panel-3-column">
          <SwaggerField
            fieldName="pm_survey_weight_estimate"
            swagger={schema}
            className="half-width"
            required
          />
          <SwaggerField
            fieldName="pm_survey_progear_weight_estimate"
            swagger={schema}
            className="half-width"
            required
          />
          <SwaggerField
            fieldName="pm_survey_spouse_progear_weight_estimate"
            swagger={schema}
            className="half-width"
            required
          />
        </div>
        <SwaggerField
          fieldName="pm_survey_notes"
          swagger={schema}
          className="half-width"
        />
      </FormSection>
    </React.Fragment>
  );
};

const formName = 'shipment_pre_move_survey';

let PremoveSurveyPanel = editablePanelify(SurveyDisplay, SurveyEdit);
PremoveSurveyPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PremoveSurveyPanel);

PremoveSurveyPanel.propTypes = {
  shipment: PropTypes.object,
};

function mapStateToProps(state, props) {
  let formValues = getFormValues(formName)(state);

  return {
    // reduxForm
    formValues: formValues,
    initialValues: {
      survey: props.shipment,
    },

    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),

    hasError: false,
    // errorMessage: get(state, 'office.error'),
    isUpdating: false,

    // editablePanelify
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      return [get(props, 'shipment.id'), formValues.survey];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PremoveSurveyPanel);
