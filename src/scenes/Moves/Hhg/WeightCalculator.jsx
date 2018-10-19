import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NumericInput from 'react-numeric-input';
import Alert from 'shared/Alert';
import helpIcon from 'shared/images/help-icon.png';

const StuffSize = Object.freeze({
  LESS: 'less',
  AVERAGE: 'average',
  MORE: 'more',
  PACKRAT: 'packrat',
});

class WeightCalculator extends Component {
  constructor(props) {
    super(props);
    this.state = {
      rooms: 1,
      stuff: StuffSize.AVERAGE,
      showInfo: false,
    };
  }

  roomsUpdated = valueAsNumber => {
    this.setState(function(state) {
      this.estimateAndReport(valueAsNumber, state.stuff);
      return {
        rooms: valueAsNumber,
      };
    });
  };

  stuffUpdated = e => {
    const newStuffValue = e.target.value;
    this.setState(function(state) {
      this.estimateAndReport(state.rooms, newStuffValue);
      return {
        stuff: newStuffValue,
      };
    });
  };

  openInfo = e => {
    e.preventDefault();
    this.setState({ showInfo: true });
  };

  closeInfo = e => {
    e.preventDefault();
    this.setState({ showInfo: false });
  };

  estimateAndReport = (rooms, stuff) => {
    let pounds = 0;
    switch (stuff) {
      case StuffSize.LESS: {
        pounds = rooms * 1000;
        break;
      }
      case StuffSize.AVERAGE: {
        pounds = rooms * 1200;
        break;
      }
      case StuffSize.MORE: {
        pounds = rooms * 1500;
        break;
      }
      case StuffSize.PACKRAT: {
        pounds = rooms * 1800;
        break;
      }
      default: {
        console.log('Invalid StuffSize: ' + stuff);
        break;
      }
    }

    this.props.onEstimate(pounds);
  };

  render() {
    return (
      <div className="weight-calculator">
        <h4>Get a quick weight estimate so you don't go over your weight entitlement</h4>
        <div className="weight-calculator-contents">
          <div className="usa-input">
            <label className="usa-input-label" htmlFor="rooms">
              How many furnished rooms in your current home?<br />
              <span className="weight-calculator-help">
                Include bedrooms, living, dining and family rooms, offices and a garage or offsite storage if you store
                a lot there.
              </span>
            </label>
            <NumericInput name="rooms" min={1} value={this.state.rooms} onChange={this.roomsUpdated} />
          </div>
          <div className="usa-input rounded">
            <label className="usa-input-label" htmlFor="stuff">
              How much stuff do you have in your rooms?{' '}
              <a href="" onClick={this.openInfo}>
                <img className="help-icon" src={helpIcon} alt="onto-move-orders" />
              </a>
              {this.state.showInfo && (
                <Alert type="info" heading="">
                  Do you have alot of heavy stuff in your home or not so much? Take your best guess. Below average means
                  you have ~ 1,000 lbs per room; Average = 1,200 lbs; Above average = 1,500 lbs; Packrat = 1,800 lbs.{' '}
                  <a href="" onClick={this.closeInfo}>
                    Close
                  </a>
                </Alert>
              )}
            </label>
            <select name="stuff" value={this.state.stuff} onChange={this.stuffUpdated}>
              <option value={StuffSize.LESS}>Less than average</option>
              <option value={StuffSize.AVERAGE}>Average Amount</option>
              <option value={StuffSize.MORE}>More than average</option>
              <option value={StuffSize.PACKRAT}>I'm a packrat who collects heavy objects!</option>
            </select>
          </div>
        </div>
      </div>
    );
  }
}

WeightCalculator.propTypes = {
  onEstimate: PropTypes.func.isRequired,
};

export default WeightCalculator;
