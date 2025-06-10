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
import { nonWeightReliantServiceItems } from 'content/serviceItems';
import { toDollarString, formatDateFromIso, formatCents, formatDollarFromMillicents } from 'utils/formatters';
import PaymentRequestDetails from 'components/Office/PaymentRequestDetails/PaymentRequestDetails';
import ConnectedAcountingCodesModal from 'components/Office/AccountingCodesModal/AccountingCodesModal';
import { groupByShipment } from 'utils/serviceItems';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { formatDateWithUTC } from 'shared/dates';

const paymentRequestStatusLabel = (status) => {
  switch (status) {
    case PAYMENT_REQUEST_STATUS.PENDING:
      return 'Needs review';
    case PAYMENT_REQUEST_STATUS.SENT_TO_GEX:
      return 'Sent to GEX';
    case PAYMENT_REQUEST_STATUS.REVIEWED:
      return 'Reviewed';
    case PAYMENT_REQUEST_STATUS.TPPS_RECEIVED:
      return 'TPPS Received';
    case PAYMENT_REQUEST_STATUS.REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED:
      return 'Rejected';
    case PAYMENT_REQUEST_STATUS.PAID:
      return 'TPPS Paid';
    case PAYMENT_REQUEST_STATUS.EDI_ERROR:
      return 'EDI Error';
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

  // do not show error details by default
  const defaultShowErrorDetails = false;
  // only show button in reviewed/paid
  const showErrorDetailsButton = !defaultShowErrorDetails;
  // state to toggle between showing details or not
  const [showErrorDetails, setShowErrorDetails] = useState(defaultShowErrorDetails);

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

  const showErrorDetailsChevron = showErrorDetails ? 'chevron-up' : 'chevron-down';
  const showErrorDetailsText = showErrorDetails ? 'Hide EDI error details' : 'Show EDI error details';
  const handleToggleErrorDetails = () => setShowErrorDetails((prevState) => !prevState);
  const {
    ediErrorCode,
    ediErrorDescription,
    ediErrorType,
    tppsInvoiceAmountPaidTotalMillicents,
    tppsInvoiceSellerPaidDate,
  } = paymentRequest;
  const ediErrorsExistForPaymentRequest = ediErrorCode || ediErrorDescription || ediErrorType;
  const tppsDataExistsForPaymentRequest = tppsInvoiceAmountPaidTotalMillicents !== undefined;
  const showViewDocuments = uploads.length > 0 ? ViewDocuments : <span>No documents provided</span>;

  const tacs = { HHG: tac, NTS: ntsTac };
  const sacs = { HHG: sac, NTS: ntsSac };

  const onEditCodesClick = () => {
    navigate(`/moves/${locator}/orders`);
  };

  const nonWeightRelatedServiceItemsOnly = () => {
    return paymentRequest.serviceItems.every((serviceItem) =>
      Object.prototype.hasOwnProperty.call(nonWeightReliantServiceItems, serviceItem.mtoServiceItemCode),
    );
  };

  const renderReviewServiceItemsBtnForTOO = () => {
    return (
      <Restricted to={permissionTypes.readPaymentServiceItemStatus}>
        <div className={styles.reviewButton}>
          <Button style={{ maxWidth: '225px' }} onClick={handleClick} disabled data-testid="reviewBtn">
            <FontAwesomeIcon icon="copy" className={`${styles['docs-icon']} fas fa-copy`} />
            Review service items
          </Button>
          {hasBillableWeightIssues && !nonWeightRelatedServiceItemsOnly() && (
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
            disabled={(hasBillableWeightIssues && !nonWeightRelatedServiceItemsOnly()) || isMoveLocked}
            data-testid="reviewBtn"
          >
            <FontAwesomeIcon icon="copy" className={`${styles['docs-icon']} fas fa-copy`} />
            Review service items
          </Button>
          {hasBillableWeightIssues && !nonWeightRelatedServiceItemsOnly() && (
            <span className={styles.errorText} data-testid="errorTxt">
              Resolve billable weight before reviewing service items.
            </span>
          )}
        </div>
      </Restricted>
    );
  };

  const renderEDIErrorDetails = () => {
    return (
      <div
        className={
          showErrorDetailsChevron === 'chevron-up' ? styles.ediErrorDetailsExpand : styles.ediErrorDetailsCondensed
        }
      >
        <div className={styles.summary}>
          <div className={styles.footer}>
            <dl>
              <dt>EDI error details:</dt>
            </dl>
            <div className={styles.toggleDrawer}>
              {showErrorDetailsButton && (
                <Button
                  aria-expanded={showErrorDetails}
                  data-testid="showErrorDetailsButton"
                  type="button"
                  unstyled
                  onClick={handleToggleErrorDetails}
                  disabled={isMoveLocked}
                >
                  <FontAwesomeIcon icon={showErrorDetailsChevron} /> {showErrorDetailsText}
                </Button>
              )}
            </div>
          </div>
          {showErrorDetails && (
            <div data-testid="toggleDrawer" className={styles.drawer}>
              <table className="table--stacked">
                <colgroup>
                  <col style={{ width: '20%' }} />
                  <col style={{ width: '20%' }} />
                  <col style={{ width: '60%' }} />
                </colgroup>
                <thead>
                  <tr>
                    <th>EDI Type</th>
                    <th className="align-left">Error Code</th>
                    <th className="align-left">Error Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td data-testid="paymentRequestEDIErrorType">
                      {ediErrorType && <div data-testid="paymentRequestEDIErrorTypeText">{ediErrorType}</div>}
                    </td>
                    <td data-testid="paymentRequestEDIErrorCode" align="top">
                      {ediErrorCode && <div data-testid="paymentRequestEDIErrorCodeText">{ediErrorCode}</div>}
                    </td>
                    <td data-testid="paymentRequestEDIErrorDescription" align="top">
                      {ediErrorDescription && (
                        <div data-testid="paymentRequestEDIErrorDescriptionText">{ediErrorDescription}</div>
                      )}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    );
  };

  const renderApprovedRejectedPaymentRequestDetails = () => {
    if (approvedAmount > 0 || rejectedAmount > 0) {
      return (
        <div data-testid="tppsPaidDetails">
          {approvedAmount > 0 && (
            <div className={styles.amountAccepted} data-testid="milMoveAcceptedDetailsDollarAmountTotal">
              <FontAwesomeIcon icon="check" />
              <div>
                <h2>{toDollarString(formatCents(approvedAmount))}</h2>
                <span>Accepted</span>
                <span> on {formatDateFromIso(paymentRequest.reviewedAt, 'DD MMM YYYY')}</span>
              </div>
            </div>
          )}
          {rejectedAmount > 0 && (
            <div className={styles.amountRejected} data-testid="milMoveRejectedDetailsDollarAmountTotal">
              <FontAwesomeIcon icon="times" />
              <div>
                <h2>{toDollarString(formatCents(rejectedAmount))}</h2>
                <span>Rejected</span>
                <span> on {formatDateFromIso(paymentRequest.reviewedAt, 'DD MMM YYYY')}</span>
              </div>
            </div>
          )}
        </div>
      );
    }
    return null;
  };

  const renderPaymentRequestDetailsForStatus = (paymentRequestStatus) => {
    if (
      (paymentRequestStatus === PAYMENT_REQUEST_STATUS.PAID ||
        paymentRequestStatus === PAYMENT_REQUEST_STATUS.EDI_ERROR) &&
      tppsInvoiceSellerPaidDate
    ) {
      return (
        <div data-testid="tppsPaidDetails">
          {tppsInvoiceAmountPaidTotalMillicents > 0 && (
            <div className={styles.amountAccepted}>
              <FontAwesomeIcon icon="check" />
              <div data-testid="tppsPaidDetailsDollarAmountTotal">
                <h2>{toDollarString(formatDollarFromMillicents(tppsInvoiceAmountPaidTotalMillicents))}</h2>
                <span>TPPS Paid</span>
                <span> on {formatDateWithUTC(tppsInvoiceSellerPaidDate, 'DD MMM YYYY')}</span>
              </div>
            </div>
          )}
        </div>
      );
    }
    if (
      (paymentRequestStatus === PAYMENT_REQUEST_STATUS.TPPS_RECEIVED ||
        paymentRequestStatus === PAYMENT_REQUEST_STATUS.EDI_ERROR) &&
      paymentRequest.receivedByGexAt
    ) {
      return (
        <div>
          {paymentRequest.receivedByGexAt && (
            <div className={styles.amountAccepted}>
              <FontAwesomeIcon icon="check" />
              <div data-testid="tppsReceivedDetailsDollarAmountTotal">
                <h2>{toDollarString(formatCents(approvedAmount))}</h2>
                <span>TPPS Received</span>
                <span> on {formatDateFromIso(paymentRequest.receivedByGexAt, 'DD MMM YYYY')}</span>
              </div>
            </div>
          )}
        </div>
      );
    }

    if (
      paymentRequestStatus === PAYMENT_REQUEST_STATUS.SENT_TO_GEX ||
      (paymentRequestStatus === PAYMENT_REQUEST_STATUS.EDI_ERROR && approvedAmount > 0)
    ) {
      return (
        <div className={styles.amountAccepted} data-testid="sentToGexDetails">
          <FontAwesomeIcon icon="check" />
          <div data-testid="sentToGexDetailsDollarAmountTotal">
            <h2>{toDollarString(formatCents(approvedAmount))}</h2>
            <span>Sent to GEX </span>
            <span data-testid="sentToGexDate">
              on {paymentRequest?.sentToGexAt ? formatDateFromIso(paymentRequest.sentToGexAt, 'DD MMM YYYY') : '-'}
            </span>
          </div>
        </div>
      );
    }
    if (
      (paymentRequestStatus === PAYMENT_REQUEST_STATUS.PENDING ||
        paymentRequestStatus === PAYMENT_REQUEST_STATUS.EDI_ERROR) &&
      requestedAmount > 0
    ) {
      return (
        <div className={styles.amountRequested}>
          <h2>{toDollarString(formatCents(requestedAmount))}</h2>
          <span>Requested</span>
        </div>
      );
    }
    return null;
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
          <div>
            {paymentRequest.status && renderApprovedRejectedPaymentRequestDetails(paymentRequest)}
            {paymentRequest.status && renderPaymentRequestDetailsForStatus(paymentRequest.status)}
          </div>
          {paymentRequest.status === PAYMENT_REQUEST_STATUS.PENDING && renderReviewServiceItemsBtnForTIOandTOO()}
        </div>
        {ediErrorsExistForPaymentRequest && renderEDIErrorDetails()}
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
                tppsDataExists={tppsDataExistsForPaymentRequest}
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
