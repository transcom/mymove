import React, { useContext, useEffect, useState } from 'react';
import { connect, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';

import styles from './GblocSwitcher.module.scss';

import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';
// import { selectIsSettingActiveOffice } from 'store/auth/selectors';
import { useListGBLOCsQueries } from 'hooks/queries';
import { selectLoggedInUser } from 'store/entities/selectors';
import SelectedGblocContext, {
  SELECTED_GBLOC_SESSION_STORAGE_KEY,
} from 'components/Office/GblocSwitcher/SelectedGblocContext';
import { roleTypes } from 'constants/userRoles';
import { setActiveOffice as setActiveOfficeAction } from 'store/auth/actions';
import { UpdateActiveOfficeServerSession } from 'utils/api';

const GBLOCSwitcher = ({ officeUser, activeRole, ariaLabel, setActiveOffice }) => {
  // const navigate = useNavigate();
  const [isInitialPageLoad, setIsInitialPageLoad] = useState(true);
  const { selectedGbloc, handleGblocChange } = useContext(SelectedGblocContext);
  // const isSettingActiveOffice = useSelector(selectIsSettingActiveOffice);
  // const [pendingOffice, setPendingOffice] = useState(null);
  // debugger;
  const handleChange = (activeGbloc) => {
    const activeOffice = officeUser.transportation_office_assignments?.find(
      (office) => office.transportationOffice.gbloc === activeGbloc,
    );
    UpdateActiveOfficeServerSession(activeOffice.transportationOffice.id).then((res) => {
      handleGblocChange(activeGbloc);
      console.log('id = ', activeOffice.transportationOffice.id);
      // debugger;
      // setPendingOffice(activeOffice.transportationOffice.id);
      console.log('res = ', res);
      // debugger;
      setActiveOffice(activeOffice.transportationOffice);
    });
  };

  let { result: gblocs } = useListGBLOCsQueries();
  if (activeRole !== roleTypes.HQ) {
    gblocs = officeUser?.transportation_office_assignments?.map((toa) => {
      return toa?.transportationOffice?.gbloc;
    });
  }

  const officeUsersDefaultGbloc = officeUser.transportation_office?.gbloc;
  if (gblocs?.indexOf(officeUsersDefaultGbloc) === -1) {
    gblocs.push(officeUsersDefaultGbloc);
  }

  // useEffect(() => {
  //   if (pendingOffice !== null && !isSettingActiveOffice) {
  //     // Pending role has been set and the auth saga
  //     // has received a response from the server.
  //     // This prevents a saga/action race condition between
  //     // index rendering on '/' and the SelectApplication component.
  //     // Previous race condition:
  //     // Select application requests to update the server session,
  //     // select application routes to index before hearing back,
  //     // index requests the current logged in user
  //     // the server begins handling the index with the old AppCtx session
  //     // the old AppCtx session has not been finished updating via SetActiveRole saga
  //     // then index is now rendering the old role, not the new role, thus a race condition
  //     setPendingOffice(null);
  //     // navigate('/');
  //   }
  // }, [pendingOffice, isSettingActiveOffice, navigate]);

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
        handleChange(e.target.value);
      }}
      value={selectedGbloc || officeUsersDefaultGbloc}
      ariaLabel={ariaLabel}
      divClassName={styles.switchGblocButton}
      testId="gbloc_switcher"
    >
      {gblocs?.map((gbloc) => (
        <option value={gbloc} key={`filterOption_${gbloc}`}>
          {gbloc}
        </option>
      ))}
    </ButtonDropdown>
  );
};

GBLOCSwitcher.defaultProps = {
  ariaLabel: 'Switch to a different GBLOC',
};

GBLOCSwitcher.propTypes = {
  officeUser: PropTypes.object.isRequired,
  activeRole: PropTypes.string.isRequired,
  ariaLabel: PropTypes.string,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
    activeRole: state.auth.activeRole,
  };
};

const mapDispatchToProps = {
  setActiveOffice: setActiveOfficeAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(GBLOCSwitcher);
