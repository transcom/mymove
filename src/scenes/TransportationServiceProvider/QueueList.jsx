import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

export default class QueueList extends Component {
  render() {
    return (
      <div>
        <h2>Queues</h2>
        <ul className="usa-sidenav-list">
          <li>
            <NavLink to="/queues/all" activeClassName="usa-current">
              <span>All Moves</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/other" activeClassName="usa-current">
              <span>Other</span>
            </NavLink>
          </li>
        </ul>
      </div>
    );
  }
}
