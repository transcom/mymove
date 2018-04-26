import React from 'react';
import { NavLink } from 'react-router-dom';

import LoginButton from 'shared/User/LoginButton';

function QueueHeader() {
  return (
    <header className="usa-header usa-header-basic" role="banner">
      <div className="usa-nav-container">
        <div className="usa-navbar">
          <button className="usa-menu-btn">Menu</button>
          <div className="usa-logo" id="basic-logo">
            <em className="usa-logo-text">
              <NavLink to="/" title="Home" aria-label="Admin Move.Mil">
                admin.move.mil
              </NavLink>
            </em>
          </div>
        </div>
        <nav className="usa-nav">
          <ul className="usa-nav-primary usa-accordion">
            <li>
              <NavLink to="/" className="usa-nav-link">
                <span>Queue List</span>
              </NavLink>
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

export default QueueHeader;
