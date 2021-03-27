import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ScrollToTop from 'components/ScrollToTop';
import { getResponseError, patchServiceMember } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { ValidateZipRateData } from 'shared/api';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';
import ResidentialAddressForm from 'components/Customer/ResidentialAddressForm/ResidentialAddressForm';
import { ResidentialAddressShape } from 'types/address';

const UnsupportedZipCodeErrorMsg =
  'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.';

const validatePostalCode = async (value) => {
  if (!value) {
    return undefined;
  }

  let responseBody;
  try {
    responseBody = await ValidateZipRateData(value, 'origin');
  } catch (e) {
    return 'Error checking ZIP';
  }

  return responseBody.valid ? undefined : UnsupportedZipCodeErrorMsg;
};

export const ResidentialAddress = ({ serviceMember, updateServiceMember, push }) => {
  const [serverError, setServerError] = useState(null);

  const formFieldsName = 'current_residence';

  const initialValues = {
    [formFieldsName]: {
      street_address_1: serviceMember.residential_address?.street_address_1 || '',
      street_address_2: serviceMember.residential_address?.street_address_2 || '',
      city: serviceMember.residential_address?.city || '',
      state: serviceMember.residential_address?.state || '',
      postal_code: serviceMember.residential_address?.postal_code || '',
    },
  };

  const handleBack = () => {
    push(customerRoutes.CURRENT_DUTY_STATION_PATH);
  };

  const handleNext = () => {
    push(customerRoutes.BACKUP_ADDRESS_PATH);
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
          <ResidentialAddressForm
            formFieldsName={formFieldsName}
            initialValues={initialValues}
            onBack={handleBack}
            onSubmit={handleSubmit}
            validators={{ postalCode: validatePostalCode }}
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
  push: PropTypes.func.isRequired,
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
)(requireCustomerState(ResidentialAddress, profileStates.DUTY_STATION_COMPLETE));
