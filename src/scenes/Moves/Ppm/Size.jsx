import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { setPendingPpmSize } from './ducks';

import BigButton from 'shared/BigButton';
import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './Size.css';

class BigButtonGroup extends Component {
  render() {
    var createButton = (value, firstLine, secondLine, icon, altTag) => {
      var onButtonClick = () => {
        this.props.onClick(value);
      };
      return (
        <BigButton
          value={value}
          selected={this.props.selectedOption === value}
          onClick={onButtonClick}
        >
          <p>{firstLine}</p>
          <p className="Todo">{secondLine}</p>
          <img className="icon" src={icon} alt={altTag} />
        </BigButton>
      );
    };
    var small = createButton(
      'S',
      'A few items in your car?',
      '(approx 100 - 800 lbs)',
      carGray,
      'car-gray',
    );
    var medium = createButton(
      'M',
      'A trailer full of household goods?',
      '(approx 400 - 1,200 lbs)',
      trailerGray,
      'trailer-gray',
    );
    var large = createButton(
      'L',
      'A moving truck that you rent yourself?',
      '(approx 1,000 - 5,000 lbs)',
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
  componentDidMount() {
    document.title = 'Transcom PPP: Size Selection';
  }

  onMoveTypeSelected = value => {
    this.props.setPendingPpmSize(value);
  };
  render() {
    const { pendingPpmSize, currentPpm } = this.props;
    const selectedOption = pendingPpmSize || (currentPpm && currentPpm.size);
    return (
      <div className="usa-grid-full ppm-size-content">
        <h3>How much of your stuff do you intend to move yourself?</h3>
        <BigButtonGroup
          selectedOption={selectedOption}
          onClick={this.onMoveTypeSelected}
        />
      </div>
    );
  }
}

PpmSize.propTypes = {
  pendingPpmSize: PropTypes.string,
  currentPpm: PropTypes.shape({ id: PropTypes.string, size: PropTypes.string }),
  setPendingPpmSize: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return state.ppm;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ setPendingPpmSize }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmSize);
