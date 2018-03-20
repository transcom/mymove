import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './PpmSize.css';

import createPpm from './ducks';

function BigButton(props) {
  return <div className="size-button">{props.children}</div>;
}

export class PpmSize extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Size Selection';
  }

  handleOptionChange = changeEvent => {
    this.setState({
      selectedOption: changeEvent.target.value,
    });
  };

  render() {
    return (
      <div className="usa-grid-full ppm-size-content">
        <h3>How much of your stuff do you intend to move yourself?</h3>

        <div className="usa-width-one-third">
          <label className="container">
            <input type="radio" value="small" name="size-selector" />
            <BigButton>
              <p>A few items in your car?</p>
              <p>(approx 100 - 800 lbs)</p>
              <img src={carGray} alt="car-gray" />
            </BigButton>
          </label>
        </div>

        <div className="usa-width-one-third">
          <label className="container">
            <input type="radio" value="medium" name="size-selector" />
            <BigButton>
              <p>A trailer full of household goods? </p>
              <p>(approx 400 - 1,200 lbs)</p>
              <img src={trailerGray} alt="trailer-gray" />
            </BigButton>
          </label>
        </div>

        <div className="usa-width-one-third">
          <label className="container">
            <input type="radio" value="large" name="size-selector" />
            <BigButton>
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
