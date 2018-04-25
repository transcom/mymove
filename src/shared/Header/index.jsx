import React from 'react';
import { NavLink } from 'react-router-dom';

import LoginButton from 'shared/User/LoginButton';
import Email from 'shared/User/Email';

import usaFlag from 'shared/images/us-flag.png';
import govIcon from 'shared/images/icon-dot-gov.svg';
import sslIcon from 'shared/images/icon-https.svg';

function Header() {
  return (
    <header className="usa-header usa-header-basic" role="banner">
      {/* Gov banner BEGIN */}
      <div className="usa-banner">
        <div className="usa-accordion">
          <header className="usa-banner-header">
            <div className="usa-grid usa-banner-inner">
              <img src={usaFlag} alt="U.S. flag" />
              <p>An official website of the United States government</p>
              <button
                className="usa-accordion-button usa-banner-button"
                aria-expanded="false"
                aria-controls="gov-banner"
              >
                <span className="usa-banner-button-text">
                  Here's how you know
                </span>
              </button>
            </div>
          </header>
          <div
            className="usa-banner-content usa-grid usa-accordion-content"
            id="gov-banner"
          >
            <div className="usa-banner-guidance-gov usa-width-one-half">
              <img
                className="usa-banner-icon usa-media_block-img"
                src={govIcon}
                alt="Dot gov"
              />
              <div className="usa-media_block-body">
                <p>
                  <strong>The .mil means it’s official.</strong>
                  <br />
                  Federal government websites always use a .gov or .mil domain.
                  Before sharing sensitive information online, make sure you’re
                  on a .gov or .mil site by inspecting your browser’s address
                  (or “location”) bar.
                </p>
              </div>
            </div>
            <div className="usa-banner-guidance-ssl usa-width-one-half">
              <img
                className="usa-banner-icon usa-media_block-img"
                src={sslIcon}
                alt="SSL"
              />
              <div className="usa-media_block-body">
                <p>
                  This site is also protected by an SSL (Secure Sockets Layer)
                  certificate that’s been signed by the U.S. government. The{' '}
                  <strong>https://</strong> means all transmitted data is
                  encrypted &nbsp;— in other words, any information or browsing
                  history that you provide is transmitted securely.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
      {/* Gov banner END */}
      <div className="usa-nav-container">
        <div className="usa-navbar">
          <button className="usa-menu-btn">Menu</button>
          <div className="usa-logo" id="basic-logo">
            <em className="usa-logo-text">
              <NavLink to="/" title="Home" aria-label="Transcom PPP Home">
                Transcom PPP
              </NavLink>
            </em>
          </div>
        </div>
        <nav className="usa-nav">
          <ul className="usa-nav-primary usa-accordion">
            <li>
              <NavLink to="/Feedback" className="usa-nav-link">
                <span>Report a bug</span>
              </NavLink>
            </li>
            <li>
              <Email />
            </li>
            <li>
              <LoginButton />
            </li>
          </ul>
          <form className="usa-search usa-search-small">
            <div role="search">
              <label className="usa-sr-only" htmlFor="search-field-small">
                Search small
              </label>
              <input id="search-field-small" type="search" name="search" />
              <button type="submit">
                <span className="usa-sr-only">Search</span>
              </button>
            </div>
          </form>
        </nav>
      </div>
    </header>
  );
}

export default Header;
