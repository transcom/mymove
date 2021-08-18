import React from 'react';
import { string, arrayOf, shape } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import ShipmentList from '../../ShipmentList';

import styles from './BillableWeightCard.module.scss';

import { formatWeight } from 'shared/formatters';

export default function BillableWeightCard({
  maxBillableWeight,
  weightRequested,
  weightAllowance,
  totalBillableWeight,
  shipments,
  entitlements,
}) {
  return (
    <div className={classnames(styles.cardContainer, 'container')}>
      <div className={styles.cardHeader}>
        <h2>Billable weights</h2>
        <Button>Review weights</Button>
      </div>
      <div className={styles.spaceBetween}>
        <div>
          <h5>Maximum billable weight</h5>
          <h4>{formatWeight(maxBillableWeight)}</h4>
          <h6>
            Weight requested<strong>{formatWeight(weightRequested)}</strong>
          </h6>
          <h6>
            Weight allowance<strong>{formatWeight(weightAllowance)}</strong>
          </h6>
        </div>
        <div className={styles.shipmentSection}>
          <h5>Total billable weight</h5>
          <h4>{formatWeight(totalBillableWeight)}</h4>
          <div className={styles.shipmentList}>
            <ShipmentList shipments={shipments} entitlements={entitlements} showShipmentWeight />
          </div>
        </div>
      </div>
    </div>
  );
}

BillableWeightCard.propTypes = {
  maxBillableWeight: string.isRequired,
  weightRequested: string.isRequired,
  weightAllowance: string.isRequired,
  totalBillableWeight: string.isRequired,
  shipments: arrayOf(
    shape({
      id: string.isRequired,
      shipmentType: string.isRequired,
      reweigh: shape({ id: string.isRequired, weight: string }),
    }),
  ).isRequired,
  entitlements: arrayOf(shape({ id: string.isRequired, authorizedWeight: string.isRequired })).isRequired,
};
