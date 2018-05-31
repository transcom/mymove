import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import Slider from 'react-rangeslider'; //todo: pull from node_modules, override
import { Field } from 'redux-form';

import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { loadEntitlements } from 'scenes/Orders/ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import Alert from 'shared/Alert';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  setPendingPpmWeight,
  getPpmWeightEstimate,
  createOrUpdatePpm,
  getPpmMaxWeightEstimate,
} from './ducks';

import 'react-rangeslider/lib/index.css';
import './Weight.css';

function getWeightInfo(ppm, entitlement) {
  const size = ppm ? ppm.size : 'L';
  switch (size) {
    case 'S':
      return {
        min: 50,
        max: 1000,
      };
    case 'M':
      return {
        min: 500,
        max: 2500,
      };
    default:
      return {
        min: 1500,
        max: entitlement.sum,
      };
  }
}

const requestedTitle = maxAdvance => {
  return (
    <Fragment>
      <div className="ppmquestion">How much advance do you want?</div>
      <div className="ppmmuted">
        Up to {maxAdvance} (60% of your PPM incentive)
      </div>
    </Fragment>
  );
};

const methodTitle = (
  <Fragment>
    <div className="ppmquestion">How do you want to get your advance?</div>
    <div className="ppmmuted">
      To direct deposit to another account you'll need to fill out a new account
      form, included in your advance paperwork, and take it to the accounting
      office.
    </div>
  </Fragment>
);

const formatMaxAdvance = maxAdvance => {
  return `$${maxAdvance.toFixed(2)}`;
};

const validateAdvanceForm = (values, form) => {
  if (values.hasEstimateInProgress) {
    return { has_requested_advance: 'Esimate in progress.' };
  }

  if (values.maxIncentive) {
    if (parseFloat(values.requested_amount) > parseFloat(values.maxIncentive)) {
      return {
        requested_amount: `Must be less than ${formatMaxAdvance(
          values.maxIncentive,
        )}`,
      };
    }
  }
};

const requestAdvanceFormName = 'request_advance';
class RequestAdvanceForm extends Component {
  state = { showInfo: false };

  openInfo = () => {
    this.setState({ showInfo: true });
  };
  closeInfo = () => {
    this.setState({ showInfo: false });
  };

  render() {
    const { schema, hasRequestedAdvance, maxIncentive } = this.props;
    let maxAdvance = '';
    if (maxIncentive) {
      maxAdvance = formatMaxAdvance(maxIncentive);
    }
    return (
      <div className="whole_box">
        <div>
          <div className="usa-width-one-whole">
            <div className="usa-width-two-thirds">
              <div className="ppmquestion">
                Would you like an advance of up to 60% of your PPM incentive? ({
                  maxAdvance
                })
              </div>
              <div className="ppmmuted">
                We recommend paying for expenses with your government travel
                card, rather than getting an advance.{' '}
                <a onClick={this.openInfo}>Why?</a>
              </div>
            </div>
            <div className="usa-width-one-third">
              <Field name="has_requested_advance" component={YesNoBoolean} />
            </div>
          </div>
          {this.state.showInfo && (
            <div className="usa-width-one-whole top-buffered">
              <Alert type="info" className="usa-width-one-whole" heading="">
                Most of the time it is simpler for you to use your government
                travel card for moving expenses rather than receiving a direct
                deposit advance. Not only do you save the effort of filling out
                the necessary forms to set up direct deposit, you eliminate any
                chance that the Government may unexpectedly garnish a part of
                your paycheck to recoup advance overages in the event that you
                move less weight than you originally estimated or take longer
                than 45 days to request payment upon arriving at your
                destination. <a onClick={this.closeInfo}>Close</a>
              </Alert>
            </div>
          )}
          {hasRequestedAdvance && (
            <div className="usa-width-one-whole top-buffered">
              <Alert type="info" heading="">
                We recommend that Service Families be cautious when requesting
                an advance on PPM expenses. Because your final incentive is
                affected by the amount of weight you actually move, if you
                request a full advance and then move less weight than
                anticipated, you may have to pay back some of your advance to
                the military. If you would like to use a more specific move
                calculator to estimate your anticipated shipment weight, you can
                do that{' '}
                <a href="https://www.move.mil/resources/weight-estimator">
                  here
                </a>.
              </Alert>
              <SwaggerField
                fieldName="requested_amount"
                swagger={schema.properties.advance}
                title={requestedTitle(maxAdvance)}
                required
              />
              <SwaggerField
                fieldName="method_of_receipt"
                swagger={schema.properties.advance}
                title={methodTitle}
                required
              />
            </div>
          )}
        </div>
      </div>
    );
  }
}

