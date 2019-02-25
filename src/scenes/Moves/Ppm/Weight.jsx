import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get, cloneDeep, without, has } from 'lodash';
import PropTypes from 'prop-types';
import Slider from 'react-rangeslider'; //todo: pull from node_modules, override
import { Field } from 'redux-form';
import { getFormValues } from 'redux-form';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import Alert from 'shared/Alert';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { formatCents, formatCentsRange, formatNumber } from 'shared/formatters';
import { convertDollarsToCents } from 'shared/utils';
import { getPpmWeightEstimate, createOrUpdatePpm, getSelectedWeightInfo, getMaxAdvance } from './ducks';
import WizardHeader from '../WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import ppmBlack from 'shared/icon/ppm-black.svg';

import 'react-rangeslider/lib/index.css';
import './Weight.css';

const requestedTitle = maxAdvance => {
  return (
    <Fragment>
      <div className="ppmquestion">How much advance do you want?</div>
      <div className="ppmmuted">Up to ${formatCents(maxAdvance)} (60% of your PPM incentive)</div>
    </Fragment>
  );
};

const methodTitle = (
  <Fragment>
    <div className="ppmquestion">How do you want to get your advance?</div>
    <div className="ppmmuted">
      To direct deposit to another account you'll need to fill out a new account form, included in your advance
      paperwork, and take it to the accounting office.
    </div>
  </Fragment>
);

const validateAdvanceForm = (values, form) => {
  if (values.hasEstimateInProgress) {
    return { has_requested_advance: 'Estimate in progress.' };
  }

  const maxAdvance = values.maxAdvance;

  if (parseFloat(values.requested_amount) > maxAdvance / 100) {
    return {
      requested_amount: `Must be less than $${formatCents(maxAdvance)}`,
    };
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
    const { hasRequestedAdvance, maxAdvance, ppmAdvanceSchema } = this.props;
    return (
      <div className="whole_box">
        <div>
          <div className="usa-width-one-whole">
            <div className="usa-width-two-thirds">
              <div className="ppmquestion">
                Would you like an advance of up to 60% of your PPM incentive? ($
                {formatCents(maxAdvance)})
              </div>
              <div className="ppmmuted">
                We recommend paying for expenses with your government travel card, rather than getting an advance.{' '}
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
                Most of the time it is simpler for you to use your government travel card for moving expenses rather
                than receiving a direct deposit advance. Not only do you save the effort of filling out the necessary
                forms to set up direct deposit, you eliminate any chance that the Government may unexpectedly garnish a
                part of your paycheck to recoup advance overages in the event that you move less weight than you
                originally estimated or take longer than 45 days to request payment upon arriving at your destination.{' '}
                <a onClick={this.closeInfo}>Close</a>
              </Alert>
            </div>
          )}
          {hasRequestedAdvance && (
            <div className="usa-width-one-whole top-buffered">
              <Alert type="info" heading="">
                We recommend that Service Families be cautious when requesting an advance on PPM expenses. Because your
                final incentive is affected by the amount of weight you actually move, if you request a full advance and
                then move less weight than anticipated, you may have to pay back some of your advance to the military.
                If you would like to use a more specific move calculator to estimate your anticipated shipment weight,
                you can do that <a href="https://www.move.mil/resources/weight-estimator">here</a>.
              </Alert>
              <SwaggerField
                fieldName="requested_amount"
                swagger={ppmAdvanceSchema}
                title={requestedTitle(maxAdvance)}
                required
              />
              <SwaggerField fieldName="method_of_receipt" swagger={ppmAdvanceSchema} title={methodTitle} required />
            </div>
          )}
        </div>
      </div>
    );
  }
}

const WeightWizardForm = reduxifyWizardForm(requestAdvanceFormName, validateAdvanceForm);

export class PpmWeight extends Component {
  constructor(props) {
    super(props);

    this.state = {
      pendingPpmWeight: null,
    };
  }

