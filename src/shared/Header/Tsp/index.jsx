import React from 'react';
import { NavLink } from 'react-router-dom';

import LoginButton from 'shared/User/LoginButton';
import Email from 'shared/User/Email';

import './index.css';

function TspHeader() {
  return (
    <header role="banner" className="header">
      <div className="tspHeaderOne">
        <div className="usa-logo" id="basic-logo">
          <em className="usa-logo-text">
            <NavLink to="/" title="Home" aria-label="Transcom PPP TSP Home">
              tsp.move.mil
            </NavLink>
          </em>
        </div>
      </div>
      <div className="tspHeaderTwo">
        <ul className="usa-nav-primary">
          <li>
            <Email />
          </li>
          <li>
            <LoginButton />
          </li>
        </ul>
      </div>
    </header>
  );
}

export default TspHeader;
