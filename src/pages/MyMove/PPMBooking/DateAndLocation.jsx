import React from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
import { DutyStationShape } from 'types';
import DateAndLocationForm from 'components/Customer/PPMBooking/DateAndLocationForm/DateAndLocationForm';
import { validatePostalCode } from 'utils/validation';
import { customerRoutes, generalRoutes } from 'constants/routes';

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

  const handleSubmit = () => {};

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
