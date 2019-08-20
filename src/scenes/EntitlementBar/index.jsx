import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { mapValues } from 'lodash';
import Alert from 'shared/Alert';

import './index.css';

class EntitlementBar extends Component {
  state = { showInfo: false };

  openInfo = () => {
    this.setState({ showInfo: true });
  };

  closeInfo = () => {
    this.setState({ showInfo: false });
  };

  render() {
    const props = this.props;

    const entitlementText = mapValues(props.entitlement, i => i.toLocaleString());
    const ppmSummaryHtml = () => {
      if (!props.entitlement) return <p />;
      if (props.entitlement.pro_gear_spouse > 0)
        return (
          <p>
            {entitlementText.weight} lbs. + {entitlementText.pro_gear} lbs. of pro-gear +{' '}
            {entitlementText.pro_gear_spouse} lbs. of spouse's pro-gear = <strong>{entitlementText.sum} lbs.</strong>
          </p>
        );
      return (
        <p>
          {entitlementText.weight} lbs. + {entitlementText.pro_gear} lbs. of pro-gear ={' '}
          <strong>{entitlementText.sum} lbs.</strong>
        </p>
      );
    };

    const hhgPPMSummaryHtml = () => {
      return (
        <p>
          {props.hhgPPMEntitlementMessage} <a onClick={this.openInfo}>What's this?</a>
        </p>
      );
    };

    return (
      <div>
        <div className="entitlement-container">
          <p>
            <strong>How much weight are you entitled to move?</strong>
          </p>
          {!props.hhgPPMEntitlementMessage ? ppmSummaryHtml() : hhgPPMSummaryHtml()}
        </div>
        {this.state.showInfo && (
          <div className="usa-width-one-whole top-bottom-buffered">
            <Alert type="info" className="usa-width-one-whole" heading="">
              Your entitlement represents the maximum weight the military is willing to move for you and/or pay you to
              move yourself, including any pro-gear. If you carry over this amount, you will not be paid for any excess
              weight. Pro-gear is any gear you need to perform your official duties at your next or later destination,
              such as reference materials, tools for a technician or mechanic or specialized clothing that's not a
              typical uniform (such as diving or astronaut suits). <a onClick={this.closeInfo}>Close</a>
            </Alert>
          </div>
        )}
      </div>
    );
  }
}
EntitlementBar.propTypes = {
  entitlement: PropTypes.shape({
    weight: PropTypes.number.isRequired,
    pro_gear: PropTypes.number.isRequired,
    pro_gear_spouse: PropTypes.number.isRequired,
    sum: PropTypes.number.isRequired,
    storage_in_transit: PropTypes.number.isRequired,
  }),
  hhgPPMEntitlementMessage: PropTypes.string,
};
export default EntitlementBar;
