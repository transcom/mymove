import React from 'react';
import PropTypes from 'prop-types';
import { Button, Tag, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useParams } from 'react-router';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './ServiceItemsTable.module.scss';

import { SERVICE_ITEM_STATUS } from 'shared/constants';
import {
  ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES,
  SIT_ADDRESS_UPDATE_STATUS,
  ALLOWED_RESUBMISSION_SI_CODES,
} from 'constants/sitUpdates';
import { formatDateFromIso } from 'utils/formatters';
import ServiceItemDetails from 'components/Office/ServiceItemDetails/ServiceItemDetails';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { selectDateFieldByStatus, selectDatePrefixByStatus } from 'utils/dates';
import { useGHCGetMoveHistory, useMovePaymentRequestsQueries } from 'hooks/queries';
import ToolTip from 'shared/ToolTip/ToolTip';

const ServiceItemsTable = ({
  serviceItems,
  statusForTableType,
  handleUpdateMTOServiceItemStatus,
  handleRequestSITAddressUpdateModal,
  handleShowRejectionDialog,
  handleShowEditSitAddressModal,
  handleShowEditSitEntryDateModal,
  serviceItemAddressUpdateAlert,
}) => {
  const hasSITAddressUpdate = (sitAddressUpdates) => {
    const requestedAddressUpdates = sitAddressUpdates.filter((s) => s.status === SIT_ADDRESS_UPDATE_STATUS.REQUESTED);
    return requestedAddressUpdates.length > 0;
  };

  const showSITAddressUpdateRequestedTag = (code, sitAddressUpdates) => {
    return (
      statusForTableType === SERVICE_ITEM_STATUS.APPROVED &&
      ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code) &&
      hasSITAddressUpdate(sitAddressUpdates)
    );
  };

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
      historyDataForMove.map((obj) => {
        let newestEventInAuditHistory = historyDataForMove[0].actionTstampTx;
        // object id of the audit history entry should match the id of the service item
        if (obj.objectId === serviceItemId) {
          // if time of curr obj is newer than the curr newestEventInAuditHistory
          if (obj.actionTstampTx > newestEventInAuditHistory) {
            newestEventInAuditHistory = obj.actionTstampTx;
            newestHistoryData = obj;
          }
        }
        return null;
      });
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
    const keys = Object.keys(details.changedValues);
    let resultStringToDisplay = '';
    keys.forEach((key) => {
      const formattedKeyString = formatKeyStringsForToolTip(key);
      const newValue = details.changedValues[key];
      const oldValue = details.oldValues[key];
      resultStringToDisplay += `${formattedKeyString}\nNew: ${newValue} \nPrevious: ${oldValue}\n\n`;
    });
    return resultStringToDisplay;
  }

  const history = useGHCGetMoveHistory({ moveCode });
  const renderToolTipWithOldDataIfResubmission = (serviceItemId) => {
    const historyDataForMove = history.queueResult.data;
    const historyDataForServiceItem = getNewestHistoryDataForServiceItem(historyDataForMove, serviceItemId);
    const isResubmitted = getResubmissionStatus(historyDataForServiceItem);
    if (isResubmitted) {
      return (
        <ToolTip
          data-testid="toolTipResubmission"
          key={serviceItemId}
          text={generateResubmissionDetailsText(historyDataForServiceItem)}
          position="bottom"
          color="#0050d8"
        />
      );
    }
    return null;
  };

  const tableRows = serviceItems.map((serviceItem, index) => {
    const { id, code, details, mtoShipmentID, sitAddressUpdates, serviceRequestDocuments, ...item } = serviceItem;
    const { makeVisible, alertType, alertMessage } = serviceItemAddressUpdateAlert;
    let hasPaymentRequestBeenMade;
    // if there are service items in the payment requests, we want to look to see if the service item is in there
    // if so, we don't want to let the TOO edit the SIT entry date
    if (serviceItemInPaymentRequests && ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code)) {
      hasPaymentRequestBeenMade = isServiceItemFoundInPaymentRequests(id);
    }

    return (
      <React.Fragment key={`sit-alert-${id}`}>
        {ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code) &&
          sitAddressUpdates &&
          showSITAddressUpdateRequestedTag(code, sitAddressUpdates) && (
            <tr key={index}>
              <td colSpan={3} style={{ borderBottom: 'none', paddingBottom: '0', paddingTop: '8px' }}>
                <Tag data-testid="sitAddressUpdateTag">UPDATE REQUESTED</Tag>
              </td>
            </tr>
          )}
        {ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code) && makeVisible && (
          <tr key={`sit-alert-${id}`}>
            <td style={{ border: 'none', paddingBottom: '0' }} colSpan={3}>
              <Alert type={alertType} slim data-testid="serviceItemAddressUpdateAlert">
                {alertMessage}
              </Alert>
            </td>
          </tr>
        )}
        <tr key={id}>
          <td className={styles.nameAndDate}>
            <div className={styles.codeName}>
              {serviceItem.serviceItem}{' '}
              {ALLOWED_RESUBMISSION_SI_CODES.includes(code) && renderToolTipWithOldDataIfResubmission(id)}
              {ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code) && hasPaymentRequestBeenMade ? (
                <ToolTip
                  text="This cannot be changed due to a payment request existing for this service item."
                  color="#d54309"
                  icon="circle-exclamation"
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
            />
          </td>
          <td>
            {statusForTableType === SERVICE_ITEM_STATUS.SUBMITTED && (
              <Restricted to={permissionTypes.updateMTOServiceItem}>
                <div className={styles.statusAction}>
                  <Button
                    type="button"
                    className="usa-button--icon usa-button--small acceptButton"
                    data-testid="acceptButton"
                    onClick={() => handleUpdateMTOServiceItemStatus(id, mtoShipmentID, SERVICE_ITEM_STATUS.APPROVED)}
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
                  >
                    <span className="icon">
                      <FontAwesomeIcon icon="times" />
                    </span>
                    <span>Reject</span>
                  </Button>
                </div>
              </Restricted>
            )}
            {statusForTableType === SERVICE_ITEM_STATUS.APPROVED && (
              <Restricted to={permissionTypes.updateMTOServiceItem}>
                <div className={styles.statusAction}>
                  <Button
                    type="button"
                    data-testid="rejectTextButton"
                    className="text-blue usa-button--unstyled margin-left-1"
                    onClick={() => handleShowRejectionDialog(id, mtoShipmentID)}
                  >
                    <span className="icon">
                      <FontAwesomeIcon icon="times" />
                    </span>{' '}
                    Reject
                  </Button>
                  {ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code) && (
                    <div>
                      {sitAddressUpdates && hasSITAddressUpdate(sitAddressUpdates) ? (
                        <Button
                          type="button"
                          data-testid="reviewRequestTextButton"
                          className="text-blue usa-button--unstyled margin-left-1"
                          onClick={() => handleRequestSITAddressUpdateModal(id, mtoShipmentID)}
                        >
                          <span>
                            <FontAwesomeIcon icon="pencil" style={{ marginRight: '5px' }} />
                          </span>{' '}
                          Review Request
                        </Button>
                      ) : (
                        <Button
                          type="button"
                          data-testid="editTextButton"
                          className="text-blue usa-button--unstyled margin-left-1"
                          disabled={hasPaymentRequestBeenMade}
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
                      )}
                    </div>
                  )}
                </div>
              </Restricted>
            )}
            {statusForTableType === SERVICE_ITEM_STATUS.REJECTED && (
              <Restricted to={permissionTypes.updateMTOServiceItem}>
                <div className={styles.statusAction}>
                  <Button
                    type="button"
                    data-testid="approveTextButton"
                    className="text-blue usa-button--unstyled"
                    onClick={() => handleUpdateMTOServiceItemStatus(id, mtoShipmentID, SERVICE_ITEM_STATUS.APPROVED)}
                  >
                    <span className="icon">
                      <FontAwesomeIcon icon="check" />
                    </span>{' '}
                    Approve
                  </Button>
                </div>
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

ServiceItemsTable.defaultProps = {
  handleRequestSITAddressUpdateModal: () => {},
};

ServiceItemsTable.propTypes = {
  handleUpdateMTOServiceItemStatus: PropTypes.func.isRequired,
  handleShowRejectionDialog: PropTypes.func.isRequired,
  statusForTableType: PropTypes.string.isRequired,
  handleRequestSITAddressUpdateModal: PropTypes.func,
  serviceItemAddressUpdateAlert: PropTypes.object.isRequired,
  serviceItems: PropTypes.arrayOf(ServiceItemDetailsShape).isRequired,
};

export default ServiceItemsTable;
