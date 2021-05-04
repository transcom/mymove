import React from 'react';
import { connect } from 'react-redux';
import { MenuItemLink, getResources } from 'react-admin';
import { withRouter } from 'react-router-dom';
import ExitIcon from '@material-ui/icons/PowerSettingsNew';
import { LogoutUser } from 'utils/api';
import { createBrowserHistory } from 'history';

const Menu = (props) => {
  const resources = props.resources;
  const history = createBrowserHistory({ basename: '' });
  return (
    <div>
      {resources
        .filter((resource) => resource.hasList || resource.hasShow)
        .map((resource) => (
          <MenuItemLink
            key={resource.name}
            to={`/${resource.name}`}
            primaryText={(resource.options && resource.options.label) || resource.name}
          />
        ))}
      <MenuItemLink
        to="/"
        primaryText="Logout"
        leftIcon={<ExitIcon />}
        onClick={(e) => {
          e.preventDefault();
          LogoutUser().then(() => {
            console.log('logoutuser then -- Menu.jsx');
            history.push({
              pathname: '/sign-in',
              state: { hasLoggedOut: true },
            });
          });
        }}
      />
    </div>
  );
};

const mapStateToProps = (state) => ({
  resources: getResources(state),
});

export default withRouter(connect(mapStateToProps)(Menu));
