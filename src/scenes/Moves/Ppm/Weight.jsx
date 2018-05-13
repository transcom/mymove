import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import Slider from 'react-rangeslider'; //todo: pull from node_modules, override

import WizardPage from 'shared/WizardPage';
import Alert from 'shared/Alert';
import {
  setPendingPpmWeight,
  loadPpm,
  getPpmEstimate,
  createOrUpdatePpm,
} from './ducks';

import 'react-rangeslider/lib/index.css';
import './Weight.css';
import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';

function getWeightInfo(ppm) {
  const size = ppm ? ppm.size : 'L';
  switch (size) {
    case 'S':
      return {
        icon: carGray,
        altTag: 'car-gray',
        min: 100,
        max: 800,
        vehicle: 'your car',
      };
    case 'M':
      return {
        icon: trailerGray,
        altTag: 'trailer-gray',
        min: 400,
        max: 1200,
        vehicle: 'a trailer',
      };
    default:
      return {
        icon: truckGray,
        altTag: 'truck-gray',
        defaultWeight: 800,
        min: 1000,
        max: 5000,
        vehicle: 'a truck',
      };
  }
}
export class PpmWeight extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Weight Selection';
    this.props.loadPpm(this.props.match.params.moveId);
  }
  handleSubmit = () => {
    const { pendingPpmWeight, incentive, createOrUpdatePpm } = this.props;
    //todo: we should make sure this move matches the redux state
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
    this.props.getPpmEstimate(
      currentPpm.planned_move_date,
      currentPpm.pickup_zip,
      currentPpm.destination_zip,
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
      error,
    } = this.props;
    const currentInfo = getWeightInfo(currentPpm);
    const setOrPendingWeight =
      pendingPpmWeight || (currentPpm && currentPpm.weight);
    const currentWeight =
      setOrPendingWeight ||
      currentInfo.min + (currentInfo.max - currentInfo.min) / 2;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={Boolean(setOrPendingWeight)}
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
        <h2>
          <img
            className="icon"
            src={currentInfo.icon}
            alt={currentInfo.altTag}
          />{' '}
          You selected {currentInfo.min} - {currentInfo.max} pounds in{' '}
          {currentInfo.vehicle}.
        </h2>
        <p>
          Use this slider to customize how much weight you think you’ll carry.
        </p>
        <div className="slider-container">
          <Slider
            min={currentInfo.min}
            max={currentInfo.max}
            value={currentWeight}
            onChange={this.onWeightSelecting}
            onChangeComplete={this.onWeightSelected}
            labels={{
              [currentInfo.min]: currentInfo.min,
              [currentInfo.max]: currentInfo.max,
              //[currentWeight]: currentWeight,
            }}
          />
        </div>
        <h4>
          {' '}
          Your PPM Incentive: <span className="incentive">{incentive}</span>
        </h4>
        <div className="info">
          <h3> How is my PPM Incentive calculated?</h3>
          <p>
            The government gives you 95% of what they would pay a mover when you
            move your own belongings, based on weight and distance. You pay
            taxes on this income. You can reduce the amount taxable incentive by
            saving receipts for approved expenses.
          </p>

          <p>
            This estimator just presents a range of possible incentives based on
            your anticipated shipment weight, anticipated moving date, and the
            specific route that you will be traveling. During your move, you
            will need to weigh the stuff you’re carrying, and submit weight
            tickets. We’ll let you know later how to weigh the stuff you carry.
          </p>
        </div>
      </WizardPage>
    );
  }
}

PpmWeight.propTypes = {
  pendingPpmWeight: PropTypes.number,
  currentPpm: PropTypes.shape({
    id: PropTypes.string,
    size: PropTypes.string,
    weight: PropTypes.number,
    incentive: PropTypes.string,
  }),
  hasSubmitSuccess: PropTypes.bool.isRequired,
  moveDistance: PropTypes.number,
  setPendingPpmWeight: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return { ...state.ppm, moveDistance: 666 };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    //todo: will want to load move to get distance
    { setPendingPpmWeight, loadPpm, getPpmEstimate, createOrUpdatePpm },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmWeight);
