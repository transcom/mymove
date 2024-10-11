import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import { customerRoutes } from 'constants/routes';
import { updateMTOShipment } from 'store/entities/actions';
import {
  selectCurrentOrders,
  selectMTOShipmentById,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';

const MobileHomeInfo = () => {
  const navigate = useNavigate();
  const { moveId, mtoShipmentId } = useParams();
  const dispatch = useDispatch();

  const handleBack = () => {
    navigate(generatePath(customerRoutes.SHIPMENT_EDIT_PATH, { moveId, mtoShipmentId }));
  };

  const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
  const orders = useSelector((state) => selectCurrentOrders(state));
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  // Loading placeholder while data loads
  if (!serviceMember || !orders || !mtoShipment) {
    return <LoadingPlaceholder />;
  }

  return (
    <MtoShipmentForm
      mtoShipment={mtoShipment}
      shipmentType={mtoShipment.shipmentType}
      currentResidence={serviceMember?.residential_address}
      newDutyLocationAddress={orders?.new_duty_location?.address}
      updateMTOShipment={(shipment) => dispatch(updateMTOShipment(shipment))}
      serviceMember={serviceMember}
      orders={orders}
      handleBack={handleBack}
    />
  );
};

export default MobileHomeInfo;
