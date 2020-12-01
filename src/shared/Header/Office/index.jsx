import React from 'react';
import { NavLink } from 'react-router-dom';

import LoginButton from 'shared/User/LoginButton';

import './index.css';

function QueueHeader() {
  return (
    <header className="usa-header usa-header--basic" role="banner">
      <div className="usa-nav-container header-widescreen">
        <div className="usa-navbar">
          <div className="usa-logo" id="basic-logo">
            <em className="usa-logo__text">
              <NavLink to="/" title="Home" aria-label="Transcom PPP Office Home">
                office.move.mil
              </NavLink>
            </em>
          </div>
        </div>
        <nav className="usa-nav" aria-label="Primary navigation">
          <ul className="usa-nav__primary usa-accordion">
            <LoginButton />
          </ul>
        </nav>
      </div>
    </header>
  );
}

export { QueueHeader };
