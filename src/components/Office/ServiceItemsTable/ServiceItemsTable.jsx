import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useParams } from 'react-router';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './ServiceItemsTable.module.scss';

import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { ALLOWED_RESUBMISSION_SI_CODES, ALLOWED_SIT_UPDATE_SI_CODES } from 'constants/sitUpdates';
import { formatDateFromIso } from 'utils/formatters';
import ServiceItemDetails from 'components/Office/ServiceItemDetails/ServiceItemDetails';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { selectDateFieldByStatus, selectDatePrefixByStatus } from 'utils/dates';
import { useGHCGetMoveHistory, useMovePaymentRequestsQueries } from 'hooks/queries';
import ToolTip from 'shared/ToolTip/ToolTip';
import { ShipmentShape } from 'types';

// Sorts service items in an order preferred by the customer
// Currently only SIT & shorthaul/linehaul receives special sorting
// this current listed order goes:
// shorthaul & linehaul
// other service items
// origin SIT
// destination SIT
function sortServiceItems(items) {
  // Prioritize service items with codes 'DSH' (shorthaul) and 'DLH' (linehaul) to be at the top of the list
  const haulTypeServiceItemCodes = ['DSH', 'DLH'];
  const haulTypeServiceItems = items.filter((item) => haulTypeServiceItemCodes.includes(item.code));
  const sortedHaulTypeServiceItems = haulTypeServiceItems.sort(
    (a, b) => haulTypeServiceItemCodes.indexOf(a.code) - haulTypeServiceItemCodes.indexOf(b.code),
  );
  // Filter and sort destination SIT. Code index is also the sort order
  const destinationServiceItemCodes = ['DDFSIT', 'DDASIT', 'DDDSIT', 'DDSFSC'];
  const destinationServiceItems = items.filter((item) => destinationServiceItemCodes.includes(item.code));
  const sortedDestinationServiceItems = destinationServiceItems.sort(
    (a, b) => destinationServiceItemCodes.indexOf(a.code) - destinationServiceItemCodes.indexOf(b.code),
  );
  // Filter origin SIT. Code index is also the sort order
  const originServiceItemCodes = ['DOFSIT', 'DOASIT', 'DOPSIT', 'DOSFSC'];
  const originServiceItems = items.filter((item) => originServiceItemCodes.includes(item.code));
  const sortedOriginServiceItems = originServiceItems.sort(
    (a, b) => originServiceItemCodes.indexOf(a.code) - originServiceItemCodes.indexOf(b.code),
  );

  // Filter all service items that are not specifically sorted
  const remainingServiceItems = items.filter(
    (item) =>
      !haulTypeServiceItemCodes.includes(item.code) &&
      !destinationServiceItemCodes.includes(item.code) &&
      !originServiceItemCodes.includes(item.code),
  );

  return [
    ...sortedHaulTypeServiceItems,
    ...remainingServiceItems,
    ...sortedOriginServiceItems,
    ...sortedDestinationServiceItems,
  ];
}

