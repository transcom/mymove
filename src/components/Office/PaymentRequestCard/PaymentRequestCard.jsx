import React, { useState } from 'react';
import { arrayOf, oneOf, shape, bool, node, string, func } from 'prop-types';
import classnames from 'classnames';
import moment from 'moment';
import { Button, Tag } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentRequestCard.module.scss';

import { HistoryShape } from 'types/router';
import { PaymentRequestShape } from 'types';
import { LOA_TYPE, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { formatCents } from 'shared/formatters';
import { toDollarString, formatDateFromIso } from 'utils/formatters';
import PaymentRequestDetails from 'components/Office/PaymentRequestDetails/PaymentRequestDetails';
import ConnectedAcountingCodesModal from 'components/Office/AccountingCodesModal/AccountingCodesModal';
import { groupByShipment } from 'utils/serviceItems';

const paymentRequestStatusLabel = (status) => {
  switch (status) {
    case 'PENDING':
      return 'Needs review';
    case 'REVIEWED':
    case 'SENT_TO_GEX':
    case 'RECEIVED_BY_GEX':
      return 'Reviewed';
    case 'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED':
      return 'Rejected';
    case 'PAID':
      return 'Paid';
    default:
      return status;
  }
};

const PaymentRequestCard = ({
  paymentRequest,
  shipmentsInfo,
  history,
  hasBillableWeightIssues,
  onEditAccountingCodes,
}) => {
  const sortedShipments = groupByShipment(paymentRequest.serviceItems);

  // show details by default if in pending/needs review
  const defaultShowDetails = paymentRequest.status === 'PENDING';
  // only show button in reviewed/paid
  const showRequestDetailsButton = !defaultShowDetails;
  // state to toggle between showing details or not
  const [showDetails, setShowDetails] = useState(defaultShowDetails);

  // show/hide AccountingCodesModal
  const [showModal, setShowModal] = useState(false);
  const [modalShipment, setModalShipment] = useState({});

  const handleModalSave = (values) => {
    onEditAccountingCodes(modalShipment.mtoShipmentID, {
      tacType: values.tacType,
      sacType: values.sacType,
    });

    setShowModal(false);
    setModalShipment({});
  };

  const handleModalCancel = () => {
    setShowModal(false);
    setModalShipment({});
  };

  const onEditClick = (shipment = {}) => {
    setShowModal(true);
    setModalShipment(shipment);
  };

  let handleClick = () => {};
  let requestedAmount = 0;
  let approvedAmount = 0;
  let rejectedAmount = 0;

  const { locator } = paymentRequest.moveTaskOrder;
  const { sac, tac, ntsTac, ntsSac } = paymentRequest.moveTaskOrder.orders;
  const { contractNumber } = paymentRequest.moveTaskOrder.contractor;

  if (paymentRequest.serviceItems) {
    paymentRequest.serviceItems.forEach((item) => {
      if (item.priceCents != null) {
        requestedAmount += item.priceCents;

        if (item.status === 'APPROVED') {
          approvedAmount += item.priceCents;
        } else if (item.status === 'DENIED') {
          rejectedAmount += item.priceCents;
        }
      }
    });

    handleClick = () => {
      history.push(`payment-requests/${paymentRequest.id}`);
    };
  }

  const showDetailsChevron = showDetails ? 'chevron-up' : 'chevron-down';
  const showDetailsText = showDetails ? 'Hide request details' : 'Show request details';
  const handleToggleDetails = () => setShowDetails((prevState) => !prevState);
  const ViewDocuments =
    paymentRequest.status !== PAYMENT_REQUEST_STATUS.DEPRECATED ? (
      <a href={`payment-requests/${paymentRequest.id}`}>
        <FontAwesomeIcon icon="copy" />
        View documents
      </a>
    ) : null;

  const tacs = { HHG: tac, NTS: ntsTac };
  const sacs = { HHG: sac, NTS: ntsSac };

  const onEditCodesClick = () => {
    history.push(`/moves/${locator}/orders`);
  };

  return (
    <div className={classnames(styles.PaymentRequestCard, 'container')}>
      <div className={styles.summary}>
        <div className={styles.header}>
          <h6>Payment Request {paymentRequest.paymentRequestNumber}</h6>
          <Tag
            className={classnames({
              pending: paymentRequest.status === 'PENDING',
              reviewed: paymentRequest.status !== 'PENDING' && paymentRequest.status !== 'PAID',
              paid: paymentRequest.status === 'PAID',
            })}
          >
            {paymentRequestStatusLabel(paymentRequest.status)}
          </Tag>
          <span className={styles.dateSubmitted}>
            Submitted {moment(paymentRequest.createdAt).fromNow()} on{' '}
            {formatDateFromIso(paymentRequest.createdAt, 'DD MMM YYYY')}
          </span>
        </div>
        <div className={styles.totalReviewed}>
          {paymentRequest.status === 'PENDING' ? (
            <div className={styles.amountRequested}>
              <h2>{toDollarString(formatCents(requestedAmount))}</h2>
              <span>Requested</span>
            </div>
          ) : (
            <>
              {approvedAmount > 0 && (
                <div className={styles.amountAccepted}>
                  <FontAwesomeIcon icon="check" />
                  <div>
                    <h2>{toDollarString(formatCents(approvedAmount))}</h2>
                    <span>Accepted</span>
                    <span> on {formatDateFromIso(paymentRequest.reviewedAt, 'DD MMM YYYY')}</span>
                  </div>
                </div>
              )}
              {rejectedAmount > 0 && (
                <div className={styles.amountRejected}>
                  <FontAwesomeIcon icon="times" />
                  <div>
                    <h2>{toDollarString(formatCents(rejectedAmount))}</h2>
                    <span>Rejected</span>
                    <span> on {formatDateFromIso(paymentRequest.reviewedAt, 'DD MMM YYYY')}</span>
                  </div>
                </div>
              )}
            </>
          )}
          {paymentRequest.status === 'PENDING' && (
            <div className={styles.reviewButton}>
              <Button
                style={{ maxWidth: '225px' }}
                onClick={handleClick}
                disabled={hasBillableWeightIssues}
                test-dataid="reviewBtn"
              >
                <FontAwesomeIcon icon="copy" className={`${styles['docs-icon']} fas fa-copy`} />
                Review service items
              </Button>
              {hasBillableWeightIssues && (
                <span className={styles.errorText} test-dataid="errorTxt">
                  Resolve billable weight before reviewing service items.
                </span>
              )}
            </div>
          )}
        </div>
        <div className={styles.footer}>
          <dl>
            <dt>Contract number:</dt>
            <dd>{contractNumber}</dd>
          </dl>
          {paymentRequest.status === 'PENDING' ? <a href="orders">View orders</a> : ViewDocuments}
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
        </div>
      </div>
      <ConnectedAcountingCodesModal
        isOpen={showModal}
        shipmentType={modalShipment.shipmentType}
        TACs={tacs}
        SACs={sacs}
        onClose={handleModalCancel}
        onSubmit={handleModalSave}
        sacType={modalShipment.sacType}
        tacType={modalShipment.tacType}
        onEditCodesClick={onEditCodesClick}
      />
      {showDetails && (
        <div data-testid="toggleDrawer" className={styles.drawer}>
          {sortedShipments.map((serviceItems) => {
            let selectedShipment = {};

            // The service items are grouped by shipment so we only need to check the first value
            const serviceItemShipmentID = serviceItems[0]?.mtoShipmentID;
            if (serviceItemShipmentID && shipmentsInfo) {
              selectedShipment = shipmentsInfo.find((shipment) => shipment.mtoShipmentID === serviceItemShipmentID);
            }

            return (
              <PaymentRequestDetails
                key={serviceItemShipmentID || 'basicServiceItems'}
                className={styles.paymentRequestDetails}
                serviceItems={serviceItems}
                shipment={selectedShipment}
                paymentRequestStatus={paymentRequest.status}
                tacs={tacs}
                sacs={sacs}
                onEditClick={onEditClick}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};

PaymentRequestCard.propTypes = {
  history: HistoryShape.isRequired,
  paymentRequest: PaymentRequestShape.isRequired,
  hasBillableWeightIssues: bool.isRequired,
  shipmentsInfo: arrayOf(
    shape({
      mtoShipmentID: string,
      shipmentAddress: node,
      departureDate: string,
      shipmentModificationType: string,
      tacType: oneOf(Object.values(LOA_TYPE)),
      sacType: oneOf(Object.values(LOA_TYPE)),
    }),
  ),
  onEditAccountingCodes: func,
};

PaymentRequestCard.defaultProps = {
  shipmentsInfo: [],
  onEditAccountingCodes: () => {},
};

export default withRouter(PaymentRequestCard);
