import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Label, Button, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './HeaderSection.module.scss';

import { formatDate, formatCents, formatWeight } from 'utils/formatters';

export const sectionTypes = {
  incentives: 'incentives',
  shipmentInfo: 'shipmentInfo',
  incentiveFactors: 'incentiveFactors',
};

const getSectionTitle = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.incentives:
      return `Incentives/Costs`;
    case sectionTypes.shipmentInfo:
      return `Shipment Info`;
    case sectionTypes.incentiveFactors:
      return `Incentive Factors`;
    default:
      return <Alert>`Error getting section title!`</Alert>;
  }
};

// Returns the markup needed for a specific section
const getSectionMarkup = (sectionInfo) => {
  const aoaRequestedValue = sectionInfo.isAdvanceRequested
    ? `$${formatCents(sectionInfo.advanceAmountRequested)}`
    : 'No';
  const aoaValue = sectionInfo.isAdvanceReceived ? `$${formatCents(sectionInfo.advanceAmountReceived)}` : 'No';

  switch (sectionInfo.type) {
    case sectionTypes.shipmentInfo:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Planned Move Start Date</Label>
            <span className={styles.light}>{formatDate(sectionInfo.plannedMoveDate, null, 'DD-MMM-YYYY')}</span>
          </div>
          <div>
            <Label>Actual Move Start Date</Label>
            <span className={styles.light}>{formatDate(sectionInfo.actualMoveDate, null, 'DD-MMM-YYYY')}</span>
          </div>
          <div>
            <Label>Starting Address</Label>
            <span className={styles.light}>{sectionInfo.pickupAddress}</span>
          </div>
          <div>
            <Label>Ending Address</Label>
            <span className={styles.light}>{sectionInfo.destinationAddress}</span>
          </div>
          <div>
            <Label>Miles</Label>
            <span className={styles.light}>{sectionInfo.miles}</span>
          </div>
          <div>
            <Label>Estimated Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.estimatedWeight)}</span>
          </div>
          <div>
            <Label>Actual Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.actualWeight)}</span>
          </div>
        </div>
      );

    case sectionTypes.incentives:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Government Constructed Cost (GCC)</Label>
            <span className={styles.light}>${formatCents(sectionInfo.gcc)}</span>
          </div>
          <div>
            <Label>Gross Incentive</Label>
            <span className={styles.light}>${formatCents(sectionInfo.grossIncentive)}</span>
          </div>
          <div>
            <Label>Advance Requested</Label>
            <span className={styles.light}>{aoaRequestedValue}</span>
          </div>
          <div>
            <Label>Advance Received</Label>
            <span className={styles.light}>{aoaValue}</span>
          </div>
          <div>
            <Label>Remaining Incentive</Label>
            <span className={styles.light}>${formatCents(sectionInfo.remainingIncentive)}</span>
          </div>
        </div>
      );

    case sectionTypes.incentiveFactors:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Haul Price</Label>
            <span className={styles.light}>${formatCents(sectionInfo.haulPrice)}</span>
          </div>
          <div>
            <Label>Haul Fuel Surcharge</Label>
            <span className={styles.light}>
              {sectionInfo.haulFSC < 0 ? '-$' : '$'}
              {formatCents(Math.abs(sectionInfo.haulFSC))}
            </span>
          </div>
          <div>
            <Label>Full Pack/Unpack Charge</Label>
            <span className={styles.light}>${formatCents(sectionInfo.fullPackUnpackCharge)}</span>
          </div>
          <div>
            <Label>Origin Price</Label>
            <span className={styles.light}>${formatCents(sectionInfo.dop)}</span>
          </div>
          <div>
            <Label>Destination Price</Label>
            <span className={styles.light}>${formatCents(sectionInfo.ddp)}</span>
          </div>
        </div>
      );

    default:
      return <Alert>An error occured while getting section markup!</Alert>;
  }
};

export default function PPMHeaderSummary({ sectionInfo }) {
  const [showDetails, setShowDetails] = useState(false);
  const showRequestDetailsButton = true;
  const handleToggleDetails = () => setShowDetails((prevState) => !prevState);
  const showDetailsChevron = showDetails ? 'chevron-up' : 'chevron-down';
  const showDetailsText = showDetails ? 'Hide details' : 'Show details';

  return (
    <section className={classnames(styles.HeaderSection)}>
      <header>
        <h4>{getSectionTitle(sectionInfo)}</h4>
      </header>
      <div className={styles.toggleDrawer}>
        {showRequestDetailsButton && (
          <Button
            aria-expanded={showDetails}
            data-testid="showRequestDetailsButton"
            type="button"
            unstyled
            onClick={handleToggleDetails}
          >
            <FontAwesomeIcon icon={showDetailsChevron} /> {showDetailsText}
          </Button>
        )}
      </div>
      {showDetails && getSectionMarkup(sectionInfo)}
    </section>
  );
}
