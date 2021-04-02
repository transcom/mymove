import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import ScrollToTop from 'components/ScrollToTop';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import CurrentDutyStationForm from 'components/Customer/CurrentDutyStationForm/CurrentDutyStationForm';
import { customerRoutes } from 'constants/routes';
import { ServiceMemberShape } from 'types/customerShapes';
import { DutyStationShape } from 'types/dutyStation';

const dutyStationFormName = 'duty_station';

export const DutyStation = ({ serviceMember, existingStation, newDutyStation, updateServiceMember, push }) => {
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    current_station: existingStation.name ? existingStation : {},
  };

  const handleBack = () => {
    push(customerRoutes.CONTACT_INFO_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      current_station_id: values.current_station.id,
    };

    return patchServiceMember(payload)
      .then((response) => {
        updateServiceMember(response);
        push(customerRoutes.CURRENT_ADDRESS_PATH);
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        setServerError(errorMessage);
      });
  };

  return (
    <GridContainer>
      <ScrollToTop otherDep={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <CurrentDutyStationForm
            onSubmit={handleSubmit}
            onBack={handleBack}
            initialValues={initialValues}
            newDutyStation={newDutyStation}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

DutyStation.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
  push: PropTypes.func.isRequired,
  existingStation: DutyStationShape,
  newDutyStation: DutyStationShape,
};

DutyStation.defaultProps = {
  existingStation: {},
  newDutyStation: {},
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  const orders = selectCurrentOrders(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    values: getFormValues(dutyStationFormName)(state),
    existingStation: serviceMember?.current_station,
    serviceMember,
    newDutyStation: orders?.new_duty_station,
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(requireCustomerState(DutyStation, profileStates.CONTACT_INFO_COMPLETE));
