import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { setPendingPpmSize } from './ducks';
import { loadEntitlements } from 'scenes/Orders/ducks';
import EntitlementBar from 'scenes/EntitlementBar';
import BigButton from 'shared/BigButton';
import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './Size.css';

class BigButtonGroup extends Component {
  render() {
    const { onClick, maxWeight } = this.props;
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
    var small = createButton(
      'S',
      'A few items in your car?',
      '(approx 50 - 1,000 lbs)',
      carGray,
      'car-gray',
    );
    var medium = createButton(
      'M',
      'A trailer full of household goods?',
      '(approx 500 - 2,500 lbs)',
      trailerGray,
      'trailer-gray',
    );
    var large = createButton(
      'L',
      'A moving truck that you rent yourself?',
      `(approx 1,500 - ${maxWeight.toLocaleString()} lbs)`,
      truckGray,
      'truck-gray',
    );

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
    const { pendingPpmSize, currentPpm, entitlement } = this.props;
    const selectedOption = pendingPpmSize || (currentPpm && currentPpm.size);
    return (
      <div className="usa-grid-full ppm-size-content">
        {entitlement && (
          <Fragment>
            <h3>How much will you move?</h3>

            <EntitlementBar entitlement={entitlement} />

            <BigButtonGroup
              selectedOption={selectedOption}
              onClick={this.onMoveTypeSelected}
              maxWeight={entitlement.sum}
            />
          </Fragment>
        )}
      </div>
    );
  }
}

PpmSize.propTypes = {
  pendingPpmSize: PropTypes.string,
  currentPpm: PropTypes.shape({ id: PropTypes.string, size: PropTypes.string }),
  setPendingPpmSize: PropTypes.func.isRequired,
  entitlement: PropTypes.object,
};

function mapStateToProps(state) {
  return {
    ...state.ppm,
    entitlement: loadEntitlements(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ setPendingPpmSize }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmSize);
