import React from 'react';
import { arrayOf, bool, func, number } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './PPMSummaryList.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { ppmShipmentStatuses } from 'constants/shipments';
import { ShipmentShape } from 'types/shipment';
import { formatCustomerDate } from 'utils/formatters';

const submittedContent = (
  <>
    <p>After a counselor approves your PPM, you will be able to:</p>
    <ul>
      <li>Download paperwork for an advance, if you requested one</li>
      <li>Upload PPM documents and start the payment request process</li>
    </ul>
  </>
);

const approvedContent = (approvedAt) => {
  return (
    <>
      <div className={styles.dateSummary}>
        <p>{`PPM approved: ${formatCustomerDate(approvedAt)}.`}</p>
      </div>
      <div>
        <p>
          When you are ready to request payment for this PPM, select Upload PPM Documents to add paperwork, calculate
          your incentive, and create a payment request packet.
        </p>
      </div>
    </>
  );
};

const paymentSubmitted = (approvedAt, submittedAt) => {
  return (
    <>
      <div className={styles.dateSummary}>
        <p>{`PPM approved: ${formatCustomerDate(approvedAt)}`}</p>
        <p>{`PPM documentation submitted: ${formatCustomerDate(submittedAt)}`}</p>
      </div>
      <div>
        <p>
          A counselor will review your documentation. When it&apos;s verified, you can visit MilMove to download the
          incentive packet that you&apos;ll need to give to Finance.
        </p>
      </div>
    </>
  );
};

const paymentReviewed = (approvedAt, submittedAt, reviewedAt) => {
  return (
    <>
      <div className={styles.dateSummary}>
        <p>{`PPM approved: ${formatCustomerDate(approvedAt)}`}</p>
        <p>{`PPM documentation submitted: ${formatCustomerDate(submittedAt)}`}</p>
        <p>{`Documentation accepted and verified: ${formatCustomerDate(reviewedAt)}`}</p>
      </div>
      <div>
        <p>
          You can now download your incentive packet and submit it to Finance to request payment. You will also need to
          include a completed DD-1351-2, and any other paperwork required by your service.
        </p>
      </div>
    </>
  );
};

const PPMSummaryStatus = (shipment, orderLabel, onButtonClick) => {
  const {
    ppmShipment: { status, approvedAt, submittedAt, reviewedAt },
  } = shipment;

  let actionButton;
  let content;

  switch (status) {
    case ppmShipmentStatuses.SUBMITTED:
      actionButton = <Button disabled>Upload PPM Documents</Button>;
      content = submittedContent;
      break;
    case ppmShipmentStatuses.WAITING_ON_CUSTOMER:
      actionButton = <Button onClick={onButtonClick}>Upload PPM Documents</Button>;
      content = approvedContent(approvedAt);
      break;
    case ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL:
      actionButton = <Button disabled>Download Incentive Packet</Button>;
      content = paymentSubmitted(approvedAt, submittedAt);
      break;
    case ppmShipmentStatuses.PAYMENT_APPROVED:
      actionButton = <Button onClick={onButtonClick}>Download Incentive Packet</Button>;
      content = paymentReviewed(approvedAt, submittedAt, reviewedAt);
      break;
    default:
  }

  return (
    <SectionWrapper className={styles['ppm-shipment']}>
      <div className={styles['ppm-shipment__heading-section']}>
        <strong>{orderLabel}</strong>
        {actionButton}
      </div>
      <div className={styles['ppm-shipment__content']}>{content}</div>
    </SectionWrapper>
  );
};

const PPMSummaryList = ({ shipments, onUploadClick }) => {
  const { length } = shipments;
  return shipments.map((shipment, i) => {
    return (
      <PPMSummaryListItem
        key={shipment.id}
        shipment={shipment}
        hasMany={length > 1}
        index={i}
        onUploadClick={() => onUploadClick(shipment.id)}
      />
    );
  });
};

const PPMSummaryListItem = ({ shipment, hasMany, index, onUploadClick }) => {
  const orderLabel = hasMany ? `PPM ${index + 1}` : 'PPM';

  return PPMSummaryStatus(shipment, orderLabel, onUploadClick);
};

PPMSummaryList.propTypes = {
  shipments: arrayOf(ShipmentShape).isRequired,
  onUploadClick: func.isRequired,
};

PPMSummaryListItem.propTypes = {
  shipment: ShipmentShape.isRequired,
  index: number.isRequired,
  hasMany: bool.isRequired,
  onUploadClick: func.isRequired,
};

export default PPMSummaryList;
