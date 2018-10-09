import React, { Fragment } from 'react';

import { formatNumber } from 'shared/formatters';
import Alert from 'shared/Alert';

function warning(name, estimatedWeight, allowedWeight) {
  const remaining = allowedWeight - estimatedWeight;

  if (remaining >= 0) {
    return null;
  }

  return (
    <Alert type="warning" heading="">
      Your {name} of {formatNumber(estimatedWeight)} is{' '}
      {formatNumber(-remaining)} lbs over your maximum entitlement of{' '}
      {formatNumber(allowedWeight)} lbs.
    </Alert>
  );
}

export default function HHGWeightSummary(props) {
  const { shipment, entitlements } = props;

  return (
    <Fragment>
      {warning(
        'weight estimate',
        shipment.weight_estimate,
        entitlements.weight,
      )}
      {warning(
        'pro-gear weight estimate',
        shipment.progear_weight_estimate,
        entitlements.pro_gear,
      )}
      {warning(
        'spouse pro-gear weight estimate',
        shipment.spouse_progear_weight_estimate,
        entitlements.pro_gear_spouse,
      )}
    </Fragment>
  );
}