  componentDidMount() {
    const { currentPpm } = this.props;
    if (currentPpm) {
      this.setState({
        pendingPpmWeight: currentPpm.weight_estimate,
      });
      this.updateIncentive();
    }
  }
  componentDidUpdate(prevProps, prevState) {
    const { currentPpm, hasLoadSuccess } = this.props;
    if (!prevProps.hasLoadSuccess && hasLoadSuccess && currentPpm) {
      this.setState({
        pendingPpmWeight: currentPpm.weight_estimate,
      });
      this.updateIncentive();
    }
  }
  // this method is used to set the incentive on page load
  // it runs even if the incentive has been set before since data changes on previous pages could
  // affect it
  updateIncentive() {
    const { currentWeight, currentPpm } = this.props;
    const weight_estimate = get(this.props, 'currentPpm.weight_estimate');
    if (![this.state.pendingPpmWeight, weight_estimate].includes(currentWeight)) {
      this.onWeightSelecting(currentWeight);
      this.props.getPpmWeightEstimate(
        currentPpm.original_move_date,
        currentPpm.pickup_postal_code,
        currentPpm.destination_postal_code,
        currentWeight,
      );
    }
  }
  handleSubmit = () => {
    const { createOrUpdatePpm, advanceFormValues } = this.props;
    const moveId = this.props.match.params.moveId;
    const ppmBody = {
      weight_estimate: this.state.pendingPpmWeight,
    };
    if (advanceFormValues.has_requested_advance) {
      ppmBody.has_requested_advance = true;
      const requestedAmount = convertDollarsToCents(advanceFormValues.requested_amount);
      ppmBody.advance = {
        requested_amount: requestedAmount,
        method_of_receipt: advanceFormValues.method_of_receipt,
      };
    } else {
      ppmBody.has_requested_advance = false;
    }
    return createOrUpdatePpm(moveId, ppmBody);
  };
  onWeightSelecting = value => {
    this.setState({
      pendingPpmWeight: value,
    });
  };
  onWeightSelected = value => {
    const { currentPpm } = this.props;
    this.props.getPpmWeightEstimate(
      currentPpm.original_move_date,
      currentPpm.pickup_postal_code,
      currentPpm.destination_postal_code,
      this.state.pendingPpmWeight,
    );
  };
  render() {
    const {
      currentPpm,
      incentive_estimate_min,
      incentive_estimate_max,
      maxAdvance,
      pages,
      pageKey,
      hasLoadSuccess,
      hasEstimateInProgress,
      error,
      hasEstimateError,
      ppmAdvanceSchema,
      advanceFormValues,
      selectedWeightInfo,
      isHHGPPMComboMove,
    } = this.props;
    const hasRequestedAdvance = get(advanceFormValues, 'has_requested_advance', false);
    let advanceInitialValues = null;
    if (currentPpm) {
      let requestedAmount = get(currentPpm, 'advance.requested_amount');
      if (requestedAmount) {
        requestedAmount = formatCents(requestedAmount);
      }
      let methodOfReceipt = get(currentPpm, 'advance.method_of_receipt');
      // GTCC is an invalid method of receipt in PPM advances, so default to direct deposit
      if (methodOfReceipt === 'GTCC') {
        methodOfReceipt = 'OTHER_DD';
      }
      advanceInitialValues = {
        has_requested_advance: currentPpm.has_requested_advance,
        requested_amount: requestedAmount,
        method_of_receipt: methodOfReceipt,
      };
    }

    return (
      <div>
        {isHHGPPMComboMove && (
          <WizardHeader
            icon={ppmBlack}
            title="Move Setup"
            right={
              <ProgressTimeline>
                <ProgressTimelineStep name="Move Setup" current />
                <ProgressTimelineStep name="Review" />
              </ProgressTimeline>
            }
          />
        )}
        <WeightWizardForm
          handleSubmit={this.handleSubmit}
          pageList={pages}
          pageKey={pageKey}
          initialValues={advanceInitialValues}
          serverError={error}
          additionalValues={{
            hasEstimateInProgress,
            incentive_estimate_max,
            maxAdvance,
          }}
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
              <p>Use this slider to customize how much weight you think you’ll carry.</p>
              <div className="slider-container">
                <Slider
                  min={selectedWeightInfo.min}
                  max={selectedWeightInfo.max}
                  value={this.state.pendingPpmWeight}
                  onChange={this.onWeightSelecting}
                  onChangeComplete={this.onWeightSelected}
                  labels={{
                    [selectedWeightInfo.min]: `${selectedWeightInfo.min} lbs`,
                    [selectedWeightInfo.max]: `${selectedWeightInfo.max} lbs`,
                  }}
                />
              </div>
              {hasEstimateError && (
                <Fragment>
                  <div className="usa-width-one-whole error-message">
                    <Alert type="warning" heading="Could not retrieve estimate">
                      There was an issue retrieving an estimate for your incentive. You still qualify, but need to talk
                      with your local transportation office which you can look up on <a href="move.mil">move.mil</a>
                    </Alert>
                  </div>
                </Fragment>
              )}
              <table className="numeric-info">
                <tbody>
                  {!isHHGPPMComboMove && (
                    <tr>
                      <th>Your PPM Weight Estimate:</th>
                      <td className="current-weight"> {formatNumber(this.state.pendingPpmWeight)} lbs.</td>
                    </tr>
                  )}
                  <tr>
                    <th>Your PPM Incentive:</th>
                    <td className="incentive">{formatCentsRange(incentive_estimate_min, incentive_estimate_max)}</td>
                  </tr>
                </tbody>
              </table>

              {!isHHGPPMComboMove && (
                <RequestAdvanceForm
                  ppmAdvanceSchema={ppmAdvanceSchema}
                  hasRequestedAdvance={hasRequestedAdvance}
                  maxAdvance={maxAdvance}
                  initialValues={advanceInitialValues}
                />
              )}

              <div className="info">
                <h3> How is my PPM Incentive calculated?</h3>
                <p>
                  The government gives you 95% of what they would pay a mover when you move your own belongings, based
                  on weight and distance. You pay taxes on this income. You can reduce the amount taxable incentive by
                  saving receipts for approved expenses.
                </p>

                <p>
                  This estimator just presents a range of possible incentives based on your anticipated shipment weight,
                  anticipated moving date, and the specific route that you will be traveling. During your move, you will
                  need to weigh the stuff you’re carrying, and submit weight tickets. We’ll let you know later how to
                  weigh the stuff you carry.
                </p>
              </div>
            </Fragment>
          )}
        </WeightWizardForm>
      </div>
    );
  }
}

