import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import DodInfoForm from 'components/Customer/DodInfoForm/DodInfoForm';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectOktaUser, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';
import { ServiceMemberShape } from 'types/customerShapes';

export const DodInfo = ({ updateServiceMember, serviceMember, oktaUser }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    affiliation: serviceMember?.affiliation || '',
    edipi: oktaUser?.cac_edipi || '',
    emplid: serviceMember?.emplid || '',
  };

  const handleNext = () => {
    navigate(customerRoutes.NAME_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      affiliation: values.affiliation,
      edipi: values.edipi,
      emplid: values.affiliation === 'COAST_GUARD' ? values.emplid : null,
    };

    return patchServiceMember(payload)
      .then(updateServiceMember)
      .then(handleNext)
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        let errorMessage;
        if (e.response.body.message === 'Unhandled data error encountered') {
          errorMessage = 'This EMPLID is already in use';
        } else {
          errorMessage = getResponseError(response, 'failed to update service member due to server error');
        }

        setServerError(errorMessage);
      });
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

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
          <DodInfoForm initialValues={initialValues} onSubmit={handleSubmit} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

DodInfo.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  oktaUser: selectOktaUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(requireCustomerState(DodInfo, profileStates.EMPTY_PROFILE));
