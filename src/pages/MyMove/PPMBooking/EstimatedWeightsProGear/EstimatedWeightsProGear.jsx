import React, { useState } from 'react';
import { connect } from 'react-redux';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import { customerRoutes } from 'constants/routes';
import EstimatedWeightsProGearForm from 'components/Customer/PPMBooking/EstimatedWeightsProGearForm/EstimatedWeightsProGearForm';
import { shipmentTypes } from 'constants/shipments';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { MtoShipmentShape, OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import {
  selectCurrentOrders,
  selectMTOShipmentById,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';

const EstimatedWeightsProGear = ({ orders, serviceMember, mtoShipment }) => {
  const [errorMessage, setErrorMessage] = useState();
  const history = useHistory();
  const { moveId, shipmentNumber } = useParams();

  const mtoShipmentId = mtoShipment.id;

  const handleBack = () => {
    history.push(generatePath(customerRoutes.SHIPMENT_EDIT_PATH, { moveId, mtoShipmentId }));
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);

    const hasProGear = values.hasProGear === 'true';

    const payload = {
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        estimatedWeight: values.estimatedWeight,
        hasProGear,
        proGearWeight: hasProGear ? values.proGearWeight : null,
        spouseProGearWeight: hasProGear ? values.spouseProGearWeight : null,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        updateMTOShipment(response);
        history.push(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, { moveId, mtoShipmentId });
      })
      .catch((err) => {
        setSubmitting(false);

        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
      });
  };

  return (
    <div>
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

EstimatedWeightsProGear.propTypes = {
  orders: OrdersShape.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
  mtoShipment: MtoShipmentShape.isRequired,
};

function mapStateToProps(state, ownProps) {
  return {
    orders: selectCurrentOrders(state) || {},
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId) || {},
  };
}

export default connect(mapStateToProps)(EstimatedWeightsProGear);
