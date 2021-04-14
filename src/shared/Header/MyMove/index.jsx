import React from 'react';
import { GovBanner } from '@trussworks/react-uswds';
import { Link } from 'react-router-dom';

import LoginButton from 'containers/LoginButton/LoginButton';
import MilMoveLogo from 'shared/images/milmove-logo.svg';

import BypassBlock from 'components/BypassBlock';

import './index.scss';

function Header() {
  return (
    <div>
      <BypassBlock />
      <GovBanner />

      <header className="usa-header usa-header--basic" role="banner">
        <div className="my-move-header">
          <div className="usa-nav-container">
            <div className="usa-navbar">
              <div className="usa-logo" id="basic-logo">
                <em className="usa-logo__text">
                  <Link to="/" title="my.move.mil" aria-label="my.move.mil">
                    <img src={MilMoveLogo} alt="MilMove" />
                  </Link>
                </em>
              </div>
              <button className="usa-menu-btn">Menu</button>
            </div>
            <nav className="usa-nav" aria-label="Primary navigation">
              <ul className="usa-nav__primary usa-accordion">
                <LoginButton />
              </ul>
            </nav>
          </div>
        </div>
      </header>
    </div>
  );
}

export default Header;
