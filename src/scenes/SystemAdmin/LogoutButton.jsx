import React from 'react';
import { connect } from 'react-redux';
import { Responsive, userLogout } from 'react-admin';
import MenuItem from '@material-ui/core/MenuItem';
import Button from '@material-ui/core/Button';
import ExitIcon from '@material-ui/icons/PowerSettingsNew';

const LogoutButton = ({ userLogout, ...rest }) => (
  <Responsive
    xsmall={
      <MenuItem onClick={userLogout} {...rest}>
        <ExitIcon /> Logout
      </MenuItem>
    }
    medium={
      <Button onClick={userLogout} size="small" {...rest}>
        <ExitIcon /> Logout
      </Button>
    }
  />
);
const redirectTo = '/';
const customUserLogout = () => userLogout(redirectTo);
export default connect(undefined, { userLogout: customUserLogout })(LogoutButton);
