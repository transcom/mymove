import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { useNavigate } from 'react-router-dom';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { getResponseError, patchServiceMember } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';
import ResidentialAddressForm from 'components/Customer/ResidentialAddressForm/ResidentialAddressForm';
import { ResidentialAddressShape } from 'types/address';

export const ResidentialAddress = ({ serviceMember, updateServiceMember }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  const formFieldsName = 'current_residence';

  const initialValues = {
    [formFieldsName]: {
      streetAddress1: serviceMember.residential_address?.streetAddress1 || '',
      streetAddress2: serviceMember.residential_address?.streetAddress2 || '',
      streetAddress3: serviceMember.residential_address?.streetAddress3 || '',
      city: serviceMember.residential_address?.city || '',
      state: serviceMember.residential_address?.state || '',
      postalCode: serviceMember.residential_address?.postalCode || '',
      county: serviceMember.residential_address?.county || '',
      usPostRegionCitiesID: serviceMember.residential_address?.usPostRegionCitiesID || '',
    },
  };

  const handleBack = () => {
    navigate(customerRoutes.CONTACT_INFO_PATH);
  };

  const handleNext = () => {
    navigate(customerRoutes.BACKUP_ADDRESS_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      residential_address: values.current_residence,
    };

    return patchServiceMember(payload)
      .then(updateServiceMember)
      .then(handleNext)
      .catch((e) => {
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');

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
          <ResidentialAddressForm
            formFieldsName={formFieldsName}
            initialValues={initialValues}
            onBack={handleBack}
            onSubmit={handleSubmit}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};
ResidentialAddress.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: PropTypes.shape({
    id: PropTypes.string.isRequired,
    residential_address: ResidentialAddressShape,
  }).isRequired,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
});

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(requireCustomerState(ResidentialAddress, profileStates.CONTACT_INFO_COMPLETE));
