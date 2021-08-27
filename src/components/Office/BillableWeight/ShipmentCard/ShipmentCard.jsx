import React from 'react';
import { string, number, shape } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import styles from './ShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight, formatAddressShort, formatDateFromIso } from 'shared/formatters';

export default function ShipmentCard({
  billableWeight,
  dateReweighRequested,
  departedDate,
  pickupAddress,
  destinationAddress,
  estimatedWeight,
  originalWeight,
  reweighRemarks,
  reweighWeight,
}) {
  return (
    <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG} className={styles.container}>
      <header>
        <h2>HHG</h2>
        <section>
          <span>
            <strong>Departed</strong> {formatDateFromIso(departedDate, 'DD MMM YYYY')}
          </span>
          <span>
            <strong>From</strong> {formatAddressShort(pickupAddress)}
          </span>
          <span>
            <strong>To</strong> {formatAddressShort(destinationAddress)}
          </span>
        </section>
      </header>
      <div className={styles.weights}>
        <div className={styles.field}>
          <strong>Estimated weight</strong>
          <span>{formatWeight(estimatedWeight)}</span>
        </div>
        <div className={styles.field}>
          <strong>Original weight</strong>
          <span>{formatWeight(originalWeight)}</span>
        </div>
        <div
          className={classnames(styles.field, {
            [styles.missing]: !reweighWeight,
          })}
        >
          <strong>Reweigh weight</strong>
          <span>{reweighWeight ? formatWeight(reweighWeight) : <strong>Missing</strong>}</span>
        </div>
        <div className={styles.field}>
          <strong>Date reweigh requested</strong>
          <span>{formatDateFromIso(dateReweighRequested, 'DD MMM YYYY')}</span>
        </div>
        <div className={classnames(styles.field, styles.remarks)}>
          <strong>Reweigh remarks</strong>
          <span>{reweighRemarks}</span>
        </div>
      </div>
      <footer>
        <h3>Billable weight</h3>
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
  reweighWeight: number,
};

ShipmentCard.defaultProps = {
  reweighWeight: null,
};
