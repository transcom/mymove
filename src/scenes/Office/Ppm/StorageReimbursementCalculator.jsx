import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';

import { getPpmSitEstimate, clearPpmSitEstimate } from '../../Moves/Ppm/ducks';

const formName = 'storage_reimbursement_calc';
const schema = {
  properties: {
    move_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Move Date',
      'x-nullable': true,
      'x-always-required': false,
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
    days_in_storage: {
      type: 'integer',
      title: 'Days in Storage',
      minimum: 0,
      maximum: 90,
      'x-nullable': true,
      'x-always-required': true,
    },
    weight: {
      type: 'integer',
      minimum: 1,
      title: 'Weight (lbs)',
      'x-nullable': true,
      'x-always-required': true,
    },
  },
};
export class StorageReimbursementCalculator extends Component {
  reset = async () => {
    const { reset, clearPpmSitEstimate } = this.props;
    await reset();
    clearPpmSitEstimate();
  };
  componentWillUnmount() {
    this.reset();
  }
  calculate = values => {
    const { pickup_postal_code, destination_postal_code, days_in_storage, weight, move_date } = values;
    this.props.getPpmSitEstimate(move_date, days_in_storage, pickup_postal_code, destination_postal_code, weight);
  };

  render() {
    const { handleSubmit, sitReimbursement, invalid, pristine, submitting, hasEstimateError } = this.props;

    return (
      <div className="calculator-panel storage-calc">
        <div className="calculator-panel-title">Storage Calculator</div>
        <form onSubmit={handleSubmit(this.calculate)}>
          <div className="usa-grid">
            {hasEstimateError && (
              <div className="usa-width-one-whole error-message">
                <Alert type="warning" heading="Could not retrieve estimate">
                  There was an issue retrieving reimbursement amount.
                </Alert>
              </div>
            )}
            <div className="usa-width-one-half">
              <SwaggerField className="date-field" fieldName="move_date" swagger={this.props.schema} required />
              <SwaggerField className="short-field" fieldName="weight" swagger={this.props.schema} required />
            </div>
            <div className="usa-width-one-half">
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
              <SwaggerField className="short-field" fieldName="days_in_storage" swagger={this.props.schema} required />
            </div>
          </div>
          <div className="usa-grid">
            <div className="usa-width-one-whole">
              <div className="buttons">
                <button data-cy="calc" type="submit" disabled={pristine || submitting || invalid}>
                  Calculate
                </button>
                <button
                  className="usa-button-secondary"
                  data-cy="reset"
                  type="button"
                  disabled={pristine || submitting}
                  onClick={this.reset}
                >
                  Reset
                </button>
              </div>
            </div>
            {sitReimbursement && (
              <div className="usa-width-one-whole">
                <div className="calculated-result">
                  Maximum Obligation: <b>{sitReimbursement}</b>
                </div>
              </div>
            )}
          </div>
        </form>
      </div>
    );
  }
}

StorageReimbursementCalculator.propTypes = {
  schema: PropTypes.object.isRequired,
  getPpmSitEstimate: PropTypes.func.isRequired,
  error: PropTypes.object,
};

function mapStateToProps(state, ownProps) {
  let ppm = selectPPMForMove(state, ownProps.moveId);
  let initialValues = pick(ppm, ['pickup_postal_code', 'destination_postal_code', 'days_in_storage']);
  initialValues.move_date = ppm.actual_move_date || ppm.original_move_date;
  return {
    schema,
    hasEstimateError: state.ppm.hasEstimateError,
    sitReimbursement: state.ppm.sitReimbursement,
    initialValues,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getPpmSitEstimate, clearPpmSitEstimate }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(
  reduxForm({ form: formName, enableReinitialize: true, keepDirtyOnReinitialize: true })(
    StorageReimbursementCalculator,
  ),
);
