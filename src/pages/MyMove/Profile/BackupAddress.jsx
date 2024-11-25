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
import BackupAddressForm from 'components/Customer/BackupAddressForm/BackupAddressForm';
import { ResidentialAddressShape } from 'types/address';

export const BackupAddress = ({ serviceMember, updateServiceMember }) => {
  const navigate = useNavigate();

  const [serverError, setServerError] = useState(null);

  const formFieldsName = 'backup_mailing_address';

  const initialValues = {
    [formFieldsName]: {
      streetAddress1: serviceMember.backup_mailing_address?.streetAddress1 || '',
      streetAddress2: serviceMember.backup_mailing_address?.streetAddress2 || '',
      city: serviceMember.backup_mailing_address?.city || '',
      state: serviceMember.backup_mailing_address?.state || '',
      postalCode: serviceMember.backup_mailing_address?.postalCode || '',
      county: serviceMember.backup_mailing_address?.county || '',
      usPostRegionCitiesId: serviceMember.backup_mailing_address?.usPostRegionCitiesId || '',
    },
  };

  const handleBack = () => {
    navigate(customerRoutes.CURRENT_ADDRESS_PATH);
  };

  const handleNext = () => {
    navigate(customerRoutes.BACKUP_CONTACTS_PATH);
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
          <BackupAddressForm
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

BackupAddress.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: PropTypes.shape({
    id: PropTypes.string.isRequired,
    backup_mailing_address: ResidentialAddressShape,
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
)(requireCustomerState(BackupAddress, profileStates.ADDRESS_COMPLETE));
