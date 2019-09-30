import React from 'react';
import { NavLink } from 'react-router-dom';

import LoginButton from 'shared/User/LoginButton';
import UserGreeting from 'shared/User/UserGreeting';

import usaFlag from 'shared/images/us-flag.png';
import govIcon from 'shared/images/icon-dot-gov.svg';
import sslIcon from 'shared/images/icon-https.svg';
import './index.css';
function Header() {
  return (
    <div>
      {/* Gov banner BEGIN */}
      <div className="usa-banner">
        <div className="usa-accordion">
          <header className="usa-banner__header">
            <div className="usa-banner__inner">
              <div className="grid-col-auto">
                <img className="usa-banner__header-flag" src={usaFlag} alt="U.S. flag" />
              </div>
              <div className="grid-col-fill tablet:grid-col-auto">
                <p className="usa-banner__header-text">An official website of the United States government</p>
                <p className="usa-banner__header-action" aria-hidden={true}>
                  Here’s how you know
                </p>
              </div>
              <button
                className="usa-accordion__button usa-banner__button"
                aria-expanded={false}
                aria-controls="gov-banner"
              >
                <span className="usa-banner__button-text">Here’s how you know</span>
              </button>
            </div>
          </header>
          <div className="usa-banner__content usa-accordion__content" id="gov-banner" hidden>
            <div className="grid-row grid-gap-lg">
              <div className="usa-banner__guidance tablet:grid-col-6">
                <img className="usa-banner__icon usa-media-block__img" src={govIcon} alt="Dot gov" />
                <div className="usa-media-block__body">
                  <p>
                    <strong>The .mil means it’s official.</strong>
                    <br />
                    Federal government websites always use a .gov or .mil domain. Before sharing sensitive information
                    online, make sure you’re on a .gov or .mil site by inspecting your browser’s address (or “location”)
                    bar.
                  </p>
                </div>
              </div>
              <div className="usa-banner__guidance tablet:grid-col-6">
                <img className="usa-banner__icon usa-media-block__img" src={sslIcon} alt="Https" />
                <div className="usa-media-block__body">
                  <p>
                    <strong>The site is secure.</strong>
                    <br />
                    This site is also protected by an SSL (Secure Sockets Layer) certificate that’s been signed by the
                    U.S. government. The <strong>https://</strong> means all transmitted data is encrypted &nbsp;— in
                    other words, any information or browsing history that you provide is transmitted securely.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Gov banner END */}

      {/* <div className="usa-nav-container">
        <div className="usa-accordion">
          <header className="usa-banner-header">
            <div className="usa-grid usa-banner-inner">
              <button
                className="usa-accordion__button usa-banner-button"
                aria-expanded="false"
                aria-controls="gov-banner"
              >
                <span className="usa-banner-button-text">Here's how you know</span>
              </button>
            </div>
          </header>
        </div>
      </div>
      {/* Gov banner END */}
      <header className="usa-header usa-header--basic" role="banner">
        <div className="my-move-header">
          <div className="usa-nav-container">
            <div className="usa-navbar">
              <div className="usa-logo" id="basic-logo">
                <em className="usa-logo__text">
                  <NavLink to="/" title="my.move.mil" aria-label="my.move.mil">
                    my.move.mil
                  </NavLink>
                </em>
              </div>
            </div>
            <nav className="usa-nav my-move-login">
              <ul className="usa-nav__primary usa-accordion">
                <li className="usa-nav__primary-item">
                  <UserGreeting />
                </li>
                <li className="usa-nav__primary-item">
                  <LoginButton />
                </li>
              </ul>
            </nav>
            {/*<div className="my-move-login">
            <UserGreeting />
            <LoginButton />
          </div>*/}
          </div>
        </div>
      </header>
    </div>
  );
}

export default Header;
