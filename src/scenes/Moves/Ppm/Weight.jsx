import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import Slider from 'react-rangeslider'; //todo: pull from node_modules, override

import { loadEntitlements } from 'scenes/Orders/ducks';
import WizardPage from 'shared/WizardPage';
import Alert from 'shared/Alert';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  setPendingPpmWeight,
  loadPpm,
  getPpmWeightEstimate,
  createOrUpdatePpm,
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
export class PpmWeight extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Weight Selection';
    if (this.props.currentPpm) {
      this.updateIncentive();
    }
  }
  componentDidUpdate(prevProps, prevState) {
    if (
      !prevProps.loggedInUser.hasSucceeded &&
      this.props.loggedInUser.hasSucceeded
    ) {
      this.props.loadPpm(this.props.match.params.moveId);
    }

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
    const { pendingPpmWeight, currentWeight, currentPpm } = this.props;
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
  }
  handleSubmit = () => {
    const { pendingPpmWeight, incentive, createOrUpdatePpm } = this.props;
    const moveId = this.props.match.params.moveId;
    createOrUpdatePpm(moveId, {
      weight_estimate: pendingPpmWeight,
      estimated_incentive: incentive,
    });
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
      pendingPpmWeight,
      currentPpm,
      incentive,
      pages,
      pageKey,
      hasSubmitSuccess,
      currentWeight,
      hasLoadSuccess,
      hasEstimateInProgress,
      error,
      entitlement,
    } = this.props;
    let currentInfo = null;
    if (hasLoadSuccess) {
      currentInfo = getWeightInfo(currentPpm, entitlement);
    }
    const isValid = incentive && !hasEstimateInProgress;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={Boolean(pendingPpmWeight)}
        hasSucceeded={hasSubmitSuccess}
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
          <React.Fragment>
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
          </React.Fragment>
        )}
      </WizardPage>
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
    loggedInUser: state.loggedInUser,
    currentWeight,
    entitlement: loadEntitlements(state),
  };

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { setPendingPpmWeight, loadPpm, getPpmWeightEstimate, createOrUpdatePpm },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmWeight);
