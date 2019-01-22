import React from 'react';
import { NavLink } from 'react-router-dom';

import LoginButton from 'shared/User/LoginButton';
import UserGreeting from 'shared/User/UserGreeting';

import './index.css';

function QueueHeader() {
  return (
    <header role="banner" className="header">
      <div className="officeHeaderOne">
        <div className="usa-logo" id="basic-logo">
          <em className="usa-logo-text">
            <NavLink to="/" title="Home" aria-label="Transcom PPP Office Home">
              office.move.mil
            </NavLink>
          </em>
        </div>
      </div>
      <div className="officeHeaderTwo">
        <ul className="usa-nav-primary">
          <li>
            <UserGreeting />
          </li>
          <li>
            <LoginButton />
          </li>
        </ul>
      </div>
    </header>
  );
}

export default QueueHeader;
