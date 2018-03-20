import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './PpmSize.css';

import { createPpm } from './ducks';

export class BigButtonGroup extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedOption: null,
    };
  }
  render() {
    return React.Children.map(props.children, child => {
      if (child.type === BigButton)
        return React.cloneElement(child, {
          isSelected: props.value === child.props.value,
          name: props.name,
          onChange: props.handleChange,
        });
      return child;
    });
  }
}

BigButtonGroup.propTypes = {
  children: PropTypes.node,
  value: PropTypes.object,
  handleChange: PropTypes.func.isRequired,
  name: PropTypes.string.isRequired,
};

function BigButton(props) {
  return (
    <div onClick={props.onChange} className="size-button">
      {props.children}
    </div>
  );
}
BigButton.propTypes = {
  children: PropTypes.node,
  value: PropTypes.object.isRequired,
  name: PropTypes.string,
  isSelected: PropTypes.bool.isRequired,
};
BigButton.defaultProps = {
  isSelected: false,
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

  handleOptionChange = evt => {
    this.setState({
      selectedOption: evt.currentTarget.value,
    });
  };

  render() {
    return (
      <div className="usa-grid-full ppm-size-content">
        <h3>How much of your stuff do you intend to move yourself?</h3>

        <BigButtonGroup
          value={this.state.selectedOption}
          handleChange={this.handleOptionChange}
          name="foo"
        >
          <BigButton value="S"> S </BigButton>
          <BigButton value="M"> M </BigButton>
          <BigButton value="L"> L </BigButton>
        </BigButtonGroup>
        <div className="usa-width-one-third">
          <label className="container">
            <input type="radio" value="small" name="size-selector" />
            <BigButton value="S">
              <p>A few items in your car?</p>
              <p>(approx 100 - 800 lbs)</p>
              <img src={carGray} alt="car-gray" />
            </BigButton>
          </label>
        </div>

        <div className="usa-width-one-third">
          <label className="container">
            <input type="radio" value="medium" name="size-selector" />
            <BigButton value="M">
              <p>A trailer full of household goods? </p>
              <p>(approx 400 - 1,200 lbs)</p>
              <img src={trailerGray} alt="trailer-gray" />
            </BigButton>
          </label>
        </div>

        <div className="usa-width-one-third">
          <label className="container">
            <input type="radio" value="large" name="size-selector" />
            <BigButton value="L">
              <p>A moving truck that you rent yourself?</p>
              <p>(approx 1,000 - 5,000 lbs)</p>
              <img src={truckGray} alt="truck-gray" />
            </BigButton>
          </label>
        </div>
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
  return state.ppmSize;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createPpm }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PpmSize);
