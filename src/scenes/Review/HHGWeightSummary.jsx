import React, { Fragment } from 'react';

import { formatNumber } from 'shared/formatters';
import Alert from 'shared/Alert';

export default function HHGShipmentSummary(props) {
  const { shipment, entitlements } = props;

  let progearSummary = '';
  if (shipment.progear_weight_estimate > 0) {
    progearSummary = ` + ${formatNumber(
      shipment.progear_weight_estimate,
    )} lbs pro-gear`;
  }

  let spouseProgearSummary = '';
  if (shipment.progear_weight_estimate > 0) {
    spouseProgearSummary = ` + ${formatNumber(
      shipment.progear_weight_estimate,
    )} lbs spouse pro-gear`;
  }

  const weightRemaining = entitlements.weight - shipment.weight_estimate;
  const progearRemaining =
    entitlements.pro_gear - shipment.progear_weight_estimate;
  const spouseProgearRemaining =
    entitlements.pro_gear_spouse - shipment.spouse_progear_weight_estimate;

  //   const progearExceeded =
  //     shipment.progear_weight_estimate > entitlements.pro_gear ||
  //     shipment.spouse_progear_weight_estimate > entitlements.pro_gear_spouse;

  //   let weightMessage;
  //   if (entitlementsExceeded) {
  //     weightMessage = <Alert type="error" heading="An error occurred" />;
  //   } else {
  //     weightMessage = 'Great! You appear to be within your weight allowance.';
  //   }

  return (
    <Fragment>
      {formatNumber(shipment.weight_estimate)} lbs
      {progearSummary}
      {spouseProgearSummary}
      <br />
      {weightRemaining < 0 && (
        <Alert type="warning" heading="">
          Your estimate of {formatNumber(shipment.weight_estimate)} is{' '}
          {formatNumber(-weightRemaining)} lbs over your maximum entitlement of{' '}
          {formatNumber(entitlements.weight)} lbs.
        </Alert>
      )}
    </Fragment>
  );
}
