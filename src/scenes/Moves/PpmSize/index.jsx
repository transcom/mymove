import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './PpmSize.css';

import { createPpm } from './ducks';

export class BigButton extends Component {
  render() {
    let className = 'ppm-size-button';
    if (this.props.selected) {
      className += ' selected';
    }
    return (
      <div className={className} onClick={this.props.onButtonClick}>
        <p>{this.props.firstLine}</p>
        <p>{this.props.secondLine}</p>
        <img className="icon" src={this.props.icon} />
      </div>
    );
  }
}

BigButton.propTypes = {
  firstLine: PropTypes.string.isRequired,
  secondLine: PropTypes.string.isRequired,
  icon: PropTypes.string.isRequired,
  selected: PropTypes.bool,
  onButtonClick: PropTypes.func,
};

export class BigButtonGroup extends Component {
  render() {
    var createButton = (value, firstLine, secondLine, icon) => {
      var onButtonClick = () => {
        this.props.onMoveTypeSelected(value);
      };
      return (
        <BigButton
          value={value}
          firstLine={firstLine}
          secondLine={secondLine}
          icon={icon}
          selected={this.props.selectedOption === value}
          onButtonClick={onButtonClick}
        />
      );
    };
    var small = createButton(
      'S',
      'A few items in your car?',
      '(approx 100 - 800 lbs)',
      carGray,
    );
    var medium = createButton(
      'M',
      'A trailer full of household goods?',
      '(approx 400 - 1,200 lbs)',
      trailerGray,
    );
    var large = createButton(
      'L',
      'A moving truck that you rent yourself?',
      '(approx 1,000 - 5,000 lbs)',
      truckGray,
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
  onMoveTypeSelected: PropTypes.func,
};

export class PpmSize extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Size Selection';
  }
  constructor(props) {
    super(props);
    this.state = {
      selectedOption: null,
    };
  }
  onMoveTypeSelected = value => {
    this.setState({ selectedOption: value });
  };
  render() {
    console.log('selected value', this.state.selectedOption);
    return (
      <div className="usa-grid-full ppm-size-content">
        <h3>How much of your stuff do you intend to move yourself?</h3>
        <BigButtonGroup
          selectedOption={this.state.selectedOption}
          onMoveTypeSelected={this.onMoveTypeSelected}
        />
      </div>
    );
  }
}

PpmSize.propTypes = {
  createPpm: PropTypes.func.isRequired,
  currentPpm: PropTypes.object,
  match: PropTypes.object.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return state.ppm;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createPpm }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmSize);
