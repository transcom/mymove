import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { customerRoutes } from 'constants/routes';
import EstimatedWeightsProGearForm from 'components/Customer/PPM/Booking/EstimatedWeightsProGearForm/EstimatedWeightsProGearForm';
import { shipmentTypes } from 'constants/shipments';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import {
  selectCurrentOrders,
  selectMTOShipmentById,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ScrollToTop from 'components/ScrollToTop';

const EstimatedWeightsProGear = () => {
  const [errorMessage, setErrorMessage] = useState();
  const history = useHistory();
  const { moveId, mtoShipmentId, shipmentNumber } = useParams();
  const dispatch = useDispatch();

  const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
  const orders = useSelector((state) => selectCurrentOrders(state));
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const handleBack = () => {
    history.push(generatePath(customerRoutes.SHIPMENT_EDIT_PATH, { moveId, mtoShipmentId }));
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);

    const hasProGear = values.hasProGear === 'true';

    const payload = {
      shipmentType: mtoShipment.shipmentType,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        estimatedWeight: Number(values.estimatedWeight),
        hasProGear,
        proGearWeight: hasProGear ? Number(values.proGearWeight) : null,
        spouseProGearWeight: hasProGear ? Number(values.spouseProGearWeight) : null,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));
        history.push(generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, { moveId, mtoShipmentId }));
      })
      .catch((err) => {
        setSubmitting(false);
        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
      });
  };

  if (!serviceMember || !orders || !mtoShipment) {
    return <LoadingPlaceholder />;
  }

  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <ScrollToTop otherDep={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>Estimated weight</h1>
            {errorMessage && (
              <Alert slim type="error">
                {errorMessage}
              </Alert>
            )}
            <EstimatedWeightsProGearForm
              orders={orders}
              serviceMember={serviceMember}
              mtoShipment={mtoShipment}
              onSubmit={handleSubmit}
              onBack={handleBack}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default EstimatedWeightsProGear;
