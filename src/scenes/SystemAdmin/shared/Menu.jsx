import React from 'react';
import { useResourceDefinitions, MenuItemLink } from 'react-admin';
import ExitIcon from '@material-ui/icons/PowerSettingsNew';
import { LogoutUser } from 'utils/api';

const Menu = () => {
  const resourcesDefinitions = useResourceDefinitions();
  const resources = Object.keys(resourcesDefinitions).map((name) => resourcesDefinitions[name]);
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
          LogoutUser().then((r) => {
            const redirectURL = r.body;
            if (redirectURL) {
              window.location.href = redirectURL;
            } else {
              window.localStorage.setItem('hasLoggedOut', true);
              window.location.href = '/';
            }
          });
        }}
      />
    </div>
  );
};

export default Menu;
