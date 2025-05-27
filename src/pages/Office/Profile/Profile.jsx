import React from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { Link } from 'react-router-dom';

import styles from './Profile.module.scss';
import 'styles/office.scss';

import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { OfficeUserInfoShape } from 'types/index';
import { selectLoggedInUser } from 'store/entities/selectors';
import ContactInfoDisplay from 'components/Office/Profile/ContactInfoDisplay';
import { formatFullName } from 'utils/formatters';
import { generalRoutes, officeRoutes } from 'constants/routes';

const Profile = ({ officeUser }) => {
  const officeUserInfo = {
    name: formatFullName(officeUser?.first_name, officeUser?.middle_name, officeUser?.last_name),
    telephone: officeUser?.telephone,
    email: officeUser?.email,
  };

  return (
    <div className={styles.Profile}>
      <GridContainer>
        <Grid col={4} desktop={{ col: 12 }} tablet={{ col: 8 }}>
          <Link to={generalRoutes.HOME_PATH}>Return to Dashboard</Link>
          <ConnectedFlashMessage />
          <div>
            <h1>Profile</h1>
          </div>
          <SectionWrapper className={formStyles.formSection}>
            <ContactInfoDisplay officeUserInfo={officeUserInfo} editURL={officeRoutes.CONTACT_INFO_EDIT_PATH} />
          </SectionWrapper>
        </Grid>
      </GridContainer>
    </div>
  );
};

Profile.propTypes = {
  officeUser: OfficeUserInfoShape,
};

Profile.defaultProps = {
  officeUser: {},
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
  };
};

export default connect(mapStateToProps)(Profile);
