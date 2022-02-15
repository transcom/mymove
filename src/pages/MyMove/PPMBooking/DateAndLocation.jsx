import React from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
import { DutyStationShape } from 'types';
import DateAndLocationForm from 'components/Customer/PPMBooking/DateAndLocationForm/DateAndLocationForm';
import { validatePostalCode } from 'utils/validation';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createMTOShipment, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const DateAndLocation = ({ mtoShipment, serviceMember, destinationDutyLocation }) => {
  const history = useHistory();
  const { moveId } = useParams();

  const isNewShipment = !mtoShipment?.id;

  const handleBack = () => {
    if (isNewShipment) {
      return history.push(generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId }));
    }

    return history.push(generalRoutes.HOME_PATH);
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    const createOrUpdateShipment = {
      moveTaskOrderID: moveId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        pickupPostalCode: values.pickupPostalCode,
        destinationPostalCode: values.destinationPostalCode,
        sitExpected: values.sitExpected,
        expectedDepartureDate: values.expectedDepartureDate,
      },
    };

    if (values.hasSecondaryPickupPostalCode) {
      createOrUpdateShipment.secondaryPickupPostalCode = values.secondaryPickupPostalCode;
    }

    if (values.hasSecondaryDestinationPostalCode) {
      createOrUpdateShipment.secondaryDestinationPostalCode = values.secondaryDestinationPostalCode;
    }

    let newShipmentId;
    if (isNewShipment) {
      createMTOShipment(createOrUpdateShipment)
        .then(() => {
          setSubmitting(false);
          newShipmentId = '00000000-0000-0000-0000-000000000000'; // TODO: replace me
          history.push(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId,
              mtoShipmentId: newShipmentId,
            }),
          );
        })
        .catch(() => {
          setSubmitting(false);
        });
    } else {
      createOrUpdateShipment.id = mtoShipment.id;
      createOrUpdateShipment.ppmShipment.id = mtoShipment.ppmShipment?.id;

      patchMTOShipment(mtoShipment.id, createOrUpdateShipment, mtoShipment.eTag)
        .then(() => {
          setSubmitting(false);
          history.push(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId,
              mtoShipmentId: mtoShipment?.id,
            }),
          );
        })
        .catch(() => {
          setSubmitting(false);
        });
    }
  };

  return (
    <GridContainer>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <h1>PPM date & location</h1>
          <DateAndLocationForm
            mtoShipment={mtoShipment}
            serviceMember={serviceMember}
            destinationDutyStation={destinationDutyLocation}
            onSubmit={handleSubmit}
            onBack={handleBack}
            postalCodeValidator={validatePostalCode}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

DateAndLocation.propTypes = {
  mtoShipment: MtoShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyLocation: DutyStationShape.isRequired,
};

DateAndLocation.defaultProps = {
  mtoShipment: {},
};

export default DateAndLocation;
