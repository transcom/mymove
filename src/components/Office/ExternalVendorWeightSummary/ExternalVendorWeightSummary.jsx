import React from 'react';
import { arrayOf, shape, number, string } from 'prop-types';
import { Link } from 'react-router-dom';

import styles from './ExternalVendorWeightSummary.module.scss';

import { formatWeight } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const totalWeight = (shipments) => {
  if (shipments.length > 1) {
    return formatWeight(
      shipments.reduce((prev, curr) => {
        // NTS shipments won't have a recorded weight, so just return existing total in that case
        return curr.ntsRecordedWeight ? prev + curr.ntsRecordedWeight : prev;
      }, 0),
    );
  }
  return formatWeight(shipments[0].ntsRecordedWeight);
};

export default function ExternalVendorWeightSummary({ shipments }) {
  return (
    <div className={styles.ExternalVendorWeightSummary}>
      <p className="text-bold">{shipments.length > 1 ? `${shipments.length} other shipments:` : '1 other shipment:'}</p>
      {shipments.some((s) => s.shipmentType === SHIPMENT_OPTIONS.NTSR) ? <p>{totalWeight(shipments)}</p> : null}
      <Link to="details">View move details</Link>
    </div>
  );
}

ExternalVendorWeightSummary.propTypes = {
  shipments: arrayOf(
    shape({
      ntsRecordedWeight: number,
      shipmentType: string.isRequired,
    }),
  ).isRequired,
};
