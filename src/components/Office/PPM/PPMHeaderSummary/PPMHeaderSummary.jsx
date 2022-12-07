import React from 'react';
import { number } from 'prop-types';
import { Label } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './PPMHeaderSummary.module.scss';

import { PPMShipmentShape } from 'types/shipment';
import { formatCentsTruncateWhole } from 'utils/formatters';

export default function PPMHeaderSummary({ ppmShipment, ppmNumber }) {
  const {
    actualPickupPostalCode,
    actualDestinationPostalCode,
    actualMoveDate,
    hasReceivedAdvance,
    advanceAmountReceived,
  } = ppmShipment || {};

  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        <section>
          <div>
            <Label className={styles.headerLabel}>Departure date</Label>
            <span className={styles.light}>{actualMoveDate}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Starting ZIP</Label>
            <span className={styles.light}>{actualPickupPostalCode}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Ending ZIP</Label>
            <span className={styles.light}>{actualDestinationPostalCode}</span>
          </div>
          <div>
            <Label className={styles.headerLabel}>Advance received</Label>
            <span className={styles.light}>
              {hasReceivedAdvance ? `Yes, $${formatCentsTruncateWhole(advanceAmountReceived)}` : 'No'}
            </span>
          </div>
        </section>
      </div>
    </header>
  );
}

PPMHeaderSummary.propTypes = {
  ppmShipment: PPMShipmentShape,
  ppmNumber: number.isRequired,
};

PPMHeaderSummary.defaultProps = {
  ppmShipment: undefined,
};
