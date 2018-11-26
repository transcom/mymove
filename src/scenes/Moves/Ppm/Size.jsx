import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { get } from 'lodash';
import { setPendingPpmSize, getRawWeightInfo, getActualRemainingWeight, getEstimatedRemainingWeight } from './ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import EntitlementBar from 'scenes/EntitlementBar';
import BigButton from 'shared/BigButton';
import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './Size.css';

function rangeFormatter(range) {
  return `${range.min.toLocaleString()} - ${range.max.toLocaleString()}`;
}

class BigButtonGroup extends Component {
  render() {
    const { onClick, weightInfo } = this.props;
    var createButton = (value, firstLine, secondLine, icon, altTag) => {
      var onButtonClick = () => {
        onClick(value);
      };
      let selected = this.props.selectedOption === value;
      let selectedClass = selected ? 'selected' : '';
      let radioClass = `radio ${selectedClass}`;
      return (
        <BigButton value={value} selected={selected} onClick={onButtonClick}>
          <div className="button-container">
            <div className="radio-container">
              <div className={radioClass}>{selected && '\u2714'}</div>
            </div>
            <div className="contents">
              <div className="text">
                <p>{firstLine}</p>
                <p>{secondLine}</p>
              </div>
              <img className="icon" src={icon} alt={altTag} />
            </div>
          </div>
        </BigButton>
      );
    };
    const smallSize = 'S';
    const medSize = 'M';
    const largeSize = 'L';
    /* eslint-disable security/detect-object-injection */
    var small = createButton(
      smallSize,
      'A few items in your car?',
      `(approx ${rangeFormatter(weightInfo[smallSize])} lbs)`,
      carGray,
      'car-gray',
    );
    var medium = createButton(
      medSize,
      'A trailer full of household goods?',
      `(approx ${rangeFormatter(weightInfo[medSize])} lbs)`,
      trailerGray,
      'trailer-gray',
    );
    var large = createButton(
      largeSize,
      'A moving truck that you rent yourself?',
      `(approx ${rangeFormatter(weightInfo[largeSize])} lbs)`,
      truckGray,
      'truck-gray',
    );
    /* eslint-enable security/detect-object-injection */
    return (
      <div>
        <div className="usa-width-one-third">{small}</div>
        <div className="usa-width-one-third">{medium}</div>
        <div className="usa-width-one-third">{large}</div>
      </div>
    );
  }
}
BigButtonGroup.propTypes = {
  selectedOption: PropTypes.string,
  onClick: PropTypes.func.isRequired,
};

export class PpmSize extends Component {
  onMoveTypeSelected = value => {
    this.props.setPendingPpmSize(value);
  };

  render() {
    const {
      pendingPpmSize,
      currentPpm,
      entitlement,
      actualRemainingWeight,
      estimatedRemainingWeight,
      weightInfo,
      isHHGPPMComboMove,
    } = this.props;
    const selectedOption = pendingPpmSize || (currentPpm && currentPpm.size);

    const weightRemainingEntitlementMsg = () => {
      if (isHHGPPMComboMove) {
        if (actualRemainingWeight < entitlement.sum) {
          return `${entitlement.sum - actualRemainingWeight} lbs. entitlement remaining
            (${entitlement.sum} lbs. - ${actualRemainingWeight} lbs. estimated HHG weight).`;
        } else {
          return `You have no entitlement remaining (Estimated ${entitlement.sum} lbs. estimated HHG weight).`;
        }

        if (estimatedRemainingWeight < entitlement.sum) {
          return `Estimated ${entitlement.sum - estimatedRemainingWeight} lbs. entitlement remaining
            (${entitlement.sum} lbs. - ${estimatedRemainingWeight} lbs. estimated HHG weight).`;
        } else {
          return `We estimate you have no entitlement remaining
            (Estimated ${entitlement.sum} lbs. estimated HHG weight).`;
        }
      }
    };

    return (
      <div className="usa-grid-full ppm-size-content">
        {weightInfo && (
          <Fragment>
            <h3>How much will you move?</h3>

            {<EntitlementBar hhgPPMEntitlementMessage={weightRemainingEntitlementMsg()} entitlement={entitlement} />}

            <BigButtonGroup selectedOption={selectedOption} onClick={this.onMoveTypeSelected} weightInfo={weightInfo} />
          </Fragment>
        )}
      </div>
    );
  }
}

PpmSize.propTypes = {
  pendingPpmSize: PropTypes.string,
  weightInfo: PropTypes.object,
  currentPpm: PropTypes.shape({ id: PropTypes.string, size: PropTypes.string }),
  setPendingPpmSize: PropTypes.func.isRequired,
  entitlement: PropTypes.object,
  estimatedRemainingWeight: PropTypes.number,
  actualRemainingWeight: PropTypes.number,
};

function mapStateToProps(state) {
  return {
    ...state.ppm,
    weightInfo: getRawWeightInfo(state),
    entitlement: loadEntitlementsFromState(state),
    isHHGPPMComboMove: get(state, 'moves.currentMove.selected_move_type') === 'HHG_PPM',
    estimatedRemainingWeight: getEstimatedRemainingWeight(state),
    actualRemainingWeight: getActualRemainingWeight(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ setPendingPpmSize }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmSize);
