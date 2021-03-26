import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';

import ScrollToTop from 'components/ScrollToTop';
import DodInfoForm from 'components/Customer/DodInfoForm/DodInfoForm';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';
import { ServiceMemberShape } from 'types/customerShapes';

export const DodInfo = ({ updateServiceMember, serviceMember, push }) => {
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    affiliation: serviceMember?.affiliation || '',
    edipi: serviceMember?.edipi || '',
    rank: serviceMember?.rank || '',
  };

  const handleBack = () => {
    push(customerRoutes.CONUS_OCONUS_PATH);
  };

  const handleNext = () => {
    push(customerRoutes.NAME_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      affiliation: values.affiliation,
      edipi: values.edipi,
      rank: values.rank,
    };

    return patchServiceMember(payload)
      .then(updateServiceMember)
      .then(handleNext)
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
          <DodInfoForm initialValues={initialValues} onSubmit={handleSubmit} onBack={handleBack} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

DodInfo.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  push: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(requireCustomerState(DodInfo, profileStates.EMPTY_PROFILE));
