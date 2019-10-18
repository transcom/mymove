import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import Slider from 'react-rangeslider'; //todo: pull from node_modules, override
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import Alert from 'shared/Alert';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCentsRange, formatNumber } from 'shared/formatters';
import { getPpmWeightEstimate, createOrUpdatePpm, getSelectedWeightInfo } from './ducks';
import { updatePPMEstimate } from 'shared/Entities/modules/ppms';
import 'react-rangeslider/lib/index.css';
import './Weight.css';

const WeightWizardForm = reduxifyWizardForm('weight-wizard-form');

export class PpmWeight extends Component {
  constructor(props) {
    super(props);

    this.state = {
      pendingPpmWeight: null,
    };

    this.onWeightSelected = this.onWeightSelected.bind(this);
  }

  getWeightClassMedian() {
    const { selectedWeightInfo } = this.props;
    return selectedWeightInfo.min + (selectedWeightInfo.max - selectedWeightInfo.min) / 2;
  }

  componentDidMount() {
    const { currentPpm } = this.props;
    if (currentPpm) {
      this.setState(
        {
          pendingPpmWeight:
            currentPpm.weight_estimate && currentPpm.weight_estimate !== 0
              ? currentPpm.weight_estimate
              : this.getWeightClassMedian(),
        },
        this.updateIncentive,
      );
    }
  }
  componentDidUpdate(prevProps, prevState) {
    const { currentPpm, hasLoadSuccess } = this.props;
    if (!prevProps.hasLoadSuccess && hasLoadSuccess && currentPpm) {
      this.setState(
        {
          pendingPpmWeight:
            currentPpm.weight_estimate && currentPpm.weight_estimate !== 0
              ? currentPpm.weight_estimate
              : this.getWeightClassMedian(),
        },
        this.updateIncentive,
      );
    }
  }
  // this method is used to set the incentive on page load
  // it runs even if the incentive has been set before since data changes on previous pages could
  // affect it
  updateIncentive() {
    const { currentWeight, currentPpm, originDutyStationZip } = this.props;
    const newWeight = currentWeight && currentWeight !== 0 ? currentWeight : this.state.pendingPpmWeight;
    this.onWeightSelecting(newWeight);
    this.props.getPpmWeightEstimate(
      currentPpm.original_move_date,
      currentPpm.pickup_postal_code,
      originDutyStationZip,
      currentPpm.destination_postal_code,
      newWeight,
    );
  }

  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const ppmBody = {
      weight_estimate: this.state.pendingPpmWeight,
      has_requested_advance: false,
    };
    return this.props
      .createOrUpdatePpm(moveId, ppmBody)
      .then(({ payload }) => this.props.updatePPMEstimate(moveId, payload.id).catch(err => err));
    // catch block returns error so that the wizard can continue on with its flow
  };

  onWeightSelecting = value => {
    this.setState({
      pendingPpmWeight: value,
    });
  };

  onWeightSelected() {
    const { currentPpm, originDutyStationZip } = this.props;
    this.props.getPpmWeightEstimate(
      currentPpm.original_move_date,
      currentPpm.pickup_postal_code,
      originDutyStationZip,
      currentPpm.destination_postal_code,
      this.state.pendingPpmWeight,
    );
  }
  render() {
    const {
      incentive_estimate_min,
      incentive_estimate_max,
      pages,
      pageKey,
      hasLoadSuccess,
      hasEstimateInProgress,
      error,
      hasEstimateError,
      selectedWeightInfo,
    } = this.props;
    return (
      <div className="grid-container usa-prose site-prose">
        <WeightWizardForm
          handleSubmit={this.handleSubmit}
          pageList={pages}
          pageKey={pageKey}
          serverError={error}
          additionalValues={{
            hasEstimateInProgress,
            incentive_estimate_max,
          }}
        >
          {error && (
            <div className="grid-row">
              <div className="grid-col-12">
                <Alert type="error" heading="An error occurred">
                  {error.message}
                </Alert>
              </div>
            </div>
          )}
          <div className="grid-row">
            <div className="grid-col-12">
              <h1>Customize Weight</h1>
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
                      <div className="error-message">
                        <Alert type="warning" heading="Could not retrieve estimate">
                          There was an issue retrieving an estimate for your incentive. You still qualify, but need to
                          talk with your local transportation office which you can look up on{' '}
                          <a href="move.mil">move.mil</a>
                        </Alert>
                      </div>
                    </Fragment>
                  )}
                  <table className="numeric-info">
                    <tbody>
                      <tr>
                        <th>Your PPM Weight Estimate:</th>
                        <td className="current-weight"> {formatNumber(this.state.pendingPpmWeight)} lbs.</td>
                      </tr>
                      <tr>
                        <th>Your PPM Incentive:</th>
                        {hasEstimateError ? (
                          <td className="incentive">
                            Not ready yet{' '}
                            <IconWithTooltip toolTipText="We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive." />
                          </td>
                        ) : (
                          <td className="incentive">
                            {formatCentsRange(incentive_estimate_min, incentive_estimate_max)}
                          </td>
                        )}
                      </tr>
                    </tbody>
                  </table>

                  <div className="info">
                    <h3> How is my PPM Incentive calculated?</h3>
                    <p>
                      The government gives you 95% of what they would pay a mover when you move your own belongings,
                      based on weight and distance. You pay taxes on this income. You can reduce the amount taxable
                      incentive by saving receipts for approved expenses.
                    </p>

                    <p>
                      This estimator just presents a range of possible incentives based on your anticipated shipment
                      weight, anticipated moving date, and the specific route that you will be traveling. During your
                      move, you will need to weigh the stuff you’re carrying, and submit weight tickets. We’ll let you
                      know later how to weigh the stuff you carry.
                    </p>
                  </div>
                </Fragment>
              )}
            </div>
          </div>
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
  const originDutyStationZip = state.serviceMember.currentServiceMember.current_station.address.postal_code;

  const props = {
    ...state.ppm,
    selectedWeightInfo: getSelectedWeightInfo(state),
    currentWeight: get(state, 'ppm.currentPpm.weight_estimate'),
    schema: schema,
    originDutyStationZip,
  };

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      getPpmWeightEstimate,
      createOrUpdatePpm,
      updatePPMEstimate,
    },
    dispatch,
  );
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(PpmWeight);
