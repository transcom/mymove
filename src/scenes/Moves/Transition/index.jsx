import React, { Component } from 'react';
import truckGrayCheckGray from './truck-gray-check-gray.svg';
import ppmBlack from './ppm-black.svg';
import './Transition.css';

export class Transition extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Transition Page';
  }
  render() {
    return (
      <div className="transition-box">
        <div className="hhg-icon">
          <img src={truckGrayCheckGray} alt={this.props.match.path} />
          <b>
            <p> Shipment 1 (HHG)</p>
          </b>
        </div>

        <div className="transition-text">
          <b>
            <p>Your moving company shipment is now set up.</p>
          </b>
          <p>Let’s go on to the stuff you’ll move yourself.</p>
        </div>

        <div className="ppm-icon">
          <img src={ppmBlack} alt={this.props.match.path} />
          <b>
            <p> Shipment 2 (PPM)</p>
          </b>
        </div>
      </div>
    );
  }
}

export default Transition;
