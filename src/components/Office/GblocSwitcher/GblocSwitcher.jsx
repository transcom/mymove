import React, { useContext, useEffect, useState } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import styles from './GblocSwitcher.module.scss';

import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';
import { useListGBLOCsQueries } from 'hooks/queries';
import { selectLoggedInUser } from 'store/entities/selectors';
import SelectedGblocContext, {
  SELECTED_GBLOC_SESSION_STORAGE_KEY,
} from 'components/Office/GblocSwitcher/SelectedGblocContext';
// import { user } from 'shared/Entities/schema';
import { roleTypes } from 'constants/userRoles';
// import { officeUser } from 'shared/Entities/schema';

const GBLOCSwitcher = ({ officeUser, activeRole, ariaLabel }) => {
  const [isInitialPageLoad, setIsInitialPageLoad] = useState(true);
  const { selectedGbloc, handleGblocChange } = useContext(SelectedGblocContext);

  let { result: gblocs } = useListGBLOCsQueries();
  if (activeRole !== roleTypes.HQ) {
    gblocs = officeUser?.transportation_office_assignments.map((toa) => {
      return toa?.transportationOffice?.gbloc;
    });
  }

  const officeUsersDefaultGbloc = officeUser.transportation_office.gbloc;
  if (gblocs.indexOf(officeUsersDefaultGbloc) === -1) {
    gblocs.push(officeUsersDefaultGbloc);
  }

  useEffect(() => {
    if (window.sessionStorage.getItem(SELECTED_GBLOC_SESSION_STORAGE_KEY) && isInitialPageLoad) {
      handleGblocChange(window.sessionStorage.getItem(SELECTED_GBLOC_SESSION_STORAGE_KEY));
      setIsInitialPageLoad(false);
    } else if (isInitialPageLoad) {
      handleGblocChange(officeUsersDefaultGbloc);
      setIsInitialPageLoad(false);
    }
  }, [selectedGbloc, officeUsersDefaultGbloc, isInitialPageLoad, handleGblocChange]);

  return (
    <ButtonDropdown
      onChange={(e) => {
        handleGblocChange(e.target.value);
      }}
      value={selectedGbloc || officeUsersDefaultGbloc}
      ariaLabel={ariaLabel}
      divClassName={styles.switchGblocButton}
      testId="gbloc_switcher"
    >
      {gblocs.map((gbloc) => (
        <option value={gbloc} key={`filterOption_${gbloc}`}>
          {gbloc}
        </option>
      ))}
    </ButtonDropdown>
  );
};

GBLOCSwitcher.defaultProps = {
  ariaLabel: '',
};

GBLOCSwitcher.propTypes = {
  officeUser: PropTypes.object.isRequired,
  activeRole: PropTypes.string.isRequired,
  // officeUsersDefaultGbloc: PropTypes.string.isRequired,
  ariaLabel: PropTypes.string,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
    activeRole: state.auth.activeRole,
  };
};

export default connect(mapStateToProps)(GBLOCSwitcher);
