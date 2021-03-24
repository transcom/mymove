import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
// import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import ScrollToTop from 'components/ScrollToTop';
import ContactInfoForm from 'components/Customer/ContactInfoForm';
// import ServiceMemberShape from 'types/customerShapes';
// import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';

export const ContactInfo = () => {
  const initialValues = {
    telephone: '',
    secondary_phone: '',
    personal_email: '',
    phone_is_preferred: null,
    email_is_preferred: null,
  };

  const handleSubmit = () => {
    console.log('do something'); // eslint-disable-line no-console
  };

  const handleBack = () => {
    console.log('do another thing'); // eslint-disable-line no-console
  };

  return (
    <GridContainer>
      <ScrollToTop />
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <ContactInfoForm initialValues={initialValues} onBack={handleBack} onSubmit={handleSubmit} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

ContactInfo.propTypes = {
  // serviceMember: ServiceMemberShape.isRequired,
  // updateServiceMember: PropTypes.func.isRequired,
  // push: PropTypes.func.isRequired,
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
)(requireCustomerState(ContactInfo, profileStates.NAME_COMPLETE));
