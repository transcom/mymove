import React from 'react';
import classnames from 'classnames';
import moment from 'moment';
import { Tag, Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentRequestCard.module.scss';

import { PaymentRequestShape } from 'types/index';
import { formatDateFromIso, formatCents, toDollarString } from 'shared/formatters';

const paymentRequestStatusLabel = (status) => {
  switch (status) {
    case 'PENDING':
      return 'Needs Review';
    case ('REVIEWED', 'SENT_TO_GEX', 'RECEIVED_BY_GEX'):
      return 'Reviewed';
    case 'PAID':
      return 'PAID';
    default:
      return status;
  }
};

const PaymentRequestCard = ({ paymentRequest }) => {
  let requestedAmount = 0;
  let approvedAmount = 0;
  let rejectedAmount = 0;

  if (paymentRequest.serviceItems) {
    paymentRequest.serviceItems.forEach((item) => {
      requestedAmount += item.priceCents;

      if (item.status === 'APPROVED') {
        approvedAmount += item.priceCents;
      } else if (item.status === 'DENIED') {
        rejectedAmount += item.priceCents;
      }
    });
  }

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
                  </div>
                </div>
              )}
              {rejectedAmount > 0 && (
                <div className={styles.amountRejected}>
                  <FontAwesomeIcon icon="times" />
                  <div>
                    <h2>{toDollarString(formatCents(rejectedAmount))}</h2>
                    <span>Rejected</span>
                  </div>
                </div>
              )}
            </>
          )}
          {paymentRequest.status === 'PENDING' && (
            <div className={styles.reviewButton}>
              <button type="button">Review service items</button>
            </div>
          )}
        </div>
        <div className={styles.footer}>
          <dl>
            <dt>Contract Number:</dt>
            <dd />
            <dt>TAC/MDC:</dt>
            <dd />
            <dt>SAC/SDN:</dt>
            <dd />
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
            <Button type="button" unstyled>
              <FontAwesomeIcon icon="chevron-down" /> Show request details
            </Button>
          </div>
        </div>
      </div>
      <div className={styles.drawer} />
    </div>
  );
};

PaymentRequestCard.propTypes = {
  paymentRequest: PaymentRequestShape.isRequired,
};

export default PaymentRequestCard;
