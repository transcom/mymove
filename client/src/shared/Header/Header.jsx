import React from 'react';

function Header() {
  return (
    <header className="usa-header usa-header-basic" role="banner">
      <div className="usa-nav-container">
        <div className="usa-navbar">
          <div className="usa-logo" id="basic-logo">
            <em className="usa-logo-text">
              <a href="/" title="Home" aria-label="Home">
                Transcom PPP
              </a>
            </em>
          </div>
        </div>
      </div>
    </header>
  );
}

export default Header;
