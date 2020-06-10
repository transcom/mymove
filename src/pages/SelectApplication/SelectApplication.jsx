import React from 'react';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { roleTypes } from 'constants/userRoles';

const SelectApplication = ({ setActiveRole, activeRole }) => {
  const history = useHistory();

  const handleSelectRole = (roleType) => {
    setActiveRole(roleType);
    history.push('/');
  };

  return (
    <>
      {/* eslint-disable-next-line react/jsx-one-expression-per-line */}
      <h2>Active role: {activeRole}</h2>
      <Button
        type="button"
        onClick={() => {
          handleSelectRole(roleTypes.PPM);
        }}
      >
        PPM Move Queue
      </Button>
      <br />
      <Button
        type="button"
        onClick={() => {
          handleSelectRole(roleTypes.TOO);
        }}
      >
        TOO Move Queue
      </Button>
      <br />
      <Button
        type="button"
        onClick={() => {
          handleSelectRole(roleTypes.TIO);
        }}
      >
        TIO Payment Request Queue
      </Button>
    </>
  );
};

SelectApplication.propTypes = {
  activeRole: PropTypes.string,
  setActiveRole: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  return {
    activeRole: state.auth.activeRole,
  };
};

const mapDispatchToProps = {
  setActiveRole: setActiveRoleAction,
};

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(SelectApplication);

export default SelectApplication;
