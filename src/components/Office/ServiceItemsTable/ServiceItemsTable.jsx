import React from 'react';
import PropTypes from 'prop-types';
import { Button, Tag, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './ServiceItemsTable.module.scss';

import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES, SIT_ADDRESS_UPDATE_STATUS } from 'constants/sitUpdates';
import { formatDateFromIso } from 'utils/formatters';
import ServiceItemDetails from 'components/Office/ServiceItemDetails/ServiceItemDetails';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { selectDateFieldByStatus, selectDatePrefixByStatus } from 'utils/dates';
import ToolTip from 'shared/ToolTip/ToolTip';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

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

  const getServiceCodeDescription = (code) => {
    let content;
    switch (code) {
      case SERVICE_ITEM_CODES.DOFSIT:
        content = 'The first day of Origin SIT.';
        break;
      case SERVICE_ITEM_CODES.DOASIT:
        content = 'Additional days of Origin SIT that occur after the first day of SIT.';
        break;
      case SERVICE_ITEM_CODES.DOPSIT:
        content = 'Picking up items from the home prior to beginning the first day of Origin SIT.';
        break;
      case SERVICE_ITEM_CODES.DDFSIT:
        content = 'The first day of Destination SIT.';
        break;
      case SERVICE_ITEM_CODES.DDASIT:
        content = 'Additional days of Destination SIT that occur after the first day of Destination SIT.';
        break;
      case SERVICE_ITEM_CODES.DDDSIT:
        content = 'Delivery of items to home following Destination SIT.';
        break;
      default:
        content = 'No definition provided.';
    }
    return content;
  };

  const tableRows = serviceItems.map((serviceItem, index) => {
    const { id, code, details, mtoShipmentID, sitAddressUpdates, serviceRequestDocuments, ...item } = serviceItem;
    const { makeVisible, alertType, alertMessage } = serviceItemAddressUpdateAlert;

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
            <p className={styles.codeName}>
              {serviceItem.serviceItem}{' '}
              <ToolTip text={getServiceCodeDescription(code)} color="black" position="right" />
            </p>
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
