import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert';
import { formatCents } from 'shared/formatters';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';

import { getPpmIncentive, clearPpmIncentive } from './ducks';
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
      title: 'Weight (lbs)',
      'x-nullable': true,
      'x-always-required': true,
    },
  },
};
export class IncentiveCalculator extends Component {
  calculate = values => {
    const { planned_move_date, pickup_postal_code, destination_postal_code, weight } = values;
    this.props.getPpmIncentive(planned_move_date, pickup_postal_code, destination_postal_code, weight);
  };
  reset = async () => {
    const { reset, clearPpmIncentive } = this.props;
    await reset();
    clearPpmIncentive();
  };
  componentWillUnmount() {
    this.reset();
  }
  render() {
    const { handleSubmit, calculation, invalid, pristine, submitting, hasErrored } = this.props;
    return (
      <div className="calculator-panel incentive-calc">
        <div className="calculator-panel-title">Incentive Calculator</div>
        <form onSubmit={handleSubmit(this.calculate)}>
          <div className="usa-grid">
            {hasErrored && (
              <div className="usa-width-one-whole error-message">
                <Alert type="warning" heading="Could not perform calculation">
                  There was an issue calculating incentive.
                </Alert>
              </div>
            )}

            <div className="usa-width-one-half">
              <SwaggerField className="date-field" fieldName="planned_move_date" swagger={this.props.schema} required />
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
            {calculation && (
              <div className="usa-width-one-whole">
                <div className="calculated-result">
                  <table className="payment-table">
                    <tbody>
                      <tr className="payment-table-column-title">
                        <th colSpan="2">PPM Incentive</th>
                      </tr>
                    </tbody>
                    <tbody>
                      <tr>
                        <td>GCC</td>
                        <td align="right">
                          <span>${formatCents(calculation.gcc)}</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          <b>PPM Incentive @ 95%</b>
                        </td>
                        <td align="right">${formatCents(calculation.incentive_percentage)}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </form>
      </div>
    );
  }
}

IncentiveCalculator.propTypes = {
  schema: PropTypes.object.isRequired,
  getPpmIncentive: PropTypes.func.isRequired,
  error: PropTypes.object,
};

function mapStateToProps(state, ownProps) {
  const initialValues = pick(selectPPMForMove(state, ownProps.moveId), [
    'planned_move_date',
    'pickup_postal_code',
    'destination_postal_code',
  ]);
  const props = {
    schema,
    ...state.ppmIncentive,
    initialValues,
  };
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getPpmIncentive, clearPpmIncentive }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(reduxForm({ form: formName })(IncentiveCalculator));
