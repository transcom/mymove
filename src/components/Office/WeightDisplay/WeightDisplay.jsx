import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';
import { connect } from 'react-redux';

import styles from 'components/Office/WeightDisplay/WeightDisplay.module.scss';
import { formatWeight } from 'utils/formatters';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { withContext } from 'shared/AppContext';
import { roleTypes } from 'constants/userRoles';

const WeightDisplay = ({ heading, weightValue, onEdit, children, activeRole }) => {
  return (
    <div className={classnames('maxw-tablet', styles.WeightDisplay)}>
      <div className={styles.heading}>
        <div>{heading}</div>
        {onEdit && (
          <Restricted to={permissionTypes.updateBillableWeight}>
            {activeRole !== roleTypes.SERVICES_COUNSELOR && (
              <Button
                unstyled
                type="button"
                className={styles.editButton}
                onClick={onEdit}
                data-testid="weightDisplayEdit"
              >
                <FontAwesomeIcon icon="pen" title="edit" alt="" />
              </Button>
            )}
          </Restricted>
        )}
      </div>
      <div data-testid="weight-display" className={styles.value}>
        {Number.isFinite(weightValue) ? formatWeight(weightValue) : '—'}
      </div>
      {children && <div className={styles.details}>{children}</div>}
    </div>
  );
};

WeightDisplay.propTypes = {
  heading: PropTypes.string.isRequired,
  weightValue: PropTypes.number,
  children: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
  onEdit: PropTypes.func,
  activeRole: PropTypes.string,
};

WeightDisplay.defaultProps = {
  weightValue: null,
  children: null,
  onEdit: null,
  activeRole: null,
};

// Checks user role such that Service Counselors cannot modify service items while on MTO read-only page
const mapStateToProps = (state) => {
  return {
    activeRole: state.auth.activeRole,
  };
};

export default withContext(connect(mapStateToProps)(WeightDisplay));
