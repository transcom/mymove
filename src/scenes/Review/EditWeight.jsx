import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { debounce, get } from 'lodash';

import { push } from 'react-router-redux';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import {
  createOrUpdatePpm,
  loadPpm,
  getPpmWeightEstimate,
} from 'scenes/Moves/Ppm/ducks';
import { loadEntitlements } from 'scenes/Orders/ducks';

import './Review.css';
import './EditWeight.css';
import profileImage from './images/profile.png';

const editWeightFormName = 'edit_weight';
const weightEstimateDebounce = 300;
const examples = [
  {
    weight: 100,
    incentive: '$0 - 100',
    description: 'A few items in your car',
  },
  { weight: 400, incentive: '$300 - 400' },
  { weight: 600, incentive: '$500 - 600' },
  {
    weight: 1000,
    incentive: '$800 - 1,000',
    description: 'A trailer full of household goods',
  },
  { weight: 2000, incentive: '$1,500 - 1,800' },
  {
    weight: 5000,
    incentive: '$3,100 - 3,700',
    description: 'Renting a moving truck yourself',
  },
  { weight: 10000, incentive: '$5,900 - 6,800' },
];

const validateWeight = (value, formValues, props, fieldName) => {
  if (value && props.entitlement && value > props.entitlement.sum) {
    return 'Cannot be more than your full entitlement';
  }
};

let EditWeightForm = props => {
  const {
    onCancel,
    schema,
    handleSubmit,
    submitting,
    valid,
    entitlement,
    dirty,
    incentive,
    onWeightChange,
    initialValues,
  } = props;

  let incentiveMax;
  if (incentive) {
    // Is of the form "$500-700"
    incentiveMax = Number(incentive.split('-')[1]);
  }

  // Error class if below advance amount, otherwise warn class if incentive has changed
  let incentiveClass = '';
  let fieldClass = dirty ? 'warn' : '';
  let advanceError = false;
  const advanceAmt = get(initialValues, 'advance.requested_amount', 0);
  if (incentiveMax && advanceAmt && incentiveMax < advanceAmt) {
    advanceError = true;
    incentiveClass = 'error';
    fieldClass = 'error';
  } else if (get(initialValues, 'estimated_incentive') !== incentive) {
    incentiveClass = 'warn';
  }

  const fullFieldClass = `weight-estimate-input ${fieldClass}`;
  return (
    <form onSubmit={handleSubmit}>
      <img src={profileImage} alt="" /> Profile
      <hr />
      <h3 className="sm-heading">Edit PPM Weight:</h3>
      <p>
        Changes could impact your move, including the estimated PPM incentive.
      </p>
      {entitlement && (
        <div className="edit-ppm-entitlement">
          <p>
            <strong>How much are you entitled to move?</strong>
          </p>
          <p>
            {entitlement.total.toLocaleString()} lbs. +{' '}
            {entitlement.pro_gear.toLocaleString()} lbs. of pro-gear +{' '}
            {entitlement.pro_gear_spouse.toLocaleString()} lbs. of spouse's
            pro-gear
          </p>
        </div>
      )}
      <div className="edit-weight-container">
        <div className="usa-width-one-half">
          <h4 className="sm-heading">Move estimate</h4>
          <div>
            <SwaggerField
              className={fullFieldClass}
              fieldName="weight_estimate"
              swagger={schema}
              onChange={onWeightChange}
              validate={validateWeight}
              required
            />
            <span> lbs</span>
          </div>
          <div>
            {!advanceError &&
              dirty && (
                <div className="usa-alert usa-alert-warning">
                  <div className="usa-alert-body">
                    <p className="usa-alert-text">
                      This update will change your incentive.
                    </p>
                  </div>
                </div>
              )}
            {advanceError && (
              <p className="advance-error">
                Weight is too low and will require paying back the advance.
              </p>
            )}
          </div>

          <div className="display-value">
            <p>Estimated Incentive</p>
            <p className={incentiveClass}>
              <strong>{incentive}</strong>
            </p>
            {initialValues &&
              initialValues.estimated_incentive !== incentive && (
                <p className="subtext">
                  Originally {initialValues.estimated_incentive}
                </p>
              )}
          </div>

          {get(initialValues, 'has_requested_advance') && (
            <div className="display-value">
              <p>Advance</p>
              <p>
                <strong>${initialValues.advance.requested_amount}</strong>
              </p>
            </div>
          )}
        </div>

        <div className="usa-width-one-half">
          <h4 className="sm-heading">Examples</h4>
          <table className="examples-table">
            <thead>
              <tr>
                <th>Weight</th>
                <th>Incentive</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {examples.map(ex => (
                <tr key={ex.weight}>
                  <td>{ex.weight.toLocaleString()}</td>
                  <td>{ex.incentive}</td>
                  <td>{ex.description || ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
      <button type="submit" disabled={submitting || !valid}>
        Save
      </button>
      <button
        type="button"
        className="usa-button-secondary"
        disabled={submitting}
        onClick={onCancel}
      >
        Cancel
      </button>
    </form>
  );
};
EditWeightForm = reduxForm({
  form: editWeightFormName,
})(EditWeightForm);

class EditWeight extends Component {
  componentDidUpdate(prevProps, prevState) {
    if (
      !prevProps.loggedInUser.hasSucceeded &&
      this.props.loggedInUser.hasSucceeded
    ) {
      this.props.loadPpm(this.props.match.params.moveId);
    }
  }

  returnToReview = () => {
    const reviewAddress = `/moves/${this.props.match.params.moveId}/review`;
    this.props.push(reviewAddress);
  };

  debouncedGetPpmWeightEstimate = debounce(
    this.props.getPpmWeightEstimate,
    weightEstimateDebounce,
  );

  onWeightChange = (e, newValue, oldValue, fieldName) => {
    const { currentPpm, entitlement } = this.props;
    if (newValue > 0 && newValue <= entitlement.sum) {
      this.debouncedGetPpmWeightEstimate(
        currentPpm.planned_move_date,
        currentPpm.pickup_postal_code,
        currentPpm.destination_postal_code,
        newValue,
      );
    } else {
      this.debouncedGetPpmWeightEstimate.cancel();
    }
  };

  updatePpm = (values, dispatch, props) => {
    const moveId = this.props.match.params.moveId;
    return this.props
      .createOrUpdatePpm(moveId, {
        weight_estimate: values.weight_estimate,
        estimated_incentive: props.incentive,
      })
      .then(() => {
        // This promise resolves regardless of error.
        if (!this.props.hasSubmitError) {
          this.returnToReview();
        } else {
          window.scrollTo(0, 0);
        }
      });
  };

  render() {
    const { error, schema, currentPpm, entitlement, incentive } = this.props;

    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditWeightForm
            initialValues={currentPpm}
            incentive={incentive}
            onSubmit={this.updatePpm}
            onCancel={this.returnToReview}
            onWeightChange={this.onWeightChange}
            entitlement={entitlement}
            schema={schema}
          />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    ...state.ppm,
    loggedInUser: get(state, 'loggedInUser'),
    error: get(state, 'serviceMember.error') || state.ppm.hasEstimateError,
    hasSubmitError: get(state, 'serviceMember.hasSubmitError'),
    entitlement: loadEntitlements(state),
    schema: get(
      state,
      'swagger.spec.definitions.UpdatePersonallyProcuredMovePayload',
      {},
    ),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { push, createOrUpdatePpm, getPpmWeightEstimate, loadPpm },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(EditWeight);
