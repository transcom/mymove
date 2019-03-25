import React from 'react';
import classnames from 'classnames';

import iconStyles from 'shared/styles/icons.module.scss';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';

const SitStatusIcon = props => {
  const { isTspSite } = props;
  return (
    <FontAwesomeIcon
      className={classnames(
        iconStyles['status-icon'],
        { [iconStyles['status-default-color']]: isTspSite },
        { [iconStyles['status-attention']]: !isTspSite },
      )}
      icon={faClock}
    />
  );
};

export default SitStatusIcon;