const WeightWizardForm = reduxifyWizardForm(
  requestAdvanceFormName,
  validateAdvanceForm,
);

export class PpmWeight extends Component {
  componentDidMount() {
    if (this.props.currentPpm) {
      this.updateIncentive();
    }
  }
  componentDidUpdate(prevProps, prevState) {
    if (
      !prevProps.hasLoadSuccess &&
      this.props.hasLoadSuccess &&
      this.props.currentPpm
    ) {
      this.updateIncentive();
    }
  }
  // this method is used to set the incentive on page load
  // it runs even if the incentive has been set before since data changes on previous pages could
  // affect it
  updateIncentive() {
    const {
      pendingPpmWeight,
      currentWeight,
      currentPpm,
      entitlement,
    } = this.props;
    const weight_estimate = get(this.props, 'currentPpm.weight_estimate');
    if (![pendingPpmWeight, weight_estimate].includes(currentWeight)) {
      this.onWeightSelecting(currentWeight);
      this.props.getPpmWeightEstimate(
        currentPpm.planned_move_date,
        currentPpm.pickup_postal_code,
        currentPpm.destination_postal_code,
        currentWeight,
      );
    }

    const currentInfo = getWeightInfo(currentPpm, entitlement);
    this.props.getPpmMaxWeightEstimate(
      currentPpm.planned_move_date,
      currentPpm.pickup_postal_code,
      currentPpm.destination_postal_code,
      currentInfo.max,
    );
  }
  handleSubmit = () => {
    const {
      pendingPpmWeight,
      incentive,
      createOrUpdatePpm,
      advanceFormData,
    } = this.props;
    const moveId = this.props.match.params.moveId;
    const ppmBody = {
      weight_estimate: pendingPpmWeight,
      estimated_incentive: incentive,
    };

    if (advanceFormData.values.has_requested_advance) {
      ppmBody.has_requested_advance = true;
      const requestedAmount = Math.round(
        parseFloat(advanceFormData.values.requested_amount) * 100,
      );
      ppmBody.advance = {
        requested_amount: requestedAmount,
        method_of_receipt: advanceFormData.values.method_of_receipt,
      };
    } else {
      ppmBody.has_requested_advance = false;
    }

    createOrUpdatePpm(moveId, ppmBody);
  };
  onWeightSelecting = value => {
    this.props.setPendingPpmWeight(value);
  };
  onWeightSelected = value => {
    const { currentPpm } = this.props;
    this.props.getPpmWeightEstimate(
      currentPpm.planned_move_date,
      currentPpm.pickup_postal_code,
      currentPpm.destination_postal_code,
      this.props.pendingPpmWeight,
    );
  };
  render() {
    const {
      currentPpm,
      incentive,
      pages,
      pageKey,
      hasSubmitSuccess,
      currentWeight,
      hasLoadSuccess,
      maxIncentive,
      hasEstimateInProgress,
      error,
      entitlement,
      hasEstimateError,
      schema,
      advanceFormData,
    } = this.props;
    let currentInfo = null;
    if (hasLoadSuccess) {
      currentInfo = getWeightInfo(currentPpm, entitlement);
    }

    const hasRequestedAdvance = get(
      advanceFormData,
      'values.has_requested_advance',
      false,
    );
    let advanceInitialValues = null;
    if (currentPpm) {
      let requestedAmount = get(currentPpm, 'advance.requested_amount');
      if (requestedAmount) {
        requestedAmount = parseFloat(requestedAmount) / 100;
      }
      advanceInitialValues = {
        has_requested_advance: currentPpm.has_requested_advance,
        requested_amount: requestedAmount,
        method_of_receipt: get(currentPpm, 'advance.method_of_receipt'),
      };
    }

    return (
      <WeightWizardForm
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        initialValues={advanceInitialValues}
        serverError={error}
        additionalValues={{ hasEstimateInProgress, maxIncentive }}
      >
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <h2>Customize Weight</h2>
        {!hasLoadSuccess && <LoadingPlaceholder />}
        {hasLoadSuccess && (
          <Fragment>
            <p>
              Use this slider to customize how much weight you think you’ll
              carry.
            </p>
            <div className="slider-container">
              <Slider
                min={currentInfo.min}
                max={currentInfo.max}
                value={currentWeight}
                onChange={this.onWeightSelecting}
                onChangeComplete={this.onWeightSelected}
                labels={{
                  [currentInfo.min]: currentInfo.min.toLocaleString(),
                  [currentInfo.max]: currentInfo.max.toLocaleString(),
                }}
              />
            </div>
            {hasEstimateError && (
              <Fragment>
                <div className="usa-width-one-whole error-message">
                  <Alert type="warning" heading="Could not retrieve estimate">
                    There was an issue retrieving an estimate for your
                    incentive. You still qualify but may need to talk with your
                    local PPPO.
                  </Alert>
                </div>
              </Fragment>
            )}
            <table className="numeric-info">
              <tbody>
                <tr>
                  <th>Your PPM Weight Estimate:</th>
                  <td className="current-weight"> {currentWeight}</td>
                </tr>
                <tr>
                  <th>Your PPM Incentive:</th>
                  <td className="incentive">{incentive}</td>
                </tr>
              </tbody>
            </table>

            <RequestAdvanceForm
              schema={schema}
              hasRequestedAdvance={hasRequestedAdvance}
              maxIncentive={maxIncentive}
              initialValues={advanceInitialValues}
            />

            <div className="info">
              <h3> How is my PPM Incentive calculated?</h3>
              <p>
                The government gives you 95% of what they would pay a mover when
                you move your own belongings, based on weight and distance. You
                pay taxes on this income. You can reduce the amount taxable
                incentive by saving receipts for approved expenses.
              </p>

              <p>
                This estimator just presents a range of possible incentives
                based on your anticipated shipment weight, anticipated moving
                date, and the specific route that you will be traveling. During
                your move, you will need to weigh the stuff you’re carrying, and
                submit weight tickets. We’ll let you know later how to weigh the
                stuff you carry.
              </p>
            </div>
          </Fragment>
        )}
      </WeightWizardForm>
    );
  }
}

