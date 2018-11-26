import React from 'react';
import PropTypes from 'prop-types';
import { mapValues } from 'lodash';
import './index.css';

const EntitlementBar = props => {
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
    return <p>{props.hhgPPMEntitlementMessage}</p>;
  };

  return (
    <div className="entitlement-container">
      <p>
        <strong>How much are you entitled to move?</strong>
      </p>
      {!props.hhgPPMEntitlementMessage ? ppmSummaryHtml() : hhgPPMSummaryHtml()}
    </div>
  );
};
EntitlementBar.propTypes = {
  entitlement: PropTypes.shape({
    weight: PropTypes.number.isRequired,
    pro_gear: PropTypes.number.isRequired,
    pro_gear_spouse: PropTypes.number.isRequired,
    sum: PropTypes.number.isRequired,
  }),
  hhgPPMEntitlementMessage: PropTypes.string,
};
export default EntitlementBar;
