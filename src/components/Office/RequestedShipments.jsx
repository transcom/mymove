import React, { useState } from 'react';
import { useFormik } from 'formik';
import * as PropTypes from 'prop-types';
import classNames from 'classnames/bind';
import { Button, Checkbox, Fieldset } from '@trussworks/react-uswds';

import { MTOAgentShape, MTOShipmentShape } from '../../types/moveOrder';

import ShipmentApprovalPreview from './ShipmentApprovalPreview';
import styles from './requestedShipments.module.scss';

import ShipmentDisplay from 'components/Office/ShipmentDisplay';

const cx = classNames.bind(styles);

const RequestedShipments = ({ mtoShipments, allowancesInfo, customerInfo, mtoAgents }) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [filteredShipments, setFilteredShipments] = useState([]);

  const handleApprovalClick = () => {
    setIsModalVisible(true);
  };

  const filterShipments = (formikShipmentIds) => {
    return mtoShipments.filter(({ id }) => formikShipmentIds.includes(id));
  };

  const formik = useFormik({
    initialValues: {
      shipmentManagementFee: false,
      counselingFee: false,
      shipments: [],
    },
    onSubmit: () => {
      setFilteredShipments(filterShipments(formik.values.shipments));
      handleApprovalClick();
    },
  });

  const isButtonEnabled =
    formik.values.shipments.length > 0 && (formik.values.counselingFee || formik.values.shipmentManagementFee);

  return (
    <div className={`${cx('requested-shipments')} container`} data-cy="requested-shipments">
      <div id="approvalConfirmationModal" style={{ display: isModalVisible ? 'block' : 'none' }}>
        <ShipmentApprovalPreview
          mtoShipments={filteredShipments}
          allowancesInfo={allowancesInfo}
          customerInfo={customerInfo}
          setIsModalVisible={setIsModalVisible}
          mtoAgents={mtoAgents}
          counselingFee={formik.values.counselingFee}
          shipmentManagementFee={formik.values.shipmentManagementFee}
        />
      </div>
      <h4>Requested shipments</h4>
      <form onSubmit={formik.handleSubmit}>
        <div className={`${cx('__content')}`}>
          {mtoShipments &&
            mtoShipments.map((shipment) => (
              <ShipmentDisplay
                key={shipment.id}
                shipmentId={shipment.id}
                shipmentType={shipment.shipmentType}
                displayInfo={{
                  heading: shipment.shipmentType,
                  requestedMoveDate: shipment.requestedPickupDate,
                  currentAddress: shipment.pickupAddress,
                  destinationAddress: shipment.destinationAddress,
                }}
                /* eslint-disable-next-line react/jsx-props-no-spreading */
                {...formik.getFieldProps(`shipments`)}
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
            disabled={!isButtonEnabled}
          >
            <span>Approve selected shipments</span>
          </Button>
        </div>
      </form>
    </div>
  );
};

RequestedShipments.propTypes = {
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape).isRequired,
  mtoAgents: PropTypes.arrayOf(MTOAgentShape),
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
};

RequestedShipments.defaultProps = {
  mtoAgents: [],
};

export default RequestedShipments;
