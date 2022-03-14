import { Button, Tag } from '@trussworks/react-uswds';
import React, { Fragment } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './shipmentApprovalPreview.module.scss';

import { serviceItemCodes } from 'content/serviceItems';
import { mtoShipmentTypes } from 'constants/shipments';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import ShipmentInfoListSelector from 'components/Office/DefinitionLists/ShipmentInfoListSelector';
import ShipmentServiceItemsTable from 'components/Office/ShipmentServiceItemsTable/ShipmentServiceItemsTable';
import { Modal, ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import { MTOShipmentShape, OrdersInfoShape } from 'types/order';

const ShipmentApprovalPreview = ({
  mtoShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  setIsModalVisible,
  onSubmit,
  counselingFee,
  shipmentManagementFee,
}) => {
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={classNames('modal', styles.approvalPreviewModal)}>
          <div className={classNames(styles.containerTop)}>
            <div>
              <button
                type="button"
                title="Close shipment approval modal"
                onClick={() => setIsModalVisible(false)}
                className={classNames(styles.approvalClose, 'usa-button--unstyled')}
                data-testid="closeShipmentApproval"
              >
                <FontAwesomeIcon icon="times" title="Close modal" aria-label="Close modal" />
              </button>
            </div>
            <h2>Preview and post move task order</h2>
            <p>Is all the information shown correct and ready to send to Global Relocation Services?</p>
            <div className="display-flex">
              <Button type="submit" onClick={onSubmit}>
                Approve and send
              </Button>
              <Button type="reset" secondary onClick={() => setIsModalVisible(false)}>
                Back
              </Button>
            </div>
          </div>

          <hr className={styles.sectionBorder} />
          <h1 className={classNames(styles.customerName, 'text-normal')}>{customerInfo.name}</h1>
          <div className={classNames(styles.previewContainer, 'container')}>
            <h2>Requested Shipments</h2>
            {mtoShipments &&
              mtoShipments.map((shipment) => (
                <ShipmentContainer
                  key={shipment.id}
                  shipmentType={shipment.shipmentType}
                  className={classNames(styles.previewShipments)}
                >
                  <div className={styles.innerWrapper}>
                    <div className={styles.shipmentTypeHeading}>
                      <h3>{mtoShipmentTypes[shipment.shipmentType]}</h3>
                      {shipment.diversion && <Tag>diversion</Tag>}
                    </div>
                    <div className={styles.shipmentDetailWrapper}>
                      <ShipmentInfoListSelector
                        className={styles.shipmentInfo}
                        shipmentType={shipment.shipmentType}
                        isExpanded
                        shipment={{
                          ...shipment,
                          destinationAddress: shipment.destinationAddress
                            ? shipment.destinationAddress
                            : { postalCode: ordersInfo.newDutyLocation.address.postalCode },
                          agents: shipment.mtoAgents,
                        }}
                      />
                      <ShipmentServiceItemsTable
                        className={classNames(styles.shipmentServiceItems)}
                        shipmentType={shipment.shipmentType}
                      />
                    </div>
                  </div>
                </ShipmentContainer>
              ))}
          </div>
          <div className={classNames(styles.previewContainer, styles.basicMoveDetails, 'container')}>
            <h2>Basic move details</h2>
            {(shipmentManagementFee || counselingFee) && (
              <>
                <h4>Approved service items for this move</h4>
                <table className={classNames(styles.basicServiceItemsTable, 'table--stacked')}>
                  <tbody>
                    {shipmentManagementFee && (
                      <tr>
                        <td>{serviceItemCodes.MS}</td>
                      </tr>
                    )}
                    {counselingFee && (
                      <tr>
                        <td>{serviceItemCodes.CS}</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </>
            )}
            <h4>Allowances</h4>
            <AllowancesList info={allowancesInfo} />
            <h4>Customer info</h4>
            <CustomerInfoList customerInfo={customerInfo} />
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ShipmentApprovalPreview.propTypes = {
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape).isRequired,
  counselingFee: PropTypes.bool.isRequired,
  shipmentManagementFee: PropTypes.bool.isRequired,
  ordersInfo: OrdersInfoShape.isRequired,
  allowancesInfo: PropTypes.shape({
    branch: PropTypes.string,
    rank: PropTypes.string,
    weightAllowance: PropTypes.number,
    authorizedWeight: PropTypes.number,
    progear: PropTypes.number,
    spouseProgear: PropTypes.number,
    storageInTransit: PropTypes.number,
    dependents: PropTypes.bool,
  }).isRequired,
  customerInfo: PropTypes.shape({
    name: PropTypes.string,
    dodId: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
    currentAddress: PropTypes.shape({
      streetAddress1: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      postalCode: PropTypes.string,
    }),
    backupContact: PropTypes.shape({
      name: PropTypes.string,
      phone: PropTypes.string,
      email: PropTypes.string,
    }),
  }).isRequired,
  setIsModalVisible: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export default ShipmentApprovalPreview;
