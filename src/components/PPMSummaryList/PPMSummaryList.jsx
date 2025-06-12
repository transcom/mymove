import React from 'react';
import { arrayOf, bool, func, number } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './PPMSummaryList.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { ppmShipmentStatuses } from 'constants/shipments';
import { ShipmentShape } from 'types/shipment';
import { formatCustomerDate, formatAddressShort } from 'utils/formatters';
import AsyncPacketDownloadLink from 'shared/AsyncPacketDownloadLink/AsyncPacketDownloadLink';
import { downloadPPMPaymentPacket } from 'services/internalApi';
import { isFeedbackAvailable } from 'constants/ppmFeedback';

const toFromAddressDisplay = (pickupAddress, destinationAddress) => {
  return (
    <div className={styles.addressDisplay}>
      <p>
        <span className={styles.bold}>From: </span>
        {formatAddressShort(pickupAddress)}
        <span className={styles.bold}> To: </span>
        {formatAddressShort(destinationAddress)}
      </p>
    </div>
  );
};

const submittedContent = (
  <>
    <p>After a counselor approves your PPM, you will be able to:</p>
    <ul>
      <li>Download paperwork for an advance, if you requested one</li>
      <li>Upload PPM documents and start the payment request process</li>
    </ul>
  </>
);

const approvedContent = (approvedAt, pickupAddress, destinationAddress) => {
  return (
    <>
      {toFromAddressDisplay(pickupAddress, destinationAddress)}
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

const paymentSubmitted = (approvedAt, submittedAt, pickupAddress, destinationAddress) => {
  return (
    <>
      {toFromAddressDisplay(pickupAddress, destinationAddress)}
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

const paymentReviewed = (approvedAt, submittedAt, reviewedAt, pickupAddress, destinationAddress) => {
  return (
    <>
      {toFromAddressDisplay(pickupAddress, destinationAddress)}
      <div className={styles.dateSummary}>
        <p>{`PPM approved: ${formatCustomerDate(approvedAt)}`}</p>
        <p>{`PPM documentation submitted: ${formatCustomerDate(submittedAt)}`}</p>
        <p>{`PPM closeout completed: ${formatCustomerDate(reviewedAt)}`}</p>
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

const PPMSummaryStatus = (shipment, orderLabel, onButtonClick, onDownloadError, onFeedbackClick) => {
  const {
    ppmShipment: { status, approvedAt, submittedAt, reviewedAt, pickupAddress, destinationAddress },
  } = shipment;

  const handleDownloadFailure = () => {
    onDownloadError();
  };

  let actionButtons;
  let content;

  switch (status) {
    case ppmShipmentStatuses.SUBMITTED:
      actionButtons = <Button disabled>Upload PPM Documents</Button>;
      content = submittedContent;
      break;
    case ppmShipmentStatuses.WAITING_ON_CUSTOMER:
      actionButtons = isFeedbackAvailable(shipment?.ppmShipment) ? (
        [
          <div>
            <Button onClick={() => onFeedbackClick()}>View Closeout Feedback</Button>
            <Button onClick={onButtonClick}>Upload PPM Documents</Button>
          </div>,
        ]
      ) : (
        <Button onClick={onButtonClick}>Upload PPM Documents</Button>
      );
      content = approvedContent(approvedAt, pickupAddress, destinationAddress);
      break;
    case ppmShipmentStatuses.NEEDS_CLOSEOUT:
      actionButtons = <Button disabled>Download Payment Packet</Button>;
      content = paymentSubmitted(approvedAt, submittedAt, pickupAddress, destinationAddress);
      break;
    case ppmShipmentStatuses.CLOSEOUT_COMPLETE:
      actionButtons = isFeedbackAvailable(shipment?.ppmShipment) ? (
        [
          <div>
            <Button onClick={() => onFeedbackClick()}>View Closeout Feedback</Button>
            <AsyncPacketDownloadLink
              id={shipment?.ppmShipment?.id}
              label="Download Payment Packet"
              asyncRetrieval={downloadPPMPaymentPacket}
              onFailure={handleDownloadFailure}
              className="styles.btn"
              loadingMessage="Downloading Payment Packet (PDF)..."
            />
          </div>,
        ]
      ) : (
        <AsyncPacketDownloadLink
          id={shipment?.ppmShipment?.id}
          label="Download Payment Packet"
          asyncRetrieval={downloadPPMPaymentPacket}
          onFailure={handleDownloadFailure}
          className="styles.btn"
          loadingMessage="Downloading Payment Packet (PDF)..."
        />
      );

      content = paymentReviewed(approvedAt, submittedAt, reviewedAt, pickupAddress, destinationAddress);
      break;
    default:
  }
  return (
    <SectionWrapper className={styles['ppm-shipment']}>
      <div className={styles['ppm-shipment__heading-section']}>
        <strong>{orderLabel}</strong>
        {actionButtons}
      </div>
      <div className={styles['ppm-shipment__content']}>{content}</div>
    </SectionWrapper>
  );
};

const PPMSummaryList = ({ shipments, onUploadClick, onDownloadError, onFeedbackClick }) => {
  const { length } = shipments;
  return shipments.map((shipment, i) => {
    return (
      <PPMSummaryListItem
        key={shipment.id}
        shipment={shipment}
        hasMany={length > 1}
        index={i}
        onUploadClick={() => onUploadClick(shipment.id)}
        onDownloadError={onDownloadError}
        onFeedbackClick={() => onFeedbackClick(shipment.id)}
      />
    );
  });
};

const PPMSummaryListItem = ({ shipment, hasMany, index, onUploadClick, onDownloadError, onFeedbackClick }) => {
  const orderLabel = hasMany ? `PPM ${index + 1}` : 'PPM';
  return PPMSummaryStatus(shipment, orderLabel, onUploadClick, onDownloadError, onFeedbackClick);
};

PPMSummaryList.propTypes = {
  shipments: arrayOf(ShipmentShape).isRequired,
  onUploadClick: func.isRequired,
  onDownloadError: func.isRequired,
};

PPMSummaryListItem.propTypes = {
  shipment: ShipmentShape.isRequired,
  index: number.isRequired,
  hasMany: bool.isRequired,
  onUploadClick: func.isRequired,
  onDownloadError: func.isRequired,
};

export default PPMSummaryList;
