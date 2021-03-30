import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ScrollToTop from 'components/ScrollToTop';
import { getResponseError, patchServiceMember } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes } from 'constants/routes';
import BackupMailingAddressForm from 'components/Customer/BackupMailingAddressForm/BackupMailingAddressForm';
import { ResidentialAddressShape } from 'types/address';

export const BackupMailingAddress = ({ serviceMember, updateServiceMember, push }) => {
  const [serverError, setServerError] = useState(null);

  const formFieldsName = 'backup_mailing_address';

  const initialValues = {
    [formFieldsName]: {
      street_address_1: serviceMember.backup_mailing_address?.street_address_1 || '',
      street_address_2: serviceMember.backup_mailing_address?.street_address_2 || '',
      city: serviceMember.backup_mailing_address?.city || '',
      state: serviceMember.backup_mailing_address?.state || '',
      postal_code: serviceMember.backup_mailing_address?.postal_code || '',
    },
  };

  const handleBack = () => {
    push(customerRoutes.CURRENT_ADDRESS_PATH);
  };

  const handleNext = () => {
    push(customerRoutes.BACKUP_CONTACTS_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      backup_mailing_address: values.backup_mailing_address,
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
          <BackupMailingAddressForm
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

BackupMailingAddress.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: PropTypes.shape({
    id: PropTypes.string.isRequired,
    backup_mailing_address: ResidentialAddressShape,
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
)(requireCustomerState(BackupMailingAddress, profileStates.ADDRESS_COMPLETE));
