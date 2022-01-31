import React from 'react';
import { string, arrayOf, shape, number } from 'prop-types';
import { Link } from 'react-router-dom';

import styles from './ExternalVendorWeightSummary.module.scss';

import { formatWeight } from 'utils/formatters';

const totalWeight = (shipments) => {
  if (shipments.length > 1) {
    return formatWeight(
      shipments.reduce((prev, curr) => {
        return prev + curr.ntsRecordedWeight;
      }, 0),
    );
  }
  return formatWeight(shipments[0].ntsRecordedWeight);
};

export default function ExternalVendorWeightSummary({ shipments }) {
  return (
    <div className={styles.ExternalVendorWeightSummary}>
      <p className="text-bold">{shipments.length > 1 ? `${shipments.length} other shipments:` : '1 other shipment:'}</p>
      <p>{totalWeight(shipments)}</p>
      <Link to="details">View move details</Link>
    </div>
  );
}

ExternalVendorWeightSummary.propTypes = {
  shipments: arrayOf(
    shape({
      id: string.isRequired,
      shipmentType: string.isRequired,
      reweigh: shape({ id: string.isRequired, weight: number }),
    }),
  ).isRequired,
};
