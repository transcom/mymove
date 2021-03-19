import React, { useState, Fragment } from 'react';
import PropTypes, { arrayOf, shape } from 'prop-types';
import classnames from 'classnames';
import moment from 'moment';
import { Button, Tag } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentRequestCard.module.scss';

import { HistoryShape } from 'types/router';
import { PaymentRequestShape } from 'types';
import { formatDateFromIso, formatCents, toDollarString } from 'shared/formatters';
import PaymentRequestDetails from 'components/Office/PaymentRequestDetails/PaymentRequestDetails';
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

const PaymentRequestCard = ({ paymentRequest, shipmentAddresses, history }) => {
  const sortedShipments = groupByShipment(paymentRequest.serviceItems);

  // show details by default if in pending/needs review
  const defaultShowDetails = paymentRequest.status === 'PENDING';
  // only show button in reviewed/paid
  const showRequestDetailsButton = !defaultShowDetails;
  // state to toggle between showing details or not
  const [showDetails, setShowDetails] = useState(defaultShowDetails);

  let handleClick = () => {};
  let requestedAmount = 0;
  let approvedAmount = 0;
  let rejectedAmount = 0;

  const { sac, tac } = paymentRequest.moveTaskOrder.orders;
  const { contractNumber } = paymentRequest.moveTaskOrder.contractor;

  if (paymentRequest.serviceItems) {
    paymentRequest.serviceItems.forEach((item) => {
      requestedAmount += item.priceCents;

      if (item.status === 'APPROVED') {
        approvedAmount += item.priceCents;
      } else if (item.status === 'DENIED') {
        rejectedAmount += item.priceCents;
      }
    });

    handleClick = () => {
      history.push(`payment-requests/${paymentRequest.id}`);
    };
  }

  const showDetailsChevron = showDetails ? 'chevron-up' : 'chevron-down';
  const showDetailsText = showDetails ? 'Hide request details' : 'Show request details';
  const handleToggleDetails = () => setShowDetails((prevState) => !prevState);

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
              <Button onClick={handleClick}>
                <FontAwesomeIcon icon="copy" className={`${styles['docs-icon']} fas fa-copy`} />
                Review service items
              </Button>
            </div>
          )}
        </div>
        <div className={styles.footer}>
          <dl>
            <dt>Contract number:</dt>
            <dd>{contractNumber}</dd>
            <dt>TAC/MDC:</dt>
            <dd>{tac}</dd>
            <dt>SAC/SDN:</dt>
            <dd>{sac}</dd>
          </dl>
          {paymentRequest.status === 'PENDING' ? (
            <a href="orders">View orders</a>
          ) : (
            <a href={`payment-requests/${paymentRequest.id}`}>
              <FontAwesomeIcon icon="copy" />
              View documents
            </a>
          )}
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
      {showDetails && (
        <div data-testid="toggleDrawer" className={styles.drawer}>
          {sortedShipments.map((serviceItems) => {
            let shipmentAddress = '';

            // The service items are grouped by shipment so we only need to check the first value
            const serviceItemShipmentID = serviceItems[0]?.mtoShipmentID;
            if (serviceItemShipmentID) {
              shipmentAddress = shipmentAddresses.find((address) => address.mtoShipmentID === serviceItemShipmentID)
                ?.shipmentAddress;
            }

            return (
              <PaymentRequestDetails
                key={serviceItemShipmentID || 'basicServiceItems'}
                className={styles.paymentRequestDetails}
                serviceItems={serviceItems}
                shipmentAddress={shipmentAddress}
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
  shipmentAddresses: arrayOf(shape({ mtoShipmentId: PropTypes.string, shipmentAddress: PropTypes.string })).isRequired,
};

export default withRouter(PaymentRequestCard);
