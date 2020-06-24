import { Modal, ModalContainer, Overlay } from '@trussworks/react-uswds';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';
import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames/bind';

import { mtoShipmentTypeToFriendlyDisplay } from '../../shared/formatters';

import styles from './shipmentApprovalPreview.module.scss';
import AllowancesTable from './AllowancesTable';
import CustomerInfoTable from './CustomerInfoTable';
import ShipmentContainer from './ShipmentContainer';

const cx = classNames.bind(styles);

const ShipmentApprovalPreview = ({ mtoShipments, allowancesInfo, customerInfo, mtoAgents, setIsModalVisible }) => {
  const getAgents = (shipment) => {
    return mtoAgents.filter((agent) => agent.shipmentId === shipment.id);
  };
  const shipmentsWithAgents = mtoAgents
    ? mtoShipments.slice().map((shipment) => ({ ...shipment, agents: getAgents(shipment) }))
    : mtoShipments;

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className="padding-4 overflow-y-auto maxh-viewport">
          <div className={`${cx('approval-close')}`}>
            <FontAwesomeIcon
              aria-hidden
              icon={faTimes}
              title="Close shipment approval modal"
              onClick={() => setIsModalVisible(false)}
              className={`${cx('approval-close')} icon`}
            />
          </div>
          <h3 className="text-bold">Preview and post move task order</h3>
          <h2 className="text-normal">{customerInfo.name}</h2>
          <div className="container">
            <h2>Requested Shipments</h2>
            {shipmentsWithAgents &&
              shipmentsWithAgents.map((shipment) => (
                <ShipmentContainer key={shipment.id} shipmentType={shipment.shipmentType}>
                  <div>
                    <h4 className="text-normal">{mtoShipmentTypeToFriendlyDisplay(shipment.shipmentType)}</h4>
                    <table className="table--stacked">
                      <tbody>
                        <tr>
                          <td>Requested Move Date</td>
                          <td>{shipment.requestedPickupDate}</td>
                        </tr>
                        <tr>
                          <td>Current Address</td>
                          <td>
                            {shipment.pickupAddress.street_address_1}
                            <br />
                            {shipment.pickupAddress.city}, {shipment.pickupAddress.state}{' '}
                            {shipment.pickupAddress.postal_code}
                          </td>
                        </tr>
                        <tr>
                          <td>Destination Address</td>
                          <td>
                            {shipment.destinationAddress.street_address_1}
                            <br />
                            {shipment.destinationAddress.city}, {shipment.destinationAddress.state}{' '}
                            {shipment.destinationAddress.postal_code}
                          </td>
                        </tr>
                        <tr>
                          <td>Customer Remarks</td>
                          <td>{shipment.customerRemarks}</td>
                        </tr>
                        {mtoAgents &&
                          mtoAgents.map((agent) => (
                            <tr>
                              <td>{agent.type === 'RELEASING_AGENT' ? 'Releasing Agent' : 'Receiving Agent'}</td>
                              <td>
                                {agent.firstName} {agent.lastName}
                                <br />
                                {agent.phone} <br /> {agent.email}
                              </td>
                            </tr>
                          ))}
                      </tbody>
                    </table>
                  </div>
                </ShipmentContainer>
              ))}
          </div>
          <div className="container">
            <h2>Basic move details</h2>
            <h4>Approved service items for this move</h4>
            <table className="table--stacked">
              <tbody>
                <tr>
                  <td>Shipment management fee</td>
                </tr>
                <tr>
                  <td>Counseling fee</td>
                </tr>
              </tbody>
            </table>
            <AllowancesTable info={allowancesInfo} />
            <CustomerInfoTable customerInfo={customerInfo} />
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ShipmentApprovalPreview.propTypes = {
  // eslint-disable-next-line react/forbid-prop-types
  mtoShipments: PropTypes.array.isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  mtoAgents: PropTypes.array,
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
    destinationAddress: PropTypes.shape({
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
};

ShipmentApprovalPreview.defaultProps = {
  mtoAgents: [],
};

export default ShipmentApprovalPreview;
