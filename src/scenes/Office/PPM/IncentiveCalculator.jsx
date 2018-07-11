import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPpmWeightEstimate as getPpmIncentive } from '../../Moves/Ppm/ducks';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const formName = 'ppm_reimbursement_calc';
const schema = {
  properties: {
    planned_move_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Move Date',
      'x-nullable': true,
      'x-always-required': true,
    },
    pickup_postal_code: {
      type: 'string',
      format: 'zip',
      title: 'Origin ZIP',
      example: '90210',
      pattern: '^(\\d{5}([\\-]\\d{4})?)$',
      'x-nullable': true,
      'x-always-required': true,
    },
    destination_postal_code: {
      type: 'string',
      format: 'zip',
      title: 'Destination ZIP',
      example: '90210',
      pattern: '^(\\d{5}([\\-]\\d{4})?)$',
      'x-nullable': true,
      'x-always-required': true,
    },
    weight: {
      type: 'integer',
      minimum: 1,
      title: 'Weight',
      'x-nullable': true,
      'x-always-required': true,
    },
  },
};
export class IncentiveCalculator extends Component {
  calculate = values => {
    const {
      planned_move_date,
      pickup_postal_code,
      destination_postal_code,
      weight,
    } = values;
    this.props.getPpmIncentive(
      planned_move_date,
      pickup_postal_code,
      destination_postal_code,
      weight,
    );
  };

  render() {
    const {
      handleSubmit,
      reimbursement,
      invalid,
      pristine,
      reset,
      submitting,
    } = this.props;
    return (
      <div className="calculator-panel">
        <div className="calculator-panel-title">Incentive Calculator</div>
        <form onSubmit={handleSubmit(this.calculate)}>
          <SwaggerField
            className="date-field"
            fieldName="planned_move_date"
            swagger={this.props.schema}
            required
          />
          <SwaggerField
            className="short-field"
            fieldName="pickup_postal_code"
            swagger={this.props.schema}
            required
          />
          <SwaggerField
            className="short-field"
            fieldName="destination_postal_code"
            swagger={this.props.schema}
            required
          />
          <SwaggerField
            className="short-field"
            fieldName="weight"
            swagger={this.props.schema}
            required
          />
          <div className="buttons">
            <button
              data-cy="calc"
              type="submit"
              disabled={pristine || submitting || invalid}
            >
              Calculate
            </button>
            <button
              className="usa-button-secondary"
              data-cy="reset"
              type="button"
              disabled={pristine || submitting}
              onClick={reset}
            >
              Reset
            </button>
          </div>
        </form>
        {reimbursement && (
          <div className="calculated-result">
            Maximum Obligation: <b>{reimbursement}</b>
          </div>
        )}
      </div>
    );
  }
}

IncentiveCalculator.propTypes = {
  schema: PropTypes.object.isRequired,
  getPpmIncentive: PropTypes.func.isRequired,
  error: PropTypes.object,
};

function mapStateToProps(state) {
  const initialValues = pick(get(state, 'office.officePPMs[0]'), [
    'planned_move_date',
    'pickup_postal_code',
    'destination_postal_code',
  ]);
  const props = {
    schema,
    hasEstimateError: state.ppm.hasEstimateError,
    reimbursement: state.ppm.incentive_estimate_min,
    initialValues,
  };
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getPpmIncentive }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(
  reduxForm({ form: formName })(IncentiveCalculator),
);
