import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Label, Button, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './HeaderSection.module.scss';

import { formatDate, formatCents, formatWeight } from 'utils/formatters';

export const sectionTypes = {
  incentives: 'incentives',
  shipmentInfo: 'shipmentInfo',
  gcc: 'gcc',
};

const getSectionTitle = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.incentives:
      return `Incentives/Costs`;
    case sectionTypes.shipmentInfo:
      return `Shipment Info`;
    case sectionTypes.gcc:
      return `GCC Factors`;
    default:
      return <Alert>`Error getting section title!`</Alert>;
  }
};

// Returns the markup needed for a specific section
const getSectionMarkup = (sectionInfo) => {
  switch (sectionInfo.type) {
    case sectionTypes.shipmentInfo:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Planned Move Start Date</Label>
            <span className={styles.light}>
              {sectionInfo.expectedDepartureDate
                ? formatDate(sectionInfo.expectedDepartureDate, null, 'DD-MMM-YYYY')
                : `TEST DATE`}
            </span>
          </div>
          <div>
            <Label>Actual Move Start Date</Label>
            <span className={styles.light}>{formatDate(sectionInfo.actualMoveDate, null, 'DD-MMM-YYYY')}</span>
          </div>
          <div>
            <Label>Starting ZIP</Label>
            <span className={styles.light}>{sectionInfo.actualPickupPostalCode}</span>
          </div>
          <div>
            <Label>Ending ZIP</Label>
            <span className={styles.light}>{sectionInfo.actualDestinationPostalCode}</span>
          </div>
          <div>
            <Label>Miles</Label>
            <span className={styles.light}>{sectionInfo.miles ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Estimated Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.estimatedWeight)}</span>
          </div>
          <div>
            <Label>Actual Net Weight</Label>
            <span className={styles.light}>{formatWeight(sectionInfo.actualWeight) ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Advance received</Label>
            <span className={styles.light}>
              {sectionInfo.hasReceivedAdvance
                ? `$${formatCents(sectionInfo.advanceAmountReceived)}`
                : 'Not requested/received.'}
            </span>
          </div>
        </div>
      );

    case sectionTypes.incentives:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Gross Incentive</Label>
            {/** TODO: Is estimatedIncentive (sent from ppmShipment in PPMHeaderSummary) actually the correct value? */}
            <span className={styles.light}>${formatCents(sectionInfo.estimatedIncentive)}</span>
          </div>
          <div>
            <Label>Government Constructive Cost (GCC)</Label>
            <span className={styles.light}>${formatCents(sectionInfo.gcc) ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Remaining Reimbursement Owed to Customer</Label>
            <span className={styles.light}>${formatCents(sectionInfo.remainingReimbursement) ?? `TEST VAL`}</span>
          </div>
        </div>
      );

    case sectionTypes.gcc:
      return (
        <div className={classnames(styles.Details)}>
          <div>
            <Label>Linehaul Price</Label>
            <span className={styles.light}>${formatCents(sectionInfo.linehaulPrice) ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Linehaul Fuel Surcharge</Label>
            <span className={styles.light}>${formatCents(sectionInfo.linehaulFuelSurcharge) ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Shorthaul Price</Label>
            <span className={styles.light}>${formatCents(sectionInfo.shorthaulPrice) ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Shorthaul Fuel Surcharge</Label>
            <span className={styles.light}>${formatCents(sectionInfo.shorthaulFuelSurcharge) ?? `TEST VAL`}</span>
          </div>
          <div>
            <Label>Full Pack/Unpack Charge</Label>
            <span className={styles.light}>${formatCents(sectionInfo.fullPackUnpackCharge) ?? `TEST VAL`}</span>
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
