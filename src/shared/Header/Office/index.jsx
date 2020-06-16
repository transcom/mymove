import React from 'react';
import { withRouter } from 'react-router';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';

import LoginButton from 'shared/User/LoginButton';

import './index.css';
import TabNav from 'components/TabNav';
import propTypes from 'prop-types';

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

function MoveTabNav(props) {
  // Should in the future become a move locator
  const moveOrderId = props.match.params.moveId;

  return (
    <header className="usa-header nav-header" role="navigation">
      <div className="grid-container-desktop-lg">
        <TabNav
          items={[
            <NavLink exact activeClassName="usa-current" className="usa-nav__link" to={`/moves/${moveOrderId}/details`}>
              <span className="tab-title">Move details</span>
              <Tag>2</Tag>
            </NavLink>,
            <NavLink exact activeClassName="usa-current" className="usa-nav__link" to={`/moves/${moveOrderId}/mto`}>
              <span className="tab-title">Move task order</span>
            </NavLink>,
            <NavLink
              exact
              activeClassName="usa-current"
              className="usa-nav__link"
              to={`/moves/${moveOrderId}/payment-requests`}
            >
              <span className="tab-title">Payment requests</span>
            </NavLink>,
            <NavLink exact activeClassName="usa-current" className="usa-nav__link" to={`/moves/${moveOrderId}/history`}>
              <span className="tab-title">History</span>
            </NavLink>,
          ]}
        ></TabNav>
      </div>
    </header>
  );
}

MoveTabNav.propTypes = {
  match: propTypes.object.isRequired,
};

const MoveTabNavWithRouter = withRouter(MoveTabNav);

export { QueueHeader, MoveTabNavWithRouter };
