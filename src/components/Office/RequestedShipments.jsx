import React, { useState } from 'react';
import { useFormik } from 'formik';
import * as PropTypes from 'prop-types';
import classNames from 'classnames/bind';
import { Button, Checkbox, Fieldset, Modal, Overlay, ModalContainer } from '@trussworks/react-uswds';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import ShipmentDisplay from 'components/Office/ShipmentDisplay';
import styles from './requestedShipments.module.scss';

const cx = classNames.bind(styles);

const RequestedShipments = ({ mtoShipments }) => {
  const [isModalVisible, setIsModalVisible] = useState(false);

  const handleApprovalClick = () => {
    setIsModalVisible(true);
  };

  const formik = useFormik({
    initialValues: {
      shipmentManagementFee: false,
      counselingFee: false,
      shipments: [],
    },
    onSubmit: () => {
      handleApprovalClick();
    },
  });

  return (
    <div className={`${cx('requested-shipments')} container`} data-cy="requested-shipments">
      <div id="approvalConfirmationModal" style={{ display: isModalVisible ? 'block' : 'none' }}>
        <Overlay />
        <ModalContainer>
          <Modal>
            <div className={`${cx('approval-close')}`}>
              <FontAwesomeIcon
                aria-hidden
                icon={faTimes}
                title="Close shipment approval modal"
                onClick={() => setIsModalVisible(false)}
                className={`${cx('approval-close')} icon`}
              />
            </div>
            <h1>Preview and post move task order</h1>
          </Modal>
        </ModalContainer>
      </div>
      <h4>Requested shipments</h4>
      <form onSubmit={formik.handleSubmit}>
        <div className={`${cx('__content')}`}>
          {mtoShipments &&
            mtoShipments.map((shipment, i) => (
              <ShipmentDisplay
                key={shipment.id}
                index={i}
                shipmentId={shipment.id}
                shipmentType={shipment.shipmentType}
                displayInfo={{
                  heading: shipment.shipmentType,
                  requestedMoveDate: shipment.requestedPickupDate,
                  currentAddress: shipment.pickupAddress,
                  destinationAddress: shipment.destinationAddress,
                }}
                /* eslint-disable-next-line react/jsx-props-no-spreading */
                {...formik.getFieldProps(`shipments[${i}]`)}
              />
            ))}
        </div>
        <div>
          <h3>Add service items to this move</h3>
          <span>{isModalVisible}</span>
          <Fieldset legend="MTO service items" legendSrOnly id="input-type-fieldset">
            <Checkbox
              id="shipmentManagementFee"
              label="Shipment management fee"
              name="shipmentManagementFee"
              onChange={formik.handleChange}
            />
            <Checkbox id="counselingFee" label="Counseling fee" name="counselingFee" onChange={formik.handleChange} />
          </Fieldset>
          <Button
            id="shipmentApproveButton"
            className={`${cx('usa-button--small')} usa-button--icon`}
            onClick={formik.handleSubmit}
            type="submit"
          >
            <span>Approve selected shipments</span>
          </Button>
        </div>
      </form>
    </div>
  );
};

RequestedShipments.propTypes = {
  // eslint-disable-next-line react/forbid-prop-types
  mtoShipments: PropTypes.array.isRequired,
};

export default RequestedShipments;