PpmWeight.propTypes = {
  pendingPpmWeight: PropTypes.number,
  currentWeight: PropTypes.number,
  currentPpm: PropTypes.shape({
    id: PropTypes.string,
    size: PropTypes.string,
    weight: PropTypes.number,
    incentive: PropTypes.string,
  }),
  hasSubmitSuccess: PropTypes.bool.isRequired,
  hasLoadSuccess: PropTypes.bool.isRequired,
  setPendingPpmWeight: PropTypes.func.isRequired,
  entitlement: PropTypes.object,
};

function getMiddleWeight(ppm, entitlement) {
  const currentInfo = getWeightInfo(ppm, entitlement);
  return currentInfo.min + (currentInfo.max - currentInfo.min) / 2;
}
function mapStateToProps(state) {
  const entitlement = loadEntitlements(state);
  const defaultWeight = state.ppm.hasLoadSuccess
    ? getMiddleWeight(state.ppm.currentPpm, entitlement)
    : null;
  const currentWeight =
    state.ppm.pendingPpmWeight ||
    get(state, 'ppm.currentPpm.weight_estimate', defaultWeight);
  const props = {
    ...state.ppm,
    currentWeight,
    entitlement: loadEntitlements(state),
    schema: get(
      state,
      'swagger.spec.definitions.UpdatePersonallyProcuredMovePayload',
      {},
    ),
    advanceFormData: state.form[requestAdvanceFormName],
  };

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      setPendingPpmWeight,
      getPpmWeightEstimate,
      getPpmMaxWeightEstimate,
      createOrUpdatePpm,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmWeight);
