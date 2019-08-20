import React from 'react';
import PropTypes from 'prop-types';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import { kebabCase } from 'lodash';

const SitAction = props => {
  const dataCyName = 'sit-' + kebabCase(props.action) + '-link';

  return (
    <span className="sit-action">
      <a onClick={props.onClick} data-cy={dataCyName}>
        {props.icon && <FontAwesomeIcon className="icon" icon={props.icon} />}
        {props.action}
      </a>
    </span>
  );
};

SitAction.propTypes = {
  action: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  icon: PropTypes.object,
};

export default SitAction;
