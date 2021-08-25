import React from 'react';
import { Button } from '@trussworks/react-uswds';

import styles from './ShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default function ShipmentCard() {
  return (
    <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG} className={styles.container}>
      <header>
        <h2>HHG</h2>
        <section>
          <span>
            <strong>Departed</strong> 09 Jun 2021
          </span>
          <span>
            <strong>From</strong> Rancho Santa Margarita, CA 92688
          </span>
          <span>
            <strong>To</strong> West Springfield Town, MA 01089
          </span>
        </section>
      </header>
      <main>
        <div className={styles.field}>
          <strong>Estimated weight</strong>
          <span> 5,000 lbs</span>
        </div>
        <div className={styles.field}>
          <strong>Original weight</strong>
          <span> 4,014 lbs</span>
        </div>
        <div className={`${styles.field} ${styles.missing}`}>
          <strong>Reweigh weight</strong>
          <span>
            <strong>Missing</strong>
          </span>
        </div>
        <div className={styles.field}>
          <strong>Date reweigh requested</strong>
          <span> 10 Jun 2021</span>
        </div>
        <div className={`${styles.field} ${styles.remarks}`}>
          <strong>Reweigh remarks</strong>
          <span>Unable to perform reweigh because shipment was already unloaded</span>
        </div>
      </main>
      <footer>
        <h3>Billable weight</h3>
        <span>4,014 lbs</span>
        <Button className={styles.editBtn}>Edit</Button>
      </footer>
    </ShipmentContainer>
  );
}
