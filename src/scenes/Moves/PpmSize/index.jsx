import React, { Component } from 'react';
import carGray from './car-gray.svg';
import trailerGray from './trailer-gray.svg';
import truckGray from './truck-gray.svg';
import './PpmSize.css';

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

export default PpmSize;
