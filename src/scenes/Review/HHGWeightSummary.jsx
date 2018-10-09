import React, { Fragment } from 'react';

import { formatNumber } from 'shared/formatters';

export default function HHGWeightSummary(props) {
  const { shipment } = props;

  let progearSummary = '';
  if (shipment.progear_weight_estimate > 0) {
    progearSummary = ` + ${formatNumber(
      shipment.progear_weight_estimate,
    )} lbs pro-gear`;
  }

  let spouseProgearSummary = '';
  if (shipment.spouse_progear_weight_estimate > 0) {
    spouseProgearSummary = ` + ${formatNumber(
      shipment.spouse_progear_weight_estimate,
    )} lbs spouse pro-gear`;
  }

  return (
    <Fragment>
      {formatNumber(shipment.weight_estimate)} lbs
      {progearSummary}
      {spouseProgearSummary}
    </Fragment>
  );
}
