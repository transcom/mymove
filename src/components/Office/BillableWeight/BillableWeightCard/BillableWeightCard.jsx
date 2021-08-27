import React from 'react';
import { string, arrayOf, shape, number } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import styles from './BillableWeightCard.module.scss';

import ShipmentList from 'components/ShipmentList';
import { formatWeight } from 'shared/formatters';

export default function BillableWeightCard({
  maxBillableWeight,
  weightRequested,
  weightAllowance,
  totalBillableWeight,
  shipments,
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
            <ShipmentList shipments={shipments} showShipmentWeight />
          </div>
        </div>
      </div>
    </div>
  );
}

BillableWeightCard.propTypes = {
  maxBillableWeight: number.isRequired,
  weightRequested: number.isRequired,
  weightAllowance: number.isRequired,
  totalBillableWeight: number.isRequired,
  shipments: arrayOf(
    shape({
      id: string.isRequired,
      shipmentType: string.isRequired,
      reweigh: shape({ id: string.isRequired, weight: number }),
    }),
  ).isRequired,
};
