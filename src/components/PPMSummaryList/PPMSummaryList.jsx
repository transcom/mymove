import React from 'react';
import { arrayOf, bool, number, shape, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './PPMSummaryList.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';
import { formatCustomerDate } from 'utils/formatters';

const ppmContent = (canUpload, approvedOn) => {
  return canUpload ? (
    <>
      <p>{`PPM approved: ${formatCustomerDate(approvedOn)}.`}</p>
      <p>
        When you are ready to request payment for this PPM, select Upload PPM Documents to add paperwork, calculate your
        incentive, and create a payment request packet.
      </p>
    </>
  ) : (
    <>
      <p>After a counselor approves your PPM, you will be able to:</p>
      <ul>
        <li>Download paperwork for an advance, if you requested one</li>
        <li>Upload PPM documents and start the payment request process</li>
      </ul>
    </>
  );
};

const PPMSummaryList = ({ shipments }) => {
  const { length } = shipments;
  return shipments.map((shipment, i) => {
    return <PPMSummaryListItem key={shipment.id} shipment={shipment} hasMany={length > 1} index={i} />;
  });
};

const PPMSummaryListItem = ({ shipment, hasMany, index }) => {
  const canUpload =
    shipment.status === shipmentStatuses.APPROVED &&
    shipment?.ppmShipment?.status === ppmShipmentStatuses.WAITING_ON_CUSTOMER;
  return (
    <SectionWrapper className={styles['ppm-shipment']}>
      <div className={styles['heading-section']}>
        <strong>{hasMany ? `PPM ${index + 1}` : 'PPM'}</strong>
        <Button disabled={!canUpload}>Upload PPM Documents</Button>
      </div>
      <div className={styles.content}>{ppmContent(canUpload, shipment?.ppmShipment?.approvedAt)}</div>
    </SectionWrapper>
  );
};

const shipmentShape = shape({
  id: string.isRequired,
});

PPMSummaryList.propTypes = {
  shipments: arrayOf(shipmentShape).isRequired,
};

PPMSummaryListItem.propTypes = {
  shipment: shipmentShape.isRequired,
  index: number.isRequired,
  hasMany: bool.isRequired,
};

export default PPMSummaryList;