const ServiceItemsTable = ({
  serviceItems,
  statusForTableType,
  handleUpdateMTOServiceItemStatus,
  handleShowRejectionDialog,
  handleShowEditSitAddressModal,
  handleShowEditSitEntryDateModal,
  shipment,
  isMoveLocked,
}) => {
  const getServiceItemDisplayDate = (item) => {
    const prefix = selectDatePrefixByStatus(statusForTableType);
    const date = formatDateFromIso(item[`${selectDateFieldByStatus(statusForTableType)}`], 'DD MMM YYYY');
    return `${prefix}: ${date}`;
  };

  // adding in payment requests to determine edit button status
  const { moveCode } = useParams();
  const { paymentRequests } = useMovePaymentRequestsQueries(moveCode);
  let serviceItemInPaymentRequests;
  if (paymentRequests.some((obj) => 'serviceItems' in obj)) {
    serviceItemInPaymentRequests = paymentRequests.map((obj) => ({
      serviceItems: obj.serviceItems.map((s) => s.mtoServiceItemID),
    }));
  }

  // function iterating through payment requests to find if a service item is in there
  const isServiceItemFoundInPaymentRequests = (id) => {
    return serviceItemInPaymentRequests.some((obj) => {
      if (obj.serviceItems.includes(id)) {
        return true; // Set the result to true when id is found
      }
      return false; // Return false when id is not found
    });
  };

  const getResubmissionStatus = (historyRecordOfServiceItem) => {
    let isResubmitted = false;
    if (historyRecordOfServiceItem) {
      if (
        historyRecordOfServiceItem.action === 'UPDATE' &&
        historyRecordOfServiceItem.oldValues.status === 'REJECTED' &&
        historyRecordOfServiceItem.eventName === 'updateMTOServiceItem' &&
        statusForTableType === SERVICE_ITEM_STATUS.SUBMITTED
      ) {
        isResubmitted = true;
      }
    }
    return isResubmitted;
  };

  const getNewestHistoryDataForServiceItem = (historyDataForMove, serviceItemId) => {
    if (historyDataForMove) {
      let newestHistoryData = historyDataForMove[0];
      for (let i = 0; i < historyDataForMove.length; i += 1) {
        // find the first event in the move history for a given serviceItemId
        if (historyDataForMove[i].objectId === serviceItemId) {
          newestHistoryData = historyDataForMove[i];
          break;
        }
      }
      return newestHistoryData;
    }
    return null;
  };

  function formatKeyStringsForToolTip(key) {
    // replace _ with ' ' and capitalize first letters of each word
    let changedKey = key.replace(/(^|_)./g, (index) => index.toUpperCase().replace('_', ' '));
    if (changedKey.indexOf('Sit') === 0) {
      const replacement = 'SIT';
      // Replace the first three characters with 'SIT'
      changedKey = replacement + changedKey.slice(3);
    }
    if (changedKey.indexOf('Id') === 0) {
      const replacement = 'ID';
      // Replace the first two characters with 'ID'
      changedKey = replacement + changedKey.slice(2);
    }
    return changedKey;
  }

  function generateResubmissionDetailsText(details) {
    if (!details) return '';

    const keys = Object.keys(details.changedValues);

    return keys.map((key) => {
      const formattedKeyString = formatKeyStringsForToolTip(key);
      const newValue = details.changedValues[key];
      const oldValue = details.oldValues[key];

      return (
        <div key={key} className={styles.resubmissionDetails}>
          <div>{formattedKeyString?.toUpperCase()}</div>
          <div>
            <strong>New:</strong> <span className={styles.detailValue}>{newValue?.toString()}</span>
          </div>
          <div>
            <strong>Previous:</strong> <span className={styles.detailValue}>{oldValue?.toString()}</span>
          </div>
        </div>
      );
    });
  }

  const history = useGHCGetMoveHistory({ moveCode });
  const renderToolTipWithOldDataIfResubmission = (serviceItemId) => {
    const historyDataForMove = history.queueResult.data;
    const historyDataForServiceItem = getNewestHistoryDataForServiceItem(historyDataForMove, serviceItemId);
    const isResubmitted = getResubmissionStatus(historyDataForServiceItem);
    let formattedResubmissionDetails = '';
    if (isResubmitted) {
      formattedResubmissionDetails = generateResubmissionDetailsText(historyDataForServiceItem);
    }
    const resubmittedServiceItemValues = {
      isResubmitted,
      formattedResubmissionDetails,
    };
    return resubmittedServiceItemValues;
  };

  const sortedServiceItems = sortServiceItems(serviceItems);
  const tableRows = sortedServiceItems.map((serviceItem) => {
    const { id, code, details, mtoShipmentID, serviceRequestDocuments, ...item } = serviceItem;
    let hasPaymentRequestBeenMade;
    // if there are service items in the payment requests, we want to look to see if the service item is in there
    // if so, we don't want to let the TOO edit the SIT entry date
    if (serviceItemInPaymentRequests && ALLOWED_SIT_UPDATE_SI_CODES.includes(code)) {
      hasPaymentRequestBeenMade = isServiceItemFoundInPaymentRequests(id);
    }
    const resubmittedToolTip = renderToolTipWithOldDataIfResubmission(id);

    // we don't want to display the "Accept" button for a DLH or DSH service item that was rejected by a shorthaul to linehaul change or vice versa
    let rejectedDSHorDLHServiceItem = false;
    if (
      (serviceItem.code === 'DLH' || serviceItem.code === 'DSH') &&
      serviceItem.details.rejectionReason ===
        'Automatically rejected due to change in delivery address affecting the ZIP code qualification for short haul / line haul.'
    ) {
      rejectedDSHorDLHServiceItem = true;
    }

    return (
      <React.Fragment key={`sit-alert-${id}`}>
        <tr key={id}>
          <td className={styles.nameAndDate}>
            <div className={styles.codeName}>
              <span className={styles.serviceItemName}>{serviceItem.serviceItem}</span>
              {(code === 'DCRT' || code === 'ICRT') && serviceItem.details.standaloneCrate && ' - Standalone'}
              {ALLOWED_RESUBMISSION_SI_CODES.includes(code) && resubmittedToolTip.isResubmitted ? (
                <ToolTip
                  data-testid="toolTipResubmission"
                  key={id}
                  text={resubmittedToolTip.formattedResubmissionDetails}
                  position="bottom"
                  color="#0050d8"
                  title={serviceItem.serviceItem}
                  closeOnLeave
                />
              ) : null}
              {ALLOWED_SIT_UPDATE_SI_CODES.includes(code) && hasPaymentRequestBeenMade ? (
                <ToolTip
                  text="This cannot be changed due to a payment request existing for this service item."
                  color="#d54309"
                  icon="circle-exclamation"
                  title={serviceItem.serviceItem}
                  closeOnLeave
                />
              ) : null}
            </div>
            <p>{getServiceItemDisplayDate(item)}</p>
          </td>
          <td className={styles.detail}>
            <ServiceItemDetails
              id={`service-${id}`}
              code={code}
              details={details}
              serviceRequestDocs={serviceRequestDocuments}
              serviceItem={serviceItem}
              shipment={shipment}
              sitStatus={shipment.sitStatus}
            />
          </td>
          <td>
            {statusForTableType === SERVICE_ITEM_STATUS.SUBMITTED && (
              <Restricted to={permissionTypes.updateMTOServiceItem}>
                <Restricted to={permissionTypes.updateMTOPage}>
                  <div className={styles.statusAction}>
                    <Button
                      type="button"
                      className="usa-button--icon usa-button--small acceptButton"
                      data-testid="acceptButton"
                      onClick={() => handleUpdateMTOServiceItemStatus(id, mtoShipmentID, SERVICE_ITEM_STATUS.APPROVED)}
                      disabled={isMoveLocked}
                    >
                      <span className="icon">
                        <FontAwesomeIcon icon="check" />
                      </span>
                      <span>Accept</span>
                    </Button>
                    <Button
                      type="button"
                      secondary
                      className="usa-button--small usa-button--icon margin-left-1 rejectButton"
                      data-testid="rejectButton"
                      onClick={() => handleShowRejectionDialog(id, mtoShipmentID)}
                      disabled={isMoveLocked}
                    >
                      <span className="icon">
                        <FontAwesomeIcon icon="times" />
                      </span>
                      <span>Reject</span>
                    </Button>
                  </div>
                </Restricted>
              </Restricted>
            )}
            {statusForTableType === SERVICE_ITEM_STATUS.APPROVED && (
              <Restricted to={permissionTypes.updateMTOServiceItem}>
                <Restricted to={permissionTypes.updateMTOPage}>
                  <div className={styles.statusAction}>
                    <Button
                      type="button"
                      data-testid="rejectTextButton"
                      className="text-blue usa-button--unstyled margin-left-1"
                      onClick={() => handleShowRejectionDialog(id, mtoShipmentID)}
                      disabled={isMoveLocked}
                    >
                      <span className="icon">
                        <FontAwesomeIcon icon="times" />
                      </span>{' '}
                      Reject
                    </Button>
                    {ALLOWED_SIT_UPDATE_SI_CODES.includes(code) && (
                      <div>
                        <Button
                          type="button"
                          data-testid="editTextButton"
                          className="text-blue usa-button--unstyled margin-left-1"
                          disabled={hasPaymentRequestBeenMade || isMoveLocked}
                          onClick={() => {
                            if (code === 'DDFSIT' || code === 'DOFSIT') {
                              handleShowEditSitEntryDateModal(id, mtoShipmentID);
                            } else {
                              handleShowEditSitAddressModal(id, mtoShipmentID);
                            }
                          }}
                        >
                          <span>
                            <FontAwesomeIcon icon="pencil" style={{ marginRight: '5px' }} />
                          </span>{' '}
                          Edit
                        </Button>
                      </div>
                    )}
                  </div>
                </Restricted>
              </Restricted>
            )}
            {statusForTableType === SERVICE_ITEM_STATUS.REJECTED && !rejectedDSHorDLHServiceItem && (
              <Restricted to={permissionTypes.updateMTOServiceItem}>
                <Restricted to={permissionTypes.updateMTOPage}>
                  <div className={styles.statusAction}>
                    <Button
                      type="button"
                      data-testid="approveTextButton"
                      className="text-blue usa-button--unstyled"
                      onClick={() => handleUpdateMTOServiceItemStatus(id, mtoShipmentID, SERVICE_ITEM_STATUS.APPROVED)}
                      disabled={isMoveLocked}
                    >
                      <span className="icon">
                        <FontAwesomeIcon icon="check" />
                      </span>{' '}
                      Approve
                    </Button>
                  </div>
                </Restricted>
              </Restricted>
            )}
          </td>
        </tr>
      </React.Fragment>
    );
  });

  return (
    <div className={classnames(styles.ServiceItemsTable, 'table--service-item', 'table--service-item--hasimg')}>
      <table>
        <thead className="table--small">
          <tr>
            <th>Service item</th>
            <th>Details</th>
            <th>&nbsp;</th>
          </tr>
        </thead>
        <tbody>{tableRows}</tbody>
      </table>
    </div>
  );
};

ServiceItemsTable.propTypes = {
  handleUpdateMTOServiceItemStatus: PropTypes.func.isRequired,
  handleShowRejectionDialog: PropTypes.func.isRequired,
  statusForTableType: PropTypes.string.isRequired,
  serviceItems: PropTypes.arrayOf(ServiceItemDetailsShape).isRequired,
  shipment: ShipmentShape,
};

ServiceItemsTable.defaultProps = {
  shipment: {},
};

export default ServiceItemsTable;
