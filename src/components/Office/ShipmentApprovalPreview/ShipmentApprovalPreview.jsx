import { Button, Tag } from '@trussworks/react-uswds';
import React, { Fragment, useEffect, useState } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './shipmentApprovalPreview.module.scss';

import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { Modal, ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import ShipmentInfoListSelector from 'components/Office/DefinitionLists/ShipmentInfoListSelector';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import ShipmentServiceItemsTable from 'components/Office/ShipmentServiceItemsTable/ShipmentServiceItemsTable';
import { serviceItemCodes } from 'content/serviceItems';
import { mtoShipmentTypes } from 'constants/shipments';
import { OrdersInfoShape } from 'types/order';
import { ShipmentShape } from 'types/shipment';

const ShipmentApprovalPreview = ({
  mtoShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  setIsModalVisible,
  onSubmit,
  counselingFee,
  shipmentManagementFee,
  isSubmitting,
}) => {
  const [isOconusMove, setIsOconusMove] = useState(false);
  const [IsMoveApprovalButtonDisabled, setIsMoveApprovalButtonDisabled] = useState(false);

  useEffect(() => {
    const fetchFlag = async () => {
      setIsMoveApprovalButtonDisabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.DISABLE_MOVE_APPROVAL));
    };
    fetchFlag();
  }, []);

  useEffect(() => {
    // Check if duty locations on the orders qualify as OCONUS to conditionally render the UB allowance details
    if (ordersInfo?.currentDutyLocation?.address?.isOconus || ordersInfo?.newDutyLocation?.address?.isOconus) {
      setIsOconusMove(true);
    } else {
      setIsOconusMove(false);
    }
  }, [ordersInfo?.currentDutyLocation, ordersInfo?.newDutyLocation]);
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={classNames('modal', styles.approvalPreviewModal)} onClose={() => setIsModalVisible(false)}>
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
            <div className={classNames(styles.displayFlexRight)}>
              <Button type="reset" secondary onClick={() => setIsModalVisible(false)}>
                Back
              </Button>
              <Button type="submit" onClick={onSubmit} disabled={isSubmitting || IsMoveApprovalButtonDisabled}>
                Approve and send
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
                      <div className={styles.typeLocator}>
                        <h3>{mtoShipmentTypes[shipment.shipmentType]}</h3>
                        <h5>#{shipment.shipmentLocator}</h5>
                      </div>
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
                        shipment={shipment}
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
            <AllowancesList info={allowancesInfo} isOconusMove={isOconusMove} />
            <h4>Customer info</h4>
            <CustomerInfoList customerInfo={customerInfo} />
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ShipmentApprovalPreview.propTypes = {
  mtoShipments: PropTypes.arrayOf(ShipmentShape).isRequired,
  counselingFee: PropTypes.bool.isRequired,
  shipmentManagementFee: PropTypes.bool.isRequired,
  ordersInfo: OrdersInfoShape.isRequired,
  allowancesInfo: PropTypes.shape({
    branch: PropTypes.string,
    grade: PropTypes.string,
    totalWeight: PropTypes.string,
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
