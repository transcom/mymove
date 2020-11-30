import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { debounce, get } from 'lodash';
import SaveCancelButtons from './SaveCancelButtons';
import { push } from 'connected-react-router';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { formatCents } from 'shared/formatters';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import {
  loadPPMs,
  updatePPM,
  selectActivePPMForMove,
  selectPPMEstimateRange,
  updatePPMEstimate,
  getPpmWeightEstimate,
} from 'shared/Entities/modules/ppms';
import { fetchLatestOrders, selectActiveOrLatestOrders } from 'shared/Entities/modules/orders';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { formatCentsRange } from 'shared/formatters';
import { editBegin, editSuccessful, entitlementChangeBegin, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

import EntitlementBar from 'scenes/EntitlementBar';
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
    description: 'A moving truck',
  },
  { weight: 10000, incentive: '$5,900 - 6,800' },
];

const validateWeight = (value, formValues, props, fieldName) => {
  if (value && props.entitlement && value > props.entitlement.sum) {
    return 'Cannot be more than your full entitlement';
  }
};

let EditWeightForm = (props) => {
  const {
    schema,
    handleSubmit,
    submitting,
    valid,
    entitlement,
    dirty,
    incentiveEstimateMin,
    incentiveEstimateMax,
    onWeightChange,
    initialValues,
  } = props;
  // Error class if below advance amount, otherwise warn class if incentive has changed
  let incentiveClass = '';
  let fieldClass = dirty ? 'warn' : '';
  let advanceError = false;
  const advanceAmt = get(initialValues, 'advance.requested_amount', 0);
  if (incentiveEstimateMax && advanceAmt && incentiveEstimateMax < formatCents(advanceAmt)) {
    advanceError = true;
    incentiveClass = 'error';
    fieldClass = 'error';
  } else if (get(initialValues, 'incentive_estimate_min') !== incentiveEstimateMin) {
    // Min and max are linked, so we only need to check one
    incentiveClass = 'warn';
  }

  const fullFieldClass = `weight-estimate-input ${fieldClass}`;
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <img src={profileImage} alt="" />
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 0,
                marginTop: 20,
              }}
            >
              Profile
            </h1>
            <hr />
            <h3>Edit PPM Weight:</h3>
            <p>Changes could impact your move, including the estimated PPM incentive.</p>
            <EntitlementBar entitlement={entitlement} />
            <div className="edit-weight-container">
              <div className="usa-width-one-half">
                <h4>Move estimate</h4>
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
                  {!advanceError && initialValues && initialValues.incentive_estimate_min && dirty && (
                    <div className="usa-alert usa-alert--warning">
                      <div className="usa-alert__body">
                        <p className="usa-alert__text">This update will change your incentive.</p>
                      </div>
                    </div>
                  )}
                  {advanceError && (
                    <p className="advance-error">Weight is too low and will require paying back the advance.</p>
                  )}
                </div>

                <div className="display-value">
                  <p>Estimated Incentive</p>
                  <p className={incentiveClass}>
                    <strong>
                      {formatCentsRange(incentiveEstimateMin, incentiveEstimateMax) || 'Unable to Calculate'}
                    </strong>
                  </p>
                  {initialValues &&
                    initialValues.incentive_estimate_min &&
                    initialValues.incentive_estimate_min !== incentiveEstimateMin && (
                      <p className="subtext">
                        Originally{' '}
                        {formatCentsRange(initialValues.incentive_estimate_min, initialValues.incentive_estimate_max)}
                      </p>
                    )}
                </div>

                {get(initialValues, 'has_requested_advance') && (
                  <div className="display-value">
                    <p>Advance</p>
                    <p>
                      <strong>${formatCents(initialValues.advance.requested_amount)}</strong>
                    </p>
                  </div>
                )}
              </div>

              <div className="usa-width-one-half">
                <h4>Examples</h4>
                <table className="examples-table">
                  <thead>
                    <tr>
                      <th>Weight</th>
                      <th>Incentive</th>
                      <th />
                    </tr>
                  </thead>
                  <tbody>
                    {examples.map((ex) => (
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
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};
EditWeightForm = reduxForm({
  form: editWeightFormName,
  enableReinitialize: true,
})(EditWeightForm);

class EditWeight extends Component {
  constructor(props) {
    super(props);
    this.state = { currentPPM: {} };
  }

  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
    this.props.loadPPMs(this.props.match.params.moveId);
    this.props.fetchLatestOrders(this.props.serviceMemberId);
    scrollToTop();
  }

  debouncedGetPpmWeightEstimate = debounce(this.props.getPpmWeightEstimate, weightEstimateDebounce);

  onWeightChange = (e, newValue) => {
    const { currentPPM, entitlement, originDutyStationZip, orders } = this.props;
    if (newValue > 0 && newValue <= entitlement.sum) {
      this.debouncedGetPpmWeightEstimate(
        currentPPM.original_move_date,
        currentPPM.pickup_postal_code,
        originDutyStationZip,
        orders.id,
        newValue,
      );
    } else {
      this.debouncedGetPpmWeightEstimate.cancel();
    }
  };

  updatePpm = (values, dispatch, props) => {
    const moveId = this.props.match.params.moveId;
    return this.props
      .updatePPM(moveId, this.props.currentPPM.id, {
        weight_estimate: values.weight_estimate,
      })
      .then(({ response }) => {
        this.props
          .updatePPMEstimate(moveId, response.body.id)
          .then(() => {
            if (!this.props.hasSubmitError) {
              this.props.editSuccessful();
              this.props.history.goBack();
              this.props.checkEntitlement(moveId);
            } else {
              scrollToTop();
            }
          })
          .catch(() => {
            if (!this.props.hasSubmitError) {
              this.props.editSuccessful();
              this.props.history.goBack();
              this.props.checkEntitlement(moveId);
            } else {
              scrollToTop();
            }
          });
      });
  };

  chooseEstimateErrorText(hasEstimateError, rateEngineError) {
    if (rateEngineError) {
      return (
        <div className="grid-row">
          <div className="grid-col-12 error-message">
            <Alert type="warning" heading="Could not retrieve estimate">
              MilMove does not presently support short-haul PPM moves. Please contact your PPPO.
            </Alert>
          </div>
        </div>
      );
    }

    if (hasEstimateError) {
      return (
        <div className="grid-row">
          <div className="grid-col-12 error-message">
            <Alert type="warning" heading="Could not retrieve estimate">
              There was an issue retrieving an estimate for your incentive. You still qualify but may need to talk with
              your local PPPO.
            </Alert>
          </div>
        </div>
      );
    }
  }

  render() {
    const {
      error,
      schema,
      entitlement,
      hasEstimateError,
      rateEngineError,
      currentPPM,
      incentiveEstimateMin,
      incentiveEstimateMax,
    } = this.props;

    return (
      <div className="grid-container usa-prose">
        {error && (
          <div className="grid-row">
            <div className="grid-col-12 error-message">
              <Alert type="error" heading="An error occurred">
                {error.message}
              </Alert>
            </div>
          </div>
        )}

        <div className="grid-container usa-prose">
          <div className="grid-row">
            <div className="grid-col-12">{this.chooseEstimateErrorText(hasEstimateError, rateEngineError)}</div>
          </div>
        </div>

        <div className="grid-row">
          <div className="grid-col-12">
            <EditWeightForm
              initialValues={currentPPM}
              incentiveEstimateMin={incentiveEstimateMin}
              incentiveEstimateMax={incentiveEstimateMax}
              onSubmit={this.updatePpm}
              onWeightChange={this.onWeightChange}
              entitlement={entitlement}
              schema={schema}
            />
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const moveID = state.moves.currentMove.id;
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;

  return {
    serviceMemberId,
    currentPPM: selectActivePPMForMove(state, moveID),
    incentiveEstimateMin: selectPPMEstimateRange(state).range_min,
    incentiveEstimateMax: selectPPMEstimateRange(state).range_max,
    entitlement: loadEntitlementsFromState(state),
    schema: get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {}),
    originDutyStationZip: serviceMember?.current_station?.address?.postal_code,
    orders: selectActiveOrLatestOrders(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
      loadPPMs,
      fetchLatestOrders,
      updatePPM,
      getPpmWeightEstimate,
      editBegin,
      editSuccessful,
      entitlementChangeBegin,
      checkEntitlement,
      updatePPMEstimate,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(EditWeight);
