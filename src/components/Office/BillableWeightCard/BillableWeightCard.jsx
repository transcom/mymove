import React from 'react';
import { string } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import styles from './BillableWeightCard.module.scss';

export default function BillableWeightCard({
  maxBillableWeight,
  weightRequested,
  weightAllowance,
  totalBillableWeight,
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
          <h4>{maxBillableWeight} lbs</h4>
          <h6>
            Weight requested<strong>{weightRequested} lbs</strong>
          </h6>
          <h6>
            Weight allowance<strong>{weightAllowance} lbs</strong>
          </h6>
        </div>
        <div className={styles.shipmentSection}>
          <h5>Total billable weight</h5>
          <h4>{totalBillableWeight} lbs</h4>
          <div className={styles.shipmentPlaceholder}>shipment list placeholder</div>
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
};
