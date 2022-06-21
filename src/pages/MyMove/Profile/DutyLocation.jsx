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
import CurrentDutyLocationForm from 'components/Customer/CurrentDutyLocationForm/CurrentDutyLocationForm';
import { customerRoutes } from 'constants/routes';
import { ServiceMemberShape } from 'types/customerShapes';
import { DutyLocationShape } from 'types/dutyLocation';

const dutyLocationFormName = 'duty_location';

export const DutyLocation = ({ serviceMember, existingDutyLocation, newDutyLocation, updateServiceMember, push }) => {
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    current_location: existingDutyLocation.name ? existingDutyLocation : {},
  };

  const handleBack = () => {
    push(customerRoutes.CONTACT_INFO_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      current_location_id: values.current_location.id,
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
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <CurrentDutyLocationForm
            onSubmit={handleSubmit}
            onBack={handleBack}
            initialValues={initialValues}
            newDutyLocation={newDutyLocation}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

DutyLocation.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
  push: PropTypes.func.isRequired,
  existingDutyLocation: DutyLocationShape,
  newDutyLocation: DutyLocationShape,
};

DutyLocation.defaultProps = {
  existingDutyLocation: {},
  newDutyLocation: {},
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  const orders = selectCurrentOrders(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    values: getFormValues(dutyLocationFormName)(state),
    existingDutyLocation: serviceMember?.current_location,
    serviceMember,
    newDutyLocation: orders?.new_duty_location,
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(requireCustomerState(DutyLocation, profileStates.CONTACT_INFO_COMPLETE));
