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
        iconStyles.statusIcon,
        { [iconStyles.statusDefaultColor]: isTspSite },
        { [iconStyles.statusAttention]: !isTspSite },
      )}
      icon={faClock}
    />
  );
};

export default SitStatusIcon;
