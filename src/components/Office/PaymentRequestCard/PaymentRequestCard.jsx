import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { arrayOf, oneOf, shape, bool, node, string, func } from 'prop-types';
import classnames from 'classnames';
import moment from 'moment';
import { Button, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentRequestCard.module.scss';

import { PaymentRequestShape } from 'types';
import { LOA_TYPE, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { toDollarString, formatDateFromIso, formatCents } from 'utils/formatters';
import PaymentRequestDetails from 'components/Office/PaymentRequestDetails/PaymentRequestDetails';
import ConnectedAcountingCodesModal from 'components/Office/AccountingCodesModal/AccountingCodesModal';
import { groupByShipment } from 'utils/serviceItems';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const paymentRequestStatusLabel = (status) => {
  switch (status) {
    case PAYMENT_REQUEST_STATUS.PENDING:
      return 'Needs review';
    case PAYMENT_REQUEST_STATUS.REVIEWED:
    case PAYMENT_REQUEST_STATUS.SENT_TO_GEX:
    case PAYMENT_REQUEST_STATUS.RECEIVED_BY_GEX:
      return 'Reviewed';
    case PAYMENT_REQUEST_STATUS.REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED:
      return 'Rejected';
    case PAYMENT_REQUEST_STATUS.PAID:
      return 'Paid';
    default:
      return status;
  }
};

const PaymentRequestCard = ({
  paymentRequest,
  shipmentsInfo,
  hasBillableWeightIssues,
  onEditAccountingCodes,
  isMoveLocked,
}) => {
  const navigate = useNavigate();
  // show details by default if in pending/needs review
  const defaultShowDetails = paymentRequest.status === PAYMENT_REQUEST_STATUS.PENDING;
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
  let sortedShipments = [];

  if (paymentRequest.serviceItems) {
    sortedShipments = groupByShipment(paymentRequest.serviceItems);

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
      navigate(paymentRequest.id);
    };
  }

  const uploads = paymentRequest.proofOfServiceDocs
    ? paymentRequest.proofOfServiceDocs.flatMap((docs) => docs.uploads.flatMap((primeUploads) => primeUploads))
    : [];

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

  const showViewDocuments = uploads.length > 0 ? ViewDocuments : <span>No documents provided</span>;

  const tacs = { HHG: tac, NTS: ntsTac };
  const sacs = { HHG: sac, NTS: ntsSac };

  const onEditCodesClick = () => {
    navigate(`/moves/${locator}/orders`);
  };

  const renderReviewServiceItemsBtnForTOO = () => {
    return (
      <Restricted to={permissionTypes.readPaymentServiceItemStatus}>
        <div className={styles.reviewButton}>
          <Button style={{ maxWidth: '225px' }} onClick={handleClick} disabled data-testid="reviewBtn">
            <FontAwesomeIcon icon="copy" className={`${styles['docs-icon']} fas fa-copy`} />
            Review service items
          </Button>
          {hasBillableWeightIssues && (
            <span className={styles.errorText} data-testid="errorTxt">
              Resolve billable weight before reviewing service items.
            </span>
          )}
        </div>
      </Restricted>
    );
  };

  // This defaults to the TIO view but if they don't have permission it tries the TOO view
  const renderReviewServiceItemsBtnForTIOandTOO = () => {
    return (
      <Restricted to={permissionTypes.updatePaymentServiceItemStatus} fallback={renderReviewServiceItemsBtnForTOO()}>
        <div className={styles.reviewButton}>
          <Button
            style={{ maxWidth: '225px' }}
            onClick={handleClick}
            disabled={hasBillableWeightIssues || isMoveLocked}
            data-testid="reviewBtn"
          >
            <FontAwesomeIcon icon="copy" className={`${styles['docs-icon']} fas fa-copy`} />
            Review service items
          </Button>
          {hasBillableWeightIssues && (
            <span className={styles.errorText} data-testid="errorTxt">
              Resolve billable weight before reviewing service items.
            </span>
          )}
        </div>
      </Restricted>
    );
  };

  return (
    <div className={classnames(styles.PaymentRequestCard, 'container')}>
      <div className={styles.summary}>
        <div className={styles.header}>
          <h6>Payment Request {paymentRequest.paymentRequestNumber}</h6>
          <Tag
            className={classnames({
              pending: paymentRequest.status === PAYMENT_REQUEST_STATUS.PENDING,
              reviewed:
                paymentRequest.status !== PAYMENT_REQUEST_STATUS.PENDING &&
                paymentRequest.status !== PAYMENT_REQUEST_STATUS.PAID,
              paid: paymentRequest.status === PAYMENT_REQUEST_STATUS.PAID,
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
          {paymentRequest.status === PAYMENT_REQUEST_STATUS.PENDING ? (
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
          {paymentRequest.status === PAYMENT_REQUEST_STATUS.PENDING && renderReviewServiceItemsBtnForTIOandTOO()}
        </div>
        <div className={styles.footer}>
          <dl>
            <dt>Contract number:</dt>
            <dd>{contractNumber}</dd>
          </dl>
          {!isMoveLocked &&
            (paymentRequest.status === PAYMENT_REQUEST_STATUS.PENDING ? (
              <Link to="../orders" state={{ from: 'paymentRequestDetails' }}>
                View orders
              </Link>
            ) : (
              showViewDocuments
            ))}
          <div className={styles.toggleDrawer}>
            {showRequestDetailsButton && (
              <Button
                aria-expanded={showDetails}
                data-testid="showRequestDetailsButton"
                type="button"
                unstyled
                onClick={handleToggleDetails}
                disabled={isMoveLocked}
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

export default PaymentRequestCard;
