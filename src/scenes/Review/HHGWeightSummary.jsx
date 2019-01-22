import React, { Fragment } from 'react';

import { formatNumber } from 'shared/formatters';

export default function HHGWeightSummary(props) {
  const { shipment, entitlements } = props;

  let progearSummary = '';
  if (shipment.progear_weight_estimate > 0) {
    progearSummary = ` + ${formatNumber(shipment.progear_weight_estimate)} lbs pro-gear`;
  }

  let spouseProgearSummary = '';
  if (shipment.spouse_progear_weight_estimate > 0) {
    spouseProgearSummary = ` + ${formatNumber(shipment.spouse_progear_weight_estimate)} lbs spouse pro-gear`;
  }

  let congrats = '';
  if (
    shipment.weight_estimate <= entitlements.weight &&
    shipment.progear_weight_estimate <= entitlements.pro_gear &&
    shipment.spouse_progear_weight_estimate <= entitlements.pro_gear_spouse
  ) {
    congrats = 'Great! You appear within your weight allowance.';
  }
  return (
    <Fragment>
      {formatNumber(shipment.weight_estimate)} lbs
      {progearSummary}
      {spouseProgearSummary}
      <br /> {congrats}
    </Fragment>
  );
}
