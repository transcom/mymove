import React from 'react';
import { string, number, shape } from 'prop-types';
import moment from 'moment';
import { Button } from '@trussworks/react-uswds';

import styles from './ShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight, formatAddressShort } from 'shared/formatters';

export default function ShipmentCard({
  billableWeight,
  dateReweighRequested,
  departedDate,
  pickupAddress,
  destinationAddress,
  estimatedWeight,
  originalWeight,
  reweighRemarks,
  reweightWeight,
}) {
  return (
    <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG} className={styles.container}>
      <header>
        <h2>HHG</h2>
        <section>
          <span>
            {/* <strong>Departed</strong> 09 Jun 2021 */}
            <strong>Departed</strong> {moment(departedDate).format('DD MMM YYYY')}
          </span>
          <span>
            {/* <strong>From</strong> Rancho Santa Margarita, CA 92688 */}
            <strong>From</strong> {formatAddressShort(pickupAddress)}
          </span>
          <span>
            <strong>To</strong> {formatAddressShort(destinationAddress)}
          </span>
        </section>
      </header>
      <main>
        <div className={styles.field}>
          <strong>Estimated weight</strong>
          {/* <span> 5,000 lbs</span> */}
          <span>{formatWeight(estimatedWeight)}</span>
        </div>
        <div className={styles.field}>
          <strong>Original weight</strong>
          {/* <span> 4,014 lbs</span> */}
          <span>{formatWeight(originalWeight)}</span>
        </div>
        <div className={`${styles.field} ${styles.missing}`}>
          <strong>Reweigh weight</strong>
          <span>{reweightWeight ? formatWeight : <strong>Missing</strong>}</span>
        </div>
        <div className={styles.field}>
          <strong>Date reweigh requested</strong>
          {/* <span> 10 Jun 2021</span> */}
          <span>{moment(dateReweighRequested).format('DD MMM YYYY')}</span>
        </div>
        <div className={`${styles.field} ${styles.remarks}`}>
          <strong>Reweigh remarks</strong>
          {/* <span>Unable to perform reweigh because shipment was already unloaded</span> */}
          <span>{reweighRemarks}</span>
        </div>
      </main>
      <footer>
        <h3>Billable weight</h3>
        {/* <span>4,014 lbs</span> */}
        <span>{formatWeight(billableWeight)}</span>
        <Button className={styles.editBtn}>Edit</Button>
      </footer>
    </ShipmentContainer>
  );
}

ShipmentCard.propTypes = {
  billableWeight: number.isRequired,
  dateReweighRequested: string.isRequired,
  departedDate: string.isRequired,
  destinationAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    postal_code: string.isRequired,
  }).isRequired,
  estimatedWeight: number.isRequired,
  originalWeight: number.isRequired,
  pickupAddress: shape({
    city: string.isRequired,
    state: string.isRequired,
    postal_code: string.isRequired,
  }).isRequired,
  reweighRemarks: string.isRequired,
  reweightWeight: number.isRequired,
};