PpmWeight.propTypes = {
  currentWeight: PropTypes.number,
  currentPpm: PropTypes.shape({
    id: PropTypes.string,
    size: PropTypes.string,
    weight: PropTypes.number,
    incentive: PropTypes.string,
  }),
  hasLoadSuccess: PropTypes.bool.isRequired,
};
function mapStateToProps(state) {
  const schema = get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {});
  // In scheduling, PPM advances cannot go to GTCC so we filter out that method of payment.
  let ppmAdvanceSchema = {};
  if (has(schema, 'properties')) {
    ppmAdvanceSchema = cloneDeep(schema.properties.advance);
    ppmAdvanceSchema.properties.method_of_receipt.enum = without(
      ppmAdvanceSchema.properties.method_of_receipt.enum,
      'GTCC',
    );
  }

  const props = {
    ...state.ppm,
    maxAdvance: getMaxAdvance(state),
    selectedWeightInfo: getSelectedWeightInfo(state),
    currentWeight: get(state, 'ppm.currentPpm.weight_estimate'),
    schema: schema,
    ppmAdvanceSchema: ppmAdvanceSchema,
    advanceFormValues: getFormValues(requestAdvanceFormName)(state),
    isHHGPPMComboMove: get(state, 'moves.currentMove.selected_move_type') === 'HHG_PPM',
  };

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      getPpmWeightEstimate,
      createOrUpdatePpm,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmWeight);
