import React from 'react';
import PropTypes from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
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

const ServiceItemsTable = ({
  serviceItems,
  statusForTableType,
  handleUpdateMTOServiceItemStatus,
  handleRequestSITAddressUpdateModal,
  handleShowRejectionDialog,
}) => {
  let dateField;
  switch (statusForTableType) {
    case SERVICE_ITEM_STATUS.SUBMITTED:
      dateField = 'createdAt';
      break;
    case SERVICE_ITEM_STATUS.APPROVED:
      dateField = 'approvedAt';
      break;
    case SERVICE_ITEM_STATUS.REJECTED:
      dateField = 'rejectedAt';
      break;
    default:
      dateField = 'createdAt';
  }

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

  const tableRows = serviceItems.map(
    ({ id, code, serviceItem, details, mtoShipmentID, sitAddressUpdates, ...item }, index) => {
      return (
        <>
          {ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES.includes(code) &&
            sitAddressUpdates &&
            showSITAddressUpdateRequestedTag(code, sitAddressUpdates) && (
              <tr key={index}>
                <td colSpan={3} style={{ borderBottom: 'none', paddingBottom: '0', paddingTop: '8px' }}>
                  <Tag data-testid="sitAddressUpdateTag">UPDATE REQUESTED</Tag>
                </td>
              </tr>
            )}
          <tr key={id}>
            <td className={styles.nameAndDate}>
              <p className={styles.codeName}>{serviceItem}</p>
              <p>{formatDateFromIso(item[`${dateField}`], 'DD MMM YYYY')}</p>
            </td>
            <td className={styles.detail}>
              <ServiceItemDetails id={id} code={code} details={details} />
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
                      <Button
                        type="button"
                        data-testid="editTextButton"
                        className="text-blue usa-button--unstyled margin-left-1"
                        onClick={() =>
                          handleRequestSITAddressUpdateModal({
                            id,
                            code,
                            details,
                            serviceItem,
                            submittedAt: item.submittedAt,
                          })
                        }
                      >
                        <span>
                          <FontAwesomeIcon icon="pencil" style={{ marginRight: '5px' }} />
                        </span>{' '}
                        {sitAddressUpdates && hasSITAddressUpdate(sitAddressUpdates) ? 'Review Request' : 'Edit'}
                      </Button>
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
                        <FontAwesomeIcon icon="times" />
                      </span>{' '}
                      Approve
                    </Button>
                  </div>
                </Restricted>
              )}
            </td>
          </tr>
        </>
      );
    },
  );

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
  serviceItems: PropTypes.arrayOf(ServiceItemDetailsShape).isRequired,
};

export default ServiceItemsTable;
