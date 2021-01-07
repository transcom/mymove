import React from 'react';
import classnames from 'classnames';
import moment from 'moment';
import { Button, Tag } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { HistoryShape } from '../../../types/router';

import styles from './PaymentRequestCard.module.scss';

import { formatDateWithoutDash, formatDateFromIso, formatCents, toDollarString } from 'shared/formatters';
import { PaymentRequestShape } from 'types/index';

const paymentRequestStatusLabel = (status) => {
  switch (status) {
    case 'PENDING':
      return 'Needs Review';
    case 'REVIEWED':
    case 'SENT_TO_GEX':
    case 'RECEIVED_BY_GEX':
      return 'Reviewed';
    case 'PAID':
      return 'Paid';
    default:
      return status;
  }
};

const PaymentRequestCard = ({ paymentRequest, history }) => {
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
                    <span> on {formatDateWithoutDash(paymentRequest.reviewedAt)}</span>
                  </div>
                </div>
              )}
              {rejectedAmount > 0 && (
                <div className={styles.amountRejected}>
                  <FontAwesomeIcon icon="times" />
                  <div>
                    <h2>{toDollarString(formatCents(rejectedAmount))}</h2>
                    <span>Rejected</span>
                    <span> on {formatDateWithoutDash(paymentRequest.reviewedAt)}</span>
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
            <dt>Contract Number:</dt>
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
  history: HistoryShape.isRequired,
  paymentRequest: PaymentRequestShape.isRequired,
};

export default withRouter(PaymentRequestCard);
