import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { reduxForm } from 'redux-form';

let PremoveSurveyForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form className="infoPanel-wizard" onSubmit={handleSubmit}>
      <div className="infoPanel-wizard-header">Pre-move Survey</div>
      <div className="editable-panel-column">
        <div className="column-subhead">Dates</div>
        <SwaggerField
          fieldName="pm_survey_conducted_date"
          title="Pre-move survey conducted"
          swagger={schema}
          required
        />
        <SwaggerField fieldName="pm_survey_method" swagger={schema} required />
        <SwaggerField
          fieldName="pm_survey_planned_pack_date"
          swagger={schema}
          title="Planned packing (first day)"
          required
        />
        <SwaggerField fieldName="pm_survey_planned_pickup_date" swagger={schema} title="Planned pickup" required />
        <SwaggerField fieldName="pm_survey_planned_delivery_date" swagger={schema} title="Planned delivery" required />
        <SwaggerField fieldName="pm_survey_notes" title="Notes about dates" swagger={schema} />
      </div>

      <div className="editable-panel-column">
        <div className="column-subhead">Weights</div>
        <SwaggerField className="short-field" fieldName="pm_survey_weight_estimate" swagger={schema} required />
        <SwaggerField className="short-field" fieldName="pm_survey_progear_weight_estimate" swagger={schema} />
        <SwaggerField className="short-field" fieldName="pm_survey_spouse_progear_weight_estimate" swagger={schema} />
      </div>

      <div className="infoPanel-wizard-actions-container">
        <a className="infoPanel-wizard-cancel" onClick={onCancel}>
          Cancel
        </a>
        <button type="submit" disabled={submitting || !valid}>
          Done
        </button>
      </div>
    </form>
  );
};

PremoveSurveyForm = reduxForm({ form: 'shipment_pre_move_survey' })(PremoveSurveyForm);

PremoveSurveyForm.propTypes = {
  schema: PropTypes.object,
  onCancel: PropTypes.func,
  handleSubmit: PropTypes.func,
  submitting: PropTypes.bool,
  valid: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const { shipment } = state.tsp;
  return {
    ...ownProps,
    initialValues: shipment,
  };
};

export default connect(mapStateToProps)(PremoveSurveyForm);
