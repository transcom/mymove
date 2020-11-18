import { Button, Modal, ModalContainer, Overlay } from '@trussworks/react-uswds';
import React, { Fragment } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTimes as XIcon } from '@fortawesome/free-solid-svg-icons';

import { mtoShipmentTypeToFriendlyDisplay } from '../../shared/formatters';
import { MTOAgentShape, MTOShipmentShape, OrdersInfoShape } from '../../types/moveOrder';
import { formatAddress } from '../../utils/shipmentDisplay';

import styles from './shipmentApprovalPreview.module.scss';
import AllowancesTable from './AllowancesTable';
import CustomerInfoTable from './CustomerInfoTable';
import ShipmentContainer from './ShipmentContainer';
import ShipmentServiceItemsTable from './ShipmentServiceItemsTable/ShipmentServiceItemsTable';

const ShipmentApprovalPreview = ({
  mtoShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  mtoAgents,
  setIsModalVisible,
  onSubmit,
  counselingFee,
  shipmentManagementFee,
}) => {
  const getAgents = (shipment) => {
    return mtoAgents.filter((agent) => agent.shipmentId === shipment.id);
  };
  const shipmentsWithAgents = mtoAgents
    ? mtoShipments.map((shipment) => ({ ...shipment, agents: getAgents(shipment) }))
    : mtoShipments;

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
                <FontAwesomeIcon icon={XIcon} title="Close modal" aria-label="Close modal" />
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
            {shipmentsWithAgents &&
              shipmentsWithAgents.map((shipment) => (
                <ShipmentContainer
                  key={shipment.id}
                  shipmentType={shipment.shipmentType}
                  className={classNames(styles.previewShipments)}
                >
                  <div className={styles.innerWrapper}>
                    <h4 className="text-normal">{mtoShipmentTypeToFriendlyDisplay(shipment.shipmentType)}</h4>
                    <div className="display-flex">
                      <table className={classNames('table--stacked', styles.shipmentInfo)}>
                        <tbody>
                          <tr>
                            <th className="text-bold" scope="row">
                              Requested Move Date
                            </th>
                            <td>{shipment.requestedPickupDate}</td>
                          </tr>
                          <tr>
                            <th className="text-bold" scope="row">
                              Current Address
                            </th>
                            <td>{shipment.pickupAddress && formatAddress(shipment.pickupAddress)}</td>
                          </tr>
                          <tr>
                            <th className="text-bold" scope="row">
                              Destination Address
                            </th>
                            <td data-testid="destinationAddress">
                              {shipment.destinationAddress
                                ? formatAddress(shipment.destinationAddress)
                                : ordersInfo.newDutyStation.address.postal_code}
                            </td>
                          </tr>
                          <tr>
                            <th className="text-bold" scope="row">
                              Customer Remarks
                            </th>
                            <td>{shipment.customerRemarks}</td>
                          </tr>
                          {mtoAgents &&
                            mtoAgents.map((agent) => (
                              <Fragment key={`${agent.type}-${agent.email}`}>
                                <tr>
                                  <th className="text-bold" scope="row">
                                    {agent.type === 'RELEASING_AGENT' ? 'Releasing Agent' : 'Receiving Agent'}
                                  </th>
                                  <td>
                                    {agent.firstName} {agent.lastName}
                                    <br />
                                    {agent.phone} <br /> {agent.email}
                                  </td>
                                </tr>
                              </Fragment>
                            ))}
                        </tbody>
                      </table>
                      <ShipmentServiceItemsTable
                        className={classNames(styles.shipmentServiceItems)}
                        shipmentType={shipment.shipmentType}
                      />
                    </div>
                  </div>
                </ShipmentContainer>
              ))}
          </div>
          <div className={classNames(styles.previewContainer, 'container')}>
            <h2>Basic move details</h2>
            {(shipmentManagementFee || counselingFee) && (
              <>
                <h4 className={classNames(styles.tableH4)}>Approved service items for this move</h4>
                <table className="table--stacked">
                  <tbody>
                    {shipmentManagementFee && (
                      <tr>
                        <td>Shipment management fee</td>
                      </tr>
                    )}
                    {counselingFee && (
                      <tr>
                        <td>Counseling fee</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </>
            )}
            <AllowancesTable info={allowancesInfo} />
            <CustomerInfoTable customerInfo={customerInfo} />
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ShipmentApprovalPreview.propTypes = {
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape).isRequired,
  mtoAgents: PropTypes.arrayOf(MTOAgentShape),
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
      street_address_1: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      postal_code: PropTypes.string,
    }),
    backupContactName: PropTypes.string,
    backupContactPhone: PropTypes.string,
    backupContactEmail: PropTypes.string,
  }).isRequired,
  setIsModalVisible: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

ShipmentApprovalPreview.defaultProps = {
  mtoAgents: [],
};

export default ShipmentApprovalPreview;
