import React from 'react';
import { string, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Link } from 'react-router-dom';
import classNames from 'classnames';

import { officeRoutes } from 'constants/routes';
import { formatOfficeProfileFirstAndLast } from 'utils/formatters';

const OfficeUserInfo = ({ handleLogout, firstName, lastName }) => {
  const content = formatOfficeProfileFirstAndLast({ firstName, lastName });

  return (
    <>
      <li className={classNames('usa-nav__primary-item')}>
        <Link to={officeRoutes.PROFILE_PATH} title="profile-link">
          {content}
        </Link>
      </li>
      <li className={classNames('usa-nav__primary-item')}>
        <Button unstyled onClick={handleLogout} type="button">
          Sign out
        </Button>
      </li>
    </>
  );
};

OfficeUserInfo.defaultProps = {
  firstName: null,
  lastName: null,
};

OfficeUserInfo.propTypes = {
  firstName: string,
  lastName: string,
  handleLogout: func.isRequired,
};

export default OfficeUserInfo;
